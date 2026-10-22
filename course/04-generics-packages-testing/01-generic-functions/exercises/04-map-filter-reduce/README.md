# Map, Filter, Reduce, Find

Build the four classic slice helpers. `Map` and `Reduce` need **two** type parameters.

```go
func Map[T, U any](xs []T, f func(T) U) []U
func Filter[T any](xs []T, keep func(T) bool) []T
func Reduce[T, A any](xs []T, init A, f func(A, T) A) A
func Find[T any](xs []T, pred func(T) bool) (T, bool)
```

## Examples

| Call | Result |
|---|---|
| `Map([]int{1, 2, 3}, strconv.Itoa)` | `[]string{"1", "2", "3"}` |
| `Map([]string{"go", "gopher"}, func(s string) int { return len(s) })` | `[]int{2, 6}` |
| `Filter([]int{1, 2, 3, 4}, isEven)` | `[]int{2, 4}` |
| `Reduce([]int{1, 2, 3}, 0, add)` | `6` |
| `Reduce([]string{"a", "bb"}, 0, func(n int, s string) int { return n + len(s) })` | `3` |
| `Find([]int{1, 4, 6}, isEven)` | `4, true` |
| `Find([]int{1, 3}, isEven)` | `0, false` |

## Rules
- `Map`: a new slice with `f(x)` for every element, in order. Its length is always `len(xs)`.
- `Filter`: a new slice with the elements for which `keep` returns `true`, in order. When nothing is kept, return a slice of length `0` (nil is fine).
- `Reduce`: start with `acc := init`, then `acc = f(acc, x)` for every element, left to right. An empty slice returns `init`.
- `Find`: the **first** element where `pred` is `true`, and `true`. If there is none: the zero value of `T` and `false`.
- Call each function argument **exactly once per element** (`Find` stops at the first match).
- Never modify `xs`.
