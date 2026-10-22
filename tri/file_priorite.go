package tri

import (
	"container/heap"

	"campuslink-algo/incident"
)

// FilePriorite = AMÉLIORATION pour le travail au quotidien.
//
// Les techniciens n'ont pas besoin de toute la liste triée : ils veulent
// « le prochain incident à traiter », et de nouveaux incidents arrivent sans arrêt.
// Retrier toute la liste à chaque arrivée coûterait O(n log n).
// Avec un tas (heap), le plus prioritaire est toujours en haut :
//   - Ajouter  : O(log n)
//   - Suivant  : O(log n) (retire et renvoie le plus prioritaire)
//   - Regarder : O(1)     (voir le plus prioritaire sans le retirer)
//   - Construire depuis une liste de n incidents : O(n)
type FilePriorite struct {
	tas tasIncidents
}

// NouvelleFile construit la file en O(n) (heap.Init), plus rapide que n ajouts (O(n log n)).
func NouvelleFile(incidents []incident.Incident) *FilePriorite {
	f := &FilePriorite{tas: append(tasIncidents(nil), incidents...)}
	heap.Init(&f.tas)
	return f
}

func (f *FilePriorite) Ajouter(i incident.Incident) { heap.Push(&f.tas, i) }
func (f *FilePriorite) Taille() int                 { return f.tas.Len() }

// Suivant retire et renvoie l'incident le plus prioritaire (ok = false si la file est vide).
func (f *FilePriorite) Suivant() (incident.Incident, bool) {
	if f.tas.Len() == 0 {
		return incident.Incident{}, false
	}
	return heap.Pop(&f.tas).(incident.Incident), true
}

// Regarder renvoie le plus prioritaire sans le retirer.
func (f *FilePriorite) Regarder() (incident.Incident, bool) {
	if f.tas.Len() == 0 {
		return incident.Incident{}, false
	}
	return f.tas[0], true
}

// LesPlusPrioritaires renvoie les k premiers incidents en O(n + k log n),
// sans trier toute la liste (utile pour un tableau de bord « Top 10 »).
func LesPlusPrioritaires(incidents []incident.Incident, k int) []incident.Incident {
	f := NouvelleFile(incidents)
	var top []incident.Incident
	for len(top) < k {
		i, ok := f.Suivant()
		if !ok {
			break
		}
		top = append(top, i)
	}
	return top
}

// tasIncidents implémente heap.Interface (Len, Less, Swap, Push, Pop).
type tasIncidents []incident.Incident

func (t tasIncidents) Len() int           { return len(t) }
func (t tasIncidents) Less(i, j int) bool { return incident.PlusPrioritaire(t[i], t[j]) }
func (t tasIncidents) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t *tasIncidents) Push(x any)        { *t = append(*t, x.(incident.Incident)) }
func (t *tasIncidents) Pop() any {
	ancien := *t
	dernier := ancien[len(ancien)-1]
	*t = ancien[:len(ancien)-1]
	return dernier
}
