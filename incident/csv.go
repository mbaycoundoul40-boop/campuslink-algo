package incident

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

// FormatDate = format des dates dans le CSV, ex. 2026-10-22 08:30.
const FormatDate = "2006-01-02 15:04"

var enTeteAttendu = []string{"id", "titre", "urgence", "criticite", "impact", "date"}

// LireCSV lit un fichier d'incidents (séparateur ;).
// Les lignes valides sont renvoyées ; chaque ligne invalide donne une erreur avec son numéro.
func LireCSV(chemin string) ([]Incident, []error, error) {
	fichier, err := os.Open(chemin)
	if err != nil {
		return nil, nil, fmt.Errorf("impossible d'ouvrir %s : %w", chemin, err)
	}
	defer fichier.Close()
	return Lire(fichier)
}

// Lire fait le travail de LireCSV à partir de n'importe quelle source (fichier, texte de test...).
func Lire(source io.Reader) ([]Incident, []error, error) {
	lecteur := csv.NewReader(source)
	lecteur.Comma = ';'
	lecteur.FieldsPerRecord = -1
	lecteur.TrimLeadingSpace = true

	entete, err := lecteur.Read()
	if err != nil {
		return nil, nil, errors.New("fichier vide ou illisible")
	}
	for i, attendu := range enTeteAttendu {
		if i >= len(entete) || strings.ToLower(strings.TrimSpace(entete[i])) != attendu {
			return nil, nil, fmt.Errorf("en-tête attendu : %s", strings.Join(enTeteAttendu, ";"))
		}
	}

	var incidents []Incident
	var erreurs []error
	vus := map[int]bool{}
	for {
		champs, err := lecteur.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		ligne, _ := lecteur.FieldPos(0)
		if err != nil {
			erreurs = append(erreurs, fmt.Errorf("ligne %d : %w", ligne, err))
			continue
		}
		inc, err := depuisChamps(champs)
		if err == nil && vus[inc.ID] {
			err = fmt.Errorf("ID %d en double", inc.ID)
		}
		if err != nil {
			erreurs = append(erreurs, fmt.Errorf("ligne %d : %w", ligne, err))
			continue
		}
		vus[inc.ID] = true
		incidents = append(incidents, inc)
	}
	return incidents, erreurs, nil
}

func depuisChamps(c []string) (Incident, error) {
	var inc Incident
	if len(c) != len(enTeteAttendu) {
		return inc, fmt.Errorf("%d colonnes au lieu de %d", len(c), len(enTeteAttendu))
	}
	var err error
	if inc.ID, err = strconv.Atoi(strings.TrimSpace(c[0])); err != nil {
		return inc, fmt.Errorf("ID %q n'est pas un nombre", c[0])
	}
	inc.Titre = strings.TrimSpace(c[1])
	if inc.Urgence, err = LireNiveau(c[2]); err != nil {
		return inc, fmt.Errorf("urgence : %w", err)
	}
	if inc.Criticite, err = LireNiveau(c[3]); err != nil {
		return inc, fmt.Errorf("criticité : %w", err)
	}
	if inc.Impact, err = strconv.Atoi(strings.TrimSpace(c[4])); err != nil {
		return inc, fmt.Errorf("impact %q n'est pas un nombre", c[4])
	}
	if inc.Date, err = time.Parse(FormatDate, strings.TrimSpace(c[5])); err != nil {
		return inc, fmt.Errorf("date %q invalide (format attendu AAAA-MM-JJ HH:MM)", c[5])
	}
	return inc, inc.Valider()
}
