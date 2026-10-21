# Merge Intervals

An interval is a `[2]int{start, end}` on the number line, with `start <= end` (you may assume every input interval is valid). Sort first, then a single pass merges everything.

```go
func Merge(intervals [][2]int) [][2]int
func Covered(intervals [][2]int) int
func Gaps(intervals [][2]int) [][2]int
```

| Function | Returns |
|---|---|
| `Merge` | the union as non-overlapping intervals, sorted by start. Intervals that **overlap or touch** (`next.start <= current.end`) merge into one |
| `Covered` | the total length of the union: the sum of `end - start` over the merged intervals |
| `Gaps` | the holes between the merged intervals, in order |

| Input | `Merge` | `Covered` | `Gaps` |
|---|---|---|---|
| `[[1 3] [8 10] [2 6]]` | `[[1 6] [8 10]]` | `7` | `[[6 8]]` |
| `[[1 3] [3 5]]` | `[[1 5]]` (touching) | `4` | `[]` |
| `[[5 7] [1 2] [3 4]]` | `[[1 2] [3 4] [5 7]]` | `4` | `[[2 3] [4 5]]` |
| `[[1 10] [2 3] [4 5]]` | `[[1 10]]` (contained) | `9` | `[]` |
| `[[4 4]]` | `[[4 4]]` | `0` | `[]` |
| `[]` | `[]` | `0` | `[]` |

## Rules
- **Do not modify** `intervals`. Tip: `[2]int` is an **array** (a value), so `slices.Clone` copies each interval completely; sorting the clone can't affect the caller.
- Sort with `slices.SortFunc` by start, then by end (`cmp.Or`).
- After sorting, one loop: either extend the last merged interval (`end = max(end, next.end)`) or append a new one.
- `Covered` and `Gaps` build on `Merge`.
