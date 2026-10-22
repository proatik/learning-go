# Max, Min and Clamp

Write three generic helpers that work for every ordered type: ints, floats, strings and named types of them.

```go
func MaxOf[T cmp.Ordered](xs []T) (T, bool)
func MinOf[T cmp.Ordered](xs []T) (T, bool)
func Clamp[T cmp.Ordered](v, lo, hi T) T
```

## Examples

| Call | Result |
|---|---|
| `MaxOf([]int{3, 9, 2})` | `9, true` |
| `MinOf([]float64{1.5, -2, 0})` | `-2, true` |
| `MaxOf([]string{"pear", "apple", "zoo"})` | `"zoo", true` |
| `MaxOf([]int{})` | `0, false` |
| `MinOf[string](nil)` | `"", false` |
| `Clamp(15, 0, 10)` | `10` |
| `Clamp(-3, 0, 10)` | `0` |
| `Clamp("m", "a", "f")` | `"f"` |

## Rules
- `MaxOf` / `MinOf`:
  - empty or nil slice → the zero value of `T` and `false`
  - otherwise the largest / smallest element and `true`
- `Clamp(v, lo, hi)`:
  - `v < lo` → `lo`; `v > hi` → `hi`; otherwise `v`
  - if `lo > hi`, swap them first, so `Clamp(5, 10, 0)` is `5` and `Clamp(20, 10, 0)` is `10`
- Write the loops yourself: don't call `slices.Max`, `slices.Min`, or the built-in `max` / `min`.
