package incident

import (
	"strings"
	"testing"
)

func TestLireCSVValide(t *testing.T) {
	incidents, erreurs, err := LireCSV("../donnees/incidents.csv")
	if err != nil || len(erreurs) != 0 || len(incidents) != 20 {
		t.Fatalf("attendu 20 incidents sans erreur, obtenu %d / %v / %v", len(incidents), erreurs, err)
	}
}

func TestLireCSVInvalide(t *testing.T) {
	incidents, erreurs, err := LireCSV("../donnees/incidents-invalides.csv")
	if err != nil {
		t.Fatal(err)
	}
	if len(incidents) != 1 || len(erreurs) != 9 {
		t.Errorf("attendu 1 valide et 9 erreurs, obtenu %d et %d : %v", len(incidents), len(erreurs), erreurs)
	}
}

func TestErreursDeFichier(t *testing.T) {
	if _, _, err := LireCSV("absent.csv"); err == nil {
		t.Error("fichier absent : erreur attendue")
	}
	if _, _, err := Lire(strings.NewReader("")); err == nil {
		t.Error("fichier vide : erreur attendue")
	}
	if _, _, err := Lire(strings.NewReader("a;b;c\n")); err == nil {
		t.Error("mauvais en-tête : erreur attendue")
	}
}

func TestLireNiveau(t *testing.T) {
	for texte, attendu := range map[string]Niveau{"haute": Haute, " Critique ": Critique, "1": Basse} {
		if n, err := LireNiveau(texte); err != nil || n != attendu {
			t.Errorf("LireNiveau(%q) = %v, %v", texte, n, err)
		}
	}
	if _, err := LireNiveau("5"); err == nil {
		t.Error("5 doit être refusé")
	}
}

func TestGenererReproductible(t *testing.T) {
	a, b := Generer(100, 7), Generer(100, 7)
	for i := range a {
		if a[i] != b[i] || a[i].Valider() != nil {
			t.Fatal("même graine = même liste, et incidents valides")
		}
	}
}
