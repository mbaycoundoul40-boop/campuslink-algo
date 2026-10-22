package incident

// PlusPrioritaire dit si a doit être traité AVANT b.
//
// Règle (dans l'ordre, on passe au critère suivant seulement en cas d'égalité) :
//  1. urgence la plus haute
//  2. criticité la plus haute
//  3. impact le plus grand (plus de personnes touchées)
//  4. date la plus ancienne (premier arrivé, premier servi)
//  5. ID le plus petit (départage final : l'ordre est toujours le même)
//
// Coût : O(1), quelques comparaisons, quelle que soit la taille de la liste.
func PlusPrioritaire(a, b Incident) bool {
	if a.Urgence != b.Urgence {
		return a.Urgence > b.Urgence
	}
	if a.Criticite != b.Criticite {
		return a.Criticite > b.Criticite
	}
	if a.Impact != b.Impact {
		return a.Impact > b.Impact
	}
	if !a.Date.Equal(b.Date) {
		return a.Date.Before(b.Date)
	}
	return a.ID < b.ID
}
