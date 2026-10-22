package tri

import (
	"reflect"
	"sort"
	"testing"
	"time"

	"campuslink-algo/incident"
)

var date = time.Date(2026, 10, 20, 8, 0, 0, 0, time.UTC)

func inc(id int, u, c incident.Niveau, impact int, decalage time.Duration) incident.Incident {
	return incident.Incident{ID: id, Titre: "test", Urgence: u, Criticite: c, Impact: impact, Date: date.Add(decalage)}
}

// reference = tri de la bibliothèque standard, qui sert d'arbitre.
func reference(l []incident.Incident) []incident.Incident {
	copie := append([]incident.Incident(nil), l...)
	sort.SliceStable(copie, func(i, j int) bool { return incident.PlusPrioritaire(copie[i], copie[j]) })
	return copie
}

func ids(l []incident.Incident) []int {
	r := []int{}
	for _, i := range l {
		r = append(r, i.ID)
	}
	return r
}

// parFile trie avec la file de priorité (on retire tout, dans l'ordre).
func parFile(l []incident.Incident) []incident.Incident { return LesPlusPrioritaires(l, len(l)) }

var algorithmes = map[string]func([]incident.Incident) []incident.Incident{
	"insertion": ParInsertion, "fusion": ParFusion, "file": parFile,
}

func TestCasLimites(t *testing.T) {
	cas := map[string][]incident.Incident{
		"liste vide":     {},
		"un seul":        {inc(1, incident.Haute, incident.Haute, 5, 0)},
		"déjà triée":     {inc(1, incident.Critique, incident.Critique, 9, 0), inc(2, incident.Haute, incident.Haute, 5, 0), inc(3, incident.Basse, incident.Basse, 1, 0)},
		"à l'envers":     {inc(3, incident.Basse, incident.Basse, 1, 0), inc(2, incident.Haute, incident.Haute, 5, 0), inc(1, incident.Critique, incident.Critique, 9, 0)},
		"tout identique": {inc(4, incident.Moyenne, incident.Moyenne, 3, 0), inc(2, incident.Moyenne, incident.Moyenne, 3, 0), inc(3, incident.Moyenne, incident.Moyenne, 3, 0)},
		"départage date": {inc(1, incident.Haute, incident.Haute, 5, time.Hour), inc(2, incident.Haute, incident.Haute, 5, 0)},
	}
	for nomCas, liste := range cas {
		attendu := ids(reference(liste))
		for nomAlgo, trier := range algorithmes {
			if obtenu := ids(trier(liste)); !reflect.DeepEqual(obtenu, attendu) {
				t.Errorf("%s / %s : %v, attendu %v", nomCas, nomAlgo, obtenu, attendu)
			}
		}
	}
}

func TestRegleDePriorite(t *testing.T) {
	a := inc(1, incident.Critique, incident.Basse, 0, 0)
	b := inc(2, incident.Haute, incident.Critique, 500, 0)
	if !incident.PlusPrioritaire(a, b) {
		t.Error("l'urgence doit passer avant la criticité et l'impact")
	}
	ancien := inc(3, incident.Haute, incident.Haute, 5, -time.Hour)
	recent := inc(4, incident.Haute, incident.Haute, 5, 0)
	if !incident.PlusPrioritaire(ancien, recent) {
		t.Error("à égalité, le plus ancien passe en premier")
	}
}

func TestListeDeDepartNonModifiee(t *testing.T) {
	liste := []incident.Incident{inc(1, incident.Basse, incident.Basse, 1, 0), inc(2, incident.Critique, incident.Critique, 9, 0)}
	for nom, trier := range algorithmes {
		trier(liste)
		if liste[0].ID != 1 {
			t.Errorf("%s a modifié la liste de départ", nom)
		}
	}
}

func TestGrandesListesIdentiques(t *testing.T) {
	for _, n := range []int{10, 100, 1000, 5000} {
		liste := incident.Generer(n, int64(n))
		attendu := ids(reference(liste))
		for nom, trier := range algorithmes {
			if !reflect.DeepEqual(ids(trier(liste)), attendu) {
				t.Errorf("n=%d : %s ne donne pas le même ordre que la référence", n, nom)
			}
		}
	}
}

func TestFilePriorite(t *testing.T) {
	f := NouvelleFile(nil)
	if _, ok := f.Suivant(); ok {
		t.Fatal("une file vide ne doit rien renvoyer")
	}
	f.Ajouter(inc(1, incident.Basse, incident.Basse, 1, 0))
	f.Ajouter(inc(2, incident.Critique, incident.Haute, 1, 0))
	f.Ajouter(inc(3, incident.Haute, incident.Haute, 1, 0))
	if haut, _ := f.Regarder(); haut.ID != 2 || f.Taille() != 3 {
		t.Fatalf("Regarder doit montrer l'incident 2 sans le retirer")
	}
	var ordre []int
	for f.Taille() > 0 {
		i, _ := f.Suivant()
		ordre = append(ordre, i.ID)
	}
	if !reflect.DeepEqual(ordre, []int{2, 3, 1}) {
		t.Errorf("ordre %v, attendu [2 3 1]", ordre)
	}
	if top := LesPlusPrioritaires(incident.Generer(50, 1), 100); len(top) != 50 {
		t.Errorf("k plus grand que la liste : %d incidents au lieu de 50", len(top))
	}
}
