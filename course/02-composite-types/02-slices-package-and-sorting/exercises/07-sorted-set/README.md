# A Sorted Set of ints

Keep a slice **always sorted**, so lookups can use binary search (O(log n)) instead of a linear scan.

```go
func Has(set []int, v int) bool
func InsertSorted(s []int, v int) []int
func Add(set []int, v int) ([]int, bool)
func Remove(set []int, v int) ([]int, bool)
```

All inputs are sorted ascending. Every function keeps them sorted.

| Function | Does |
|---|---|
| `Has` | reports whether `v` is in `set` |
| `InsertSorted` | inserts `v` at its sorted position, **even if it is already there** (duplicates allowed) |
| `Add` | inserts `v` only if absent; the `bool` says whether it was added |
| `Remove` | removes **one** occurrence of `v`; the `bool` says whether one was found |

| Call | Result |
|---|---|
| `Has([]int{1, 3, 5}, 3)` | `true` |
| `InsertSorted([]int{1, 3, 5}, 4)` | `[1 3 4 5]` |
| `InsertSorted([]int{1, 3, 5}, 3)` | `[1 3 3 5]` |
| `InsertSorted(nil, 7)` | `[7]` |
| `Add([]int{1, 3, 5}, 0)` | `[0 1 3 5], true` |
| `Add([]int{1, 3, 5}, 5)` | `[1 3 5], false` |
| `Remove([]int{1, 3, 5}, 3)` | `[1 5], true` |
| `Remove([]int{1, 3, 5}, 2)` | `[1 3 5], false` |

## Rules
- Use `slices.BinarySearch` to find the position; **no** `for` loops and **no** re-sorting.
- Use `slices.Insert` and `slices.Delete`. Like `append`, they may reuse the input's backing array, so these functions return the new slice and callers must use it: `set, _ = Add(set, 4)`.
