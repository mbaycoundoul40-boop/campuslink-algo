package tri

import "campuslink-algo/incident"

// ParFusion = VERSION 2 (solution optimisée).
//
// Principe « diviser pour régner » : on coupe la liste en deux, on trie chaque
// moitié (même méthode, récursivement), puis on fusionne les deux moitiés triées.
//
// Complexité en temps : O(n log n) dans TOUS les cas.
//   - on coupe en deux log₂(n) fois (profondeur de la récursion)
//   - à chaque niveau, la fusion parcourt les n incidents une fois
//
// Mémoire : O(n) (listes intermédiaires pendant la fusion).
// Stable : en cas d'égalité on prend d'abord l'élément de gauche.
func ParFusion(incidents []incident.Incident) []incident.Incident {
	liste := append([]incident.Incident(nil), incidents...)
	tampon := make([]incident.Incident, len(liste))
	trierFusion(liste, tampon)
	return liste
}

func trierFusion(liste, tampon []incident.Incident) {
	if len(liste) < 2 {
		return // 0 ou 1 incident : déjà trié
	}
	milieu := len(liste) / 2
	trierFusion(liste[:milieu], tampon[:milieu])
	trierFusion(liste[milieu:], tampon[milieu:])
	fusionner(liste, milieu, tampon)
}

// fusionner réunit liste[:milieu] et liste[milieu:], déjà triées chacune.
func fusionner(liste []incident.Incident, milieu int, tampon []incident.Incident) {
	copy(tampon, liste)
	gauche, droite, k := 0, milieu, 0
	for gauche < milieu && droite < len(liste) {
		// on ne prend la droite que si elle est STRICTEMENT prioritaire → stabilité
		if incident.PlusPrioritaire(tampon[droite], tampon[gauche]) {
			liste[k] = tampon[droite]
			droite++
		} else {
			liste[k] = tampon[gauche]
			gauche++
		}
		k++
	}
	k += copy(liste[k:], tampon[gauche:milieu])
	copy(liste[k:], tampon[droite:len(liste)])
}
