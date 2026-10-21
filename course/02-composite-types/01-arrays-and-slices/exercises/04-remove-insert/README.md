# Remove and Insert by Hand

Implement the three classic slice edits with `append`, `copy` and slice expressions:

```go
func RemoveAt(s []int, i int) []int
func RemoveUnordered(s []int, i int) []int
func InsertAt(s []int, i int, v int) []int
```

| Call | Result |
|---|---|
| `RemoveAt([]int{10, 20, 30, 40}, 1)` | `[10 30 40]` |
| `RemoveAt([]int{10}, 0)` | `[]` |
| `RemoveUnordered([]int{10, 20, 30, 40}, 0)` | `[40 20 30]` (last element moved into the hole) |
| `InsertAt([]int{1, 2, 4}, 2, 3)` | `[1 2 3 4]` |
| `InsertAt([]int{1, 2}, 0, 0)` | `[0 1 2]` |
| `InsertAt([]int{1, 2}, 2, 3)` | `[1 2 3]` (at the end) |
| `InsertAt(nil, 0, 7)` | `[7]` |

## Rules
- Like `append`, each function **returns** the edited slice; the caller writes `s = RemoveAt(s, i)`.
- They may reuse (and overwrite) the backing array of `s`; don't allocate a new slice unless `InsertAt` needs more capacity (let `append` decide).
- `RemoveAt` and `InsertAt` keep the order of the other elements. `RemoveUnordered` doesn't, and runs in O(1).
- Assume valid indices: `0 <= i < len(s)` for removals, `0 <= i <= len(s)` for `InsertAt`.
- Don't use the `slices` package.
