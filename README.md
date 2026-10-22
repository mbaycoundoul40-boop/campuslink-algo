# CampusLink — Priorisation des incidents (FR03 Algorithmie)

Ce programme classe les incidents du campus par ordre de priorité, selon l'urgence, la criticité, l'impact et la date. Il compare une solution simple, le **tri par insertion** en O(n²), avec une solution optimisée, le **tri fusion** en O(n log n). Il ajoute une **file de priorité**, avec un tas, pour traiter les incidents au fil de l'eau.

L'analyse complète (règle de priorité, Big-O, mesures, choix final) se trouve dans **[NOTE.md](NOTE.md)**.

## Lancer

Il faut Go 1.22 ou plus récent.

```bash
go run .                                   # classe donnees/incidents.csv (tri fusion)
go run . -algo insertion                   # même chose avec le tri par insertion (V1)
go run . -top 5                            # les 5 plus prioritaires (file de priorité)
go run . -fichier donnees/incidents-invalides.csv   # montre la gestion des erreurs
```

Exemple de sortie (`go run . -top 3`) :

```
RANG  ID  URGENCE   CRITICITÉ  IMPACT  DATE              TITRE
1     15  Critique  Critique   400     2026-10-20 10:45  Odeur de brûlé salle serveur
2     9   Critique  Critique   400     2026-10-20 11:00  Porte d'issue de secours bloquée
3     5   Critique  Critique   120     2026-10-20 10:30  Panne d'électricité amphithéâtre
```

## Tester

```bash
go test ./...                          # tous les tests (cas limites inclus)
go test -cover ./...                   # avec la couverture
go test ./tri -bench . -benchmem       # mesures de performance (V1, V2, top 10)
```

## Format du fichier CSV

Le séparateur est `;`, et la première ligne est l'en-tête :

```
id;titre;urgence;criticite;impact;date
3;Wi-Fi coupé dans tout le bâtiment A;Critique;Haute;250;2026-10-20 08:45
```

- **urgence** et **criticite** : `Basse`, `Moyenne`, `Haute`, `Critique`, ou de 1 à 4.
- **impact** : le nombre de personnes touchées (0 ou plus).
- **date** : au format `AAAA-MM-JJ HH:MM`.

## Organisation du code

```
main.go                    lecture des options, affichage du classement
incident/
  incident.go              type Incident, niveaux, validation
  priorite.go              PlusPrioritaire(a, b) : LA règle de priorité
  csv.go                   lecture du fichier, une erreur par ligne invalide
  generer.go               génération de n incidents pour les mesures
tri/
  insertion.go             V1 : tri par insertion, O(n²)
  fusion.go                V2 : tri fusion, O(n log n)
  file_priorite.go         V3 : file de priorité (container/heap)
  tri_test.go              cas limites + comparaison des 3 méthodes
  bench_test.go            benchmarks
donnees/
  incidents.csv            20 incidents valides
  incidents-invalides.csv  9 lignes invalides + 1 valide
NOTE.md                    note d'analyse Big-O et d'optimisation
```
