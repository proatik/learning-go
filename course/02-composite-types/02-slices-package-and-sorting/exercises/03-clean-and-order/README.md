# Clean and Order

Return **new** slices built from the input: filtered, reversed, and sorted in descending order.

```go
func RemoveNegatives(nums []int) []int
func Reversed(words []string) []string
func Descending(nums []int) []int
```

| Call | Result |
|---|---|
| `RemoveNegatives([]int{3, -1, 0, -7, 5})` | `[3 0 5]` |
| `RemoveNegatives([]int{-1, -2})` | `[]` (length 0) |
| `Reversed([]string{"a", "b", "c"})` | `[c b a]` |
| `Descending([]int{2, 9, 4, 9})` | `[9 9 4 2]` |
| `Descending(nil)` | `[]` (length 0) |

## Rules
- **None of the functions may modify its argument.** Start each one from `slices.Clone`.
- `RemoveNegatives` uses `slices.DeleteFunc` (and keeps `0`, which isn't negative). Keep the order of the remaining numbers.
- `Reversed` uses `slices.Reverse`.
- `Descending` uses `slices.Sort` + `slices.Reverse`, **or** `slices.SortFunc` with a comparison that swaps its arguments. Try both.
