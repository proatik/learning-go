# Split a Package

Move the number helpers out of package `grades` into their own package `stats`, export them, and make `Report` call them.

```go
// package stats  (stats/stats.go)
func Sum(xs []int) int
func Mean(xs []int) float64
func MinMax(xs []int) (lo, hi int, ok bool)

// package grades  (grades.go)
func Report(scores []int) string
```

## Examples

| Call | Result |
|---|---|
| `stats.Sum([]int{70, 80, 95})` | `245` |
| `stats.Mean([]int{70, 80, 95})` | `81.666…` |
| `stats.Mean(nil)` | `0` |
| `stats.MinMax([]int{3, -1, 7})` | `-1, 7, true` |
| `stats.MinMax(nil)` | `0, 0, false` |
| `grades.Report([]int{70, 80, 95})` | `"n=3 mean=81.67 min=70 max=95"` |
| `grades.Report(nil)` | `"n=0"` |

## Rules
- `grades.go` already works, using unexported helpers `sum`, `mean` and `minMax`. **Move** them into `stats/stats.go` as the exported `Sum`, `Mean` and `MinMax`.
- After the move, `grades.go` declares **only** `Report`, which calls `stats.…` (a test checks this).
- `Mean` of an empty or nil slice is `0`; `MinMax` of one returns `0, 0, false`.
- `Report` output stays exactly the same: the mean uses `%.2f`; an empty input gives `"n=0"`.
- Import the sub-package by its full path: `learngo/course/04-generics-packages-testing/04-packages/exercises/01-split-package/stats`.
