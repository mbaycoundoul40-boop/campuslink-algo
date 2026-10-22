package tri

import (
	"fmt"
	"testing"

	"campuslink-algo/incident"
)

// Lancer : go test ./tri -bench . -benchmem
var tailles = []int{100, 1000, 10000, 100000}

func BenchmarkParInsertion(b *testing.B) {
	for _, n := range tailles {
		if n > 10000 {
			continue // trop long en O(n²) (plusieurs minutes)
		}
		liste := incident.Generer(n, 42)
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				ParInsertion(liste)
			}
		})
	}
}

func BenchmarkParFusion(b *testing.B) {
	for _, n := range tailles {
		liste := incident.Generer(n, 42)
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				ParFusion(liste)
			}
		})
	}
}

func BenchmarkTop10ParFile(b *testing.B) {
	for _, n := range tailles {
		liste := incident.Generer(n, 42)
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				LesPlusPrioritaires(liste, 10)
			}
		})
	}
}

func BenchmarkInsertionListeDejaTriee(b *testing.B) {
	liste := ParFusion(incident.Generer(10000, 42))
	for i := 0; i < b.N; i++ {
		ParInsertion(liste) // meilleur cas : O(n)
	}
}
