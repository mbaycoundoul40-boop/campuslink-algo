// Package incident décrit un incident CampusLink et la règle qui dit
// lequel de deux incidents doit être traité en premier.
package incident

import (
	"fmt"
	"strings"
	"time"
)

// Niveau = urgence ou criticité, de 1 (Basse) à 4 (Critique).
type Niveau int

const (
	Basse Niveau = iota + 1
	Moyenne
	Haute
	Critique
)

var nomsNiveaux = map[Niveau]string{Basse: "Basse", Moyenne: "Moyenne", Haute: "Haute", Critique: "Critique"}

func (n Niveau) String() string {
	if nom, ok := nomsNiveaux[n]; ok {
		return nom
	}
	return fmt.Sprintf("Niveau(%d)", int(n))
}

// LireNiveau accepte le mot ("haute", "Critique") ou le chiffre ("3").
func LireNiveau(texte string) (Niveau, error) {
	texte = strings.ToLower(strings.TrimSpace(texte))
	for n, nom := range nomsNiveaux {
		if texte == strings.ToLower(nom) || texte == fmt.Sprint(int(n)) {
			return n, nil
		}
	}
	return 0, fmt.Errorf("niveau %q inconnu (Basse, Moyenne, Haute, Critique ou 1 à 4)", texte)
}

// Incident = un problème signalé sur le campus.
type Incident struct {
	ID        int
	Titre     string
	Urgence   Niveau    // à quelle vitesse il faut agir
	Criticite Niveau    // gravité des conséquences
	Impact    int       // nombre de personnes touchées
	Date      time.Time // date et heure de déclaration
}

// Valider vérifie qu'un incident respecte les règles.
func (i Incident) Valider() error {
	switch {
	case i.ID <= 0:
		return fmt.Errorf("ID %d invalide (doit être > 0)", i.ID)
	case strings.TrimSpace(i.Titre) == "":
		return fmt.Errorf("incident %d : titre vide", i.ID)
	case i.Urgence < Basse || i.Urgence > Critique:
		return fmt.Errorf("incident %d : urgence %d hors de 1..4", i.ID, i.Urgence)
	case i.Criticite < Basse || i.Criticite > Critique:
		return fmt.Errorf("incident %d : criticité %d hors de 1..4", i.ID, i.Criticite)
	case i.Impact < 0:
		return fmt.Errorf("incident %d : impact négatif", i.ID)
	case i.Date.IsZero():
		return fmt.Errorf("incident %d : date manquante", i.ID)
	}
	return nil
}
