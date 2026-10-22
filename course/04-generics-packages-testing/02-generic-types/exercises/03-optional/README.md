# Optional Values

Build `Optional[T]`: a value that is either **present** (`Some(v)`) or **absent** (`None`). It makes "maybe missing" explicit, even when the zero value (`0`, `""`) is a valid value.

```go
type Optional[T any] struct {
	value T
	ok    bool
}

func Some[T any](v T) Optional[T]
func None[T any]() Optional[T]

func (o Optional[T]) Get() (T, bool)
func (o Optional[T]) IsSome() bool
func (o Optional[T]) OrElse(fallback T) T
func (o Optional[T]) String() string

func First[T any](s []T) Optional[T]
func Lookup[K comparable, V any](m map[K]V, k K) Optional[V]
```

## Examples

| Code | Result |
|---|---|
| `Some(5).Get()` | `5, true` |
| `Some(0).IsSome()` | `true` (zero is still a value) |
| `None[int]().Get()` | `0, false` |
| `var o Optional[string]; o.IsSome()` | `false` (the zero value is `None`) |
| `None[string]().OrElse("guest")` | `"guest"` |
| `Some("ann").OrElse("guest")` | `"ann"` |
| `Some(42).String()` / `None[int]().String()` | `"Some(42)"` / `"None"` |
| `First([]string{"a", "b"})` | `Some(a)` |
| `First([]int{})` | `None` |
| `Lookup(map[string]int{"x": 0}, "x")` | `Some(0)` |
| `Lookup(map[string]int{"x": 0}, "y")` | `None` |

## Rules
- The zero value of `Optional[T]` is `None`.
- `Get` on `None` returns the zero value of `T` and `false`.
- `String` formats a present value with `%v`: `"Some(<v>)"`; an absent one is `"None"`.
- `None` has no argument to infer `T` from, so callers write `None[int]()`.
- `First` returns `None` for a nil or empty slice. `Lookup` works on a nil map.
