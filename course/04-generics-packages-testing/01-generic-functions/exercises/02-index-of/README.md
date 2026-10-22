# Find Things with `comparable`

Write search helpers that work for any type supporting `==`: numbers, strings, structs, pointers.

```go
func IndexOf[T comparable](xs []T, target T) int
func LastIndexOf[T comparable](xs []T, target T) int
func Contains[T comparable](xs []T, target T) bool
func Count[T comparable](xs []T, target T) int
```

## Examples

| Call | Result |
|---|---|
| `IndexOf([]string{"a", "b", "a"}, "a")` | `0` |
| `LastIndexOf([]string{"a", "b", "a"}, "a")` | `2` |
| `IndexOf([]int{1, 2, 3}, 9)` | `-1` |
| `Contains([]int{1, 2, 3}, 2)` | `true` |
| `Count([]int{1, 2, 1, 1}, 1)` | `3` |
| `Count([]Point{{1, 2}, {1, 2}}, Point{1, 2})` | `2` |
| `IndexOf([]int(nil), 0)` | `-1` |

## Rules
- `IndexOf` / `LastIndexOf`: the first / last index `i` with `xs[i] == target`, or `-1` if there is none.
- `Contains`: `true` if `target` is in `xs`. Build it on top of `IndexOf`.
- `Count`: how many elements equal `target` (`0` for nil or empty).
- Don't call `slices.Index` or `slices.Contains`: write the loops.
