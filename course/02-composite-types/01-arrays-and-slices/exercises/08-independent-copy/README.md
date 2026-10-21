# Independent Copies (Beat Aliasing)

Assigning a slice copies only its header; both still share one backing array. Write functions that return results that are **fully independent** of their inputs:

```go
func Clone(s []int) []int
func WithAppended(s []int, v int) []int
func Clone2D(m [][]int) [][]int
```

| Call | Result | Independence guarantee |
|---|---|---|
| `Clone([]int{1, 2, 3})` | `[1 2 3]` | writing to the result never changes `s`, and vice versa |
| `Clone(nil)` | `nil` | |
| `Clone([]int{})` | `[]` (empty, **non-nil**) | |
| `WithAppended([]int{1, 2}, 3)` | `[1 2 3]` | even if `s` has spare capacity, `s`'s array is never written |
| `Clone2D([][]int{{1, 2}, {3}})` | `[[1 2] [3]]` | every **row** is copied too (a deep copy) |
| `Clone2D(nil)` | `nil` | |

## Rules
- Never modify the inputs.
- `Clone` keeps nil-ness: nil in → nil out; non-nil empty in → non-nil empty out.
- `WithAppended` must be safe even here:
  ```go
  base := make([]int, 2, 10)
  a := WithAppended(base, 1)
  b := WithAppended(base, 2) // must not change a
  ```
- `Clone2D` keeps nil rows nil and copies every other row with `Clone`.
- Don't use the `slices` package (in real code `slices.Clone` exists; it's next lesson).
