# Note d'analyse — Priorisation des incidents CampusLink

## 1. Formalisation du problème

**Données.** Un incident a :
- un identifiant ;
- un titre ;
- une **urgence** de 1 à 4 (Basse, Moyenne, Haute, Critique) ;
- une **criticité** de 1 à 4 ;
- un **impact**, c'est-à-dire le nombre de personnes touchées ;
- une **date** de déclaration.

**Entrée.** Une liste de *n* incidents, lue depuis un fichier CSV ou générée au hasard.

**Sortie.** La même liste, classée du plus prioritaire au moins prioritaire. La liste de départ n'est pas modifiée.

**Règle de priorité** (fonction `PlusPrioritaire(a, b)`). On compare les critères dans l'ordre. On ne passe au critère suivant qu'en cas d'égalité :

| Ordre | Critère | Le plus prioritaire | Pourquoi |
|---|---|---|---|
| 1 | Urgence | la plus haute | ce qui doit être réglé vite passe d'abord |
| 2 | Criticité | la plus haute | à urgence égale, le plus grave d'abord |
| 3 | Impact | le plus grand | à gravité égale, celui qui gêne le plus de monde |
| 4 | Date | la plus ancienne | équité : premier arrivé, premier servi |
| 5 | ID | le plus petit | départage final : l'ordre est toujours le même (déterministe) |

Grâce au critère 5, deux incidents ne sont jamais « égaux ». Toutes les méthodes donnent donc **exactement le même ordre**, ce que les tests vérifient.

Comparer deux incidents coûte **O(1)** : au plus 5 comparaisons, quelle que soit la taille de la liste.

## 2. Les solutions

| Version | Fichier | Idée |
|---|---|---|
| **V1 : tri par insertion** | `tri/insertion.go` | on glisse chaque incident à sa place dans la partie déjà triée |
| **V2 : tri fusion** | `tri/fusion.go` | on coupe la liste en deux, on trie chaque moitié, puis on fusionne |
| **V3 : file de priorité (tas)** | `tri/file_priorite.go` | le plus prioritaire est toujours « en haut » ; on n'a pas besoin de tout trier |

## 3. Analyse de complexité (Big-O)

| Opération | Meilleur cas | Pire cas | Mémoire | Stable |
|---|---|---|---|---|
| V1 insertion : trier n incidents | **O(n)**, liste déjà triée | **O(n²)**, liste à l'envers | O(n) (copie) | oui |
| V2 fusion : trier n incidents | **O(n log n)** | **O(n log n)** | O(n) | oui |
| V3 tas : construire la file | O(n) | O(n) | O(n) | — |
| V3 tas : ajouter un incident | O(1) | O(log n) | — | — |
| V3 tas : prendre le suivant | O(log n) | O(log n) | — | — |
| V3 tas : top k | O(n + k log n) | O(n + k log n) | O(n) | — |

**Pourquoi O(n²) pour V1.** La boucle externe passe sur les n incidents. Pour chacun, la boucle interne peut reculer jusqu'au début. Au pire, cela fait 1 + 2 + … + (n−1) = n(n−1)/2 décalages, donc environ n²/2, soit **O(n²)**. Si la liste est déjà triée, la boucle interne s'arrête tout de suite : on fait **n** comparaisons, donc **O(n)**.

**Pourquoi O(n log n) pour V2.** On coupe la liste en deux jusqu'à avoir des morceaux d'un seul incident, ce qui fait **log₂(n)** niveaux. À chaque niveau, les fusions parcourent les n incidents une fois. Le total est donc **n × log₂(n)**, et c'est vrai dans tous les cas, même le pire.

**Pourquoi V3 est utile.** Au quotidien, de nouveaux incidents arrivent sans arrêt, et on veut seulement « le prochain à traiter ». Retrier toute la liste à chaque arrivée coûterait O(n log n). Avec un tas, chaque arrivée coûte **O(log n)**. Pour 100 000 incidents, log₂(100 000) ≈ 17 opérations, au lieu de 100 000 × 17.

## 4. Mesures (benchmarks)

Commande : `go test ./tri -bench . -benchmem`. Machine de test : processeur Intel Xeon 2,8 GHz. Les temps changent selon l'ordinateur, mais les **rapports** entre les tailles restent les mêmes.

| n incidents | V1 insertion | V2 fusion | V3 top 10 (tas) |
|---|---|---|---|
| 100 | 0,027 ms | 0,016 ms | 0,010 ms |
| 1 000 | 2,8 ms | 0,27 ms | 0,058 ms |
| 10 000 | 270 ms | 4,7 ms | 1,1 ms |
| 100 000 | (plusieurs dizaines de secondes, non mesuré) | 49 ms | 10 ms |

**Lecture.**
- **V1 : n × 10 donne un temps × 100.** On passe de 2,8 ms à 270 ms, soit environ ×97. C'est exactement ce que prévoit O(n²) : (10n)² = 100n².
- **V2 : n × 10 donne un temps × 10 à 17.** On passe de 0,27 ms à 4,7 ms, puis à 49 ms. C'est ce que prévoit O(n log n) : un peu plus que ×10, à cause du facteur log.
- **À 10 000 incidents, V2 est environ 57 fois plus rapide que V1.** L'écart grandit avec n.
- **Le top 10 par tas** est encore 4 à 5 fois plus rapide que le tri complet, parce qu'il ne trie pas tout.
- **Meilleur cas de V1** (10 000 incidents déjà triés) : 0,8 ms au lieu de 270 ms, ce qui confirme le O(n).

## 5. Optimisation et choix final

- **Trier toute la liste**, pour un export ou un affichage complet : **V2, le tri fusion**. Il est en O(n log n) garanti et stable.
- **Traiter les incidents au fil de l'eau**, pour le tableau de bord des techniciens : **V3, la file de priorité**. On ajoute en O(log n) et on prend le suivant en O(log n).
- **V1 n'est acceptable que pour de très petites listes**, ou des listes presque triées. On la garde comme première solution et comme point de comparaison.

## 6. Cas limites testés

Tous les cas ci-dessous sont dans `go test ./...` :

**Tri**
- liste vide ;
- un seul incident ;
- liste déjà triée ;
- liste triée à l'envers ;
- incidents tous identiques ;
- départage par la date.

**Garanties**
- la liste de départ n'est pas modifiée ;
- les trois méthodes donnent le même ordre que le tri de référence de Go, sur des listes de 10 à 5 000 incidents ;
- file vide, et k plus grand que la liste.

**Fichier CSV**
- fichier absent, vide, ou avec un mauvais en-tête ;
- titre vide ;
- niveau inconnu ou hors de 1..4 ;
- impact négatif ou qui n'est pas un nombre ;
- date mal écrite ;
- ID qui n'est pas un nombre, ou en double ;
- mauvais nombre de colonnes.

Chaque ligne invalide est ignorée avec un message qui donne son numéro. Les autres lignes sont gardées.
