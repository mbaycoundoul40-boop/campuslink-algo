// Package tri contient les algorithmes qui classent les incidents par priorité.
// Chaque fonction renvoie une NOUVELLE liste triée : la liste de départ n'est pas modifiée.
package tri

import "campuslink-algo/incident"

// ParInsertion = VERSION 1 (solution simple).
//
// Principe : on prend les incidents un par un et on glisse chacun à sa place
// dans la partie déjà triée, comme on range des cartes dans sa main.
//
// Complexité en temps :
//   - pire cas (liste à l'envers) : O(n²), chaque incident recule jusqu'au début
//   - meilleur cas (liste déjà triée) : O(n), aucun décalage
//
// Mémoire : O(n) pour la copie (le tri lui-même n'utilise rien de plus).
// Stable : deux incidents « égaux » gardent leur ordre de départ.
func ParInsertion(incidents []incident.Incident) []incident.Incident {
	liste := append([]incident.Incident(nil), incidents...)
	for i := 1; i < len(liste); i++ {
		courant := liste[i]
		j := i - 1
		for j >= 0 && incident.PlusPrioritaire(courant, liste[j]) {
			liste[j+1] = liste[j] // on décale vers la droite
			j--
		}
		liste[j+1] = courant
	}
	return liste
}
