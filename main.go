// CampusLink — priorisation des incidents.
//
// Exemples :
//
//	go run .                                  (fichier donnees/incidents.csv, tri par fusion)
//	go run . -fichier donnees/incidents.csv -algo insertion
//	go run . -top 5                           (les 5 plus prioritaires, avec la file de priorité)
package main

import (
	"flag"
	"fmt"
	"os"
	"text/tabwriter"

	"campuslink-algo/incident"
	"campuslink-algo/tri"
)

func main() {
	fichier := flag.String("fichier", "donnees/incidents.csv", "fichier CSV des incidents")
	algo := flag.String("algo", "fusion", "algorithme : insertion (v1) ou fusion (v2)")
	top := flag.Int("top", 0, "n'afficher que les N plus prioritaires (file de priorité)")
	flag.Parse()

	incidents, erreurs, err := incident.LireCSV(*fichier)
	if err != nil {
		fmt.Fprintln(os.Stderr, "erreur :", err)
		os.Exit(1)
	}
	for _, e := range erreurs {
		fmt.Fprintln(os.Stderr, "ligne ignorée :", e)
	}

	var classes []incident.Incident
	switch {
	case *top > 0:
		classes = tri.LesPlusPrioritaires(incidents, *top)
	case *algo == "insertion":
		classes = tri.ParInsertion(incidents)
	case *algo == "fusion":
		classes = tri.ParFusion(incidents)
	default:
		fmt.Fprintf(os.Stderr, "erreur : algorithme %q inconnu (insertion ou fusion)\n", *algo)
		os.Exit(2)
	}
	afficher(classes)
	fmt.Printf("\n%d incident(s) classé(s), %d ligne(s) ignorée(s)\n", len(classes), len(erreurs))
}

func afficher(incidents []incident.Incident) {
	t := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(t, "RANG\tID\tURGENCE\tCRITICITÉ\tIMPACT\tDATE\tTITRE")
	for rang, i := range incidents {
		fmt.Fprintf(t, "%d\t%d\t%s\t%s\t%d\t%s\t%s\n", rang+1, i.ID, i.Urgence, i.Criticite,
			i.Impact, i.Date.Format(incident.FormatDate), i.Titre)
	}
	t.Flush()
}
