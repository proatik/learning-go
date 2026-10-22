# Zero Values and `Must`

Practise `var zero T` and explicit instantiation with four small helpers you will meet in real code.

```go
func Must[T any](v T, err error) T
func FirstOr[T any](xs []T, fallback T) T
func IsZero[T comparable](v T) bool
func FirstNonZero[T comparable](vals ...T) T
```

## Examples

| Call | Result |
|---|---|
| `Must(strconv.Atoi("42"))` | `42` |
| `Must(strconv.Atoi("x"))` | panics with the `*strconv.NumError` |
| `FirstOr([]string{"a", "b"}, "none")` | `"a"` |
| `FirstOr([]string{}, "none")` | `"none"` |
| `IsZero(0)`, `IsZero("")`, `IsZero(Point{})` | `true` |
| `IsZero(Point{X: 1})` | `false` |
| `FirstNonZero("", "", "guest", "admin")` | `"guest"` |
| `FirstNonZero(0, 0)` | `0` |
| `FirstNonZero[int]()` | `0` |

## Rules
- `Must(v, err)`: if `err != nil`, `panic(err)` (the error value itself). Otherwise return `v`.
  - `Must(strconv.Atoi("42"))` works because a call returning `(int, error)` can be passed straight to a function taking `(T, error)`.
- `FirstOr`: `xs[0]`, or `fallback` when `xs` is nil or empty.
- `IsZero(v)`: `true` when `v` equals the zero value of its type.
- `FirstNonZero(vals...)`: the first value that is not zero. If every value is zero (or there are none), return the zero value.
- Don't call `cmp.Or` (it does the same as `FirstNonZero`) or `reflect`.
