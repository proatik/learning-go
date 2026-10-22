# Result Chains

Build a generic `Result[T]` — a value **or** an error — whose Go 1.27 **generic methods** let you chain steps that change the type, then collect a whole sequence of results.

```go
var ErrPanic = errors.New("result: panic")

type Result[T any] struct {
	val T
	err error
}

func Ok[T any](v T) Result[T]
func Fail[T any](err error) Result[T]
func Of[T any](v T, err error) Result[T]
func Try[T any](f func() T) (r Result[T])

func (r Result[T]) Get() (T, error)
func (r Result[T]) IsOk() bool
func (r Result[T]) OrElse(fallback T) T
func (r Result[T]) Context(msg string) Result[T]
func (r Result[T]) Map[U any](f func(T) U) Result[U]
func (r Result[T]) Then[U any](f func(T) (U, error)) Result[U]

func Collect[T any](seq iter.Seq[Result[T]]) ([]T, error)
```

## Examples

| Code | Result |
|---|---|
| `Of(strconv.Atoi("12")).Get()` | `12, nil` |
| `Ok(21).Map(double).Map(strconv.Itoa).Get()` | `"42", nil` — the type changes `int → string` |
| `Ok("x").Then(parse).Get()` | `0, strconv.Atoi: parsing "x": invalid syntax` |
| `Fail[int](err).Context("load").Context("start")` | error `"start: load: <err>"`, still `errors.Is(_, err)` |
| `Fail[int](err).OrElse(-1)` | `-1` |
| `Try(func() string { panic("x") })` | error `"result: panic: x"`, `errors.Is(_, ErrPanic)` |
| `Collect(seq of Ok(1), Ok(2), Fail(boom), Ok(4))` | `nil, "item 2: boom"` |

## Rules
- **Of(v, err)** is `Fail(err)` when `err != nil`, otherwise `Ok(v)`. A failed result's `Get` returns the **zero** `T`.
- **Map** and **Then** on a failed result return a failed `Result[U]` with the **same** error and never call `f`.
- **Then** turns `f`'s `(U, error)` into a result (hint: `Of`).
- **Context** wraps with `%w` as `"msg: err"`; a successful result is returned as is.
- **Try** recovers a panic in `f` and returns `Fail` with `fmt.Errorf("%w: %v", ErrPanic, p)`.
- **Collect** stops pulling at the **first** failure and returns `nil` plus `"item i: err"` (`i` from 0, wrapped with `%w`). On success it returns a non-nil slice, empty for an empty sequence.
