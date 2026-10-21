# Sum and Largest: Array vs Slice

Write the same kind of loop for an array and for a slice:

```go
func SumArray(a [5]int) int
func Sum(s []int) int
func Largest(s []int) (largest int, ok bool)
```

| Call | Result |
|---|---|
| `SumArray([5]int{1, 2, 3, 4, 5})` | `15` |
| `SumArray([5]int{})` | `0` |
| `Sum([]int{4, -1, 7})` | `10` |
| `Sum(nil)` | `0` |
| `Largest([]int{3, 9, 2})` | `9, true` |
| `Largest([]int{-5, -2, -8})` | `-2, true` |
| `Largest(nil)` / `Largest([]int{})` | `0, false` |

## Rules
- Use `for ... range` loops.
- `Largest` must work for all-negative input: don't start from `0`, start from the first element.
- `ok` is `false` (and `largest` is `0`) only when the slice is empty.
- Don't use the built-in `max` or the `slices` package.
