package incident

import (
	"fmt"
	"math/rand"
	"time"
)

var titresExemples = []string{
	"Vidéoprojecteur en panne", "Wi-Fi coupé", "Fuite d'eau", "Chauffage arrêté",
	"PC qui ne démarre pas", "Porte bloquée", "Imprimante hors service", "Panne d'électricité",
}

// Generer fabrique n incidents au hasard, pour les mesures de performance.
// La même graine donne toujours la même liste : les tests sont reproductibles.
func Generer(n int, graine int64) []Incident {
	hasard := rand.New(rand.NewSource(graine))
	debut := time.Date(2026, 9, 1, 8, 0, 0, 0, time.UTC)
	liste := make([]Incident, n)
	for i := range liste {
		liste[i] = Incident{
			ID:        i + 1,
			Titre:     fmt.Sprintf("%s (#%d)", titresExemples[hasard.Intn(len(titresExemples))], i+1),
			Urgence:   Niveau(hasard.Intn(4) + 1),
			Criticite: Niveau(hasard.Intn(4) + 1),
			Impact:    hasard.Intn(300),
			Date:      debut.Add(time.Duration(hasard.Intn(60*24*60)) * time.Minute),
		}
	}
	return liste
}
