# Filter in Place

Filter a slice **without allocating**, reusing its backing array via the `s[:0]` trick:

```go
func KeepEven(s []int) []int
func RemoveValue(s []int, v int) []int
```

| Call | Result |
|---|---|
| `KeepEven([]int{1, 2, 3, 4, 5, 6})` | `[2 4 6]` |
| `KeepEven([]int{-3, -2, 0, 7})` | `[-2 0]` |
| `KeepEven([]int{1, 3})` | `[]` |
| `RemoveValue([]int{3, 1, 3, 3, 2}, 3)` | `[1 2]` |
| `RemoveValue([]int{1, 2}, 9)` | `[1 2]` |
| `RemoveValue(nil, 1)` | `[]` (len 0) |

## Rules
- Keep the original order of the kept elements.
- **Zero allocations**: the result must share `s`'s backing array (`&result[0] == &s[0]` when non-empty). Tests measure this with `testing.AllocsPerRun`.
- After filtering, **zero the leftover tail** `s[len(result):len(s)]` with the `clear` built-in, so the caller's original slice doesn't show stale values.
- Don't use the `slices` package.
