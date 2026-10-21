# Protect Middleware

A server must survive a buggy handler. Wrap handlers so a panic becomes a typed error — except a deliberate **abort**, which must keep panicking (like `net/http`'s `http.ErrAbortHandler`).

```go
var ErrAbort = errors.New("abort handler")

type PanicError struct {
	Input string // the input the handler was called with
	Value any    // the recovered panic value
}

func (e *PanicError) Error() string // panic handling "<Input>": <Value>
func (e *PanicError) Unwrap() error // Value if it is an error, else nil

func Protect(h func(string) string) func(string) (string, error)
func Pipeline(input string, steps []func(string) string) (string, error)
```

| Call | Result |
|---|---|
| `Protect(strings.ToUpper)("go")` | `"GO", nil` |
| `Protect(first)("")` where `first` does `s[:1]` | `""`, `panic handling "": runtime error: slice bounds out of range [:1] with length 0` |
| `Protect(func(string) string { panic("bad") })("x")` | `""`, `*PanicError{Input: "x", Value: "bad"}` |
| `Protect(func(string) string { panic(ErrAbort) })("x")` | panics again with `ErrAbort` |
| `Pipeline("go", []func(string) string{strings.ToUpper, reverse})` | `"OG", nil` |
| `Pipeline("x", []func(string) string{strings.ToUpper, boom})` | `""`, `step 1: panic handling "X": boom` |

## Rules
- `Protect` returns a function that calls `h(input)`:
  - no panic → `h`'s result and `nil`
  - panic → `""` and a **`*PanicError`** with the input and panic value
  - if the panic value is an `error` and `errors.Is(value, ErrAbort)` → **re-panic** with the same value (don't convert it)
- `(*PanicError).Error()` uses `%q` for the input and `%v` for the value.
- `(*PanicError).Unwrap()` returns the value if it is an `error` (so `errors.Is(err, io.EOF)` works after `panic(io.EOF)`), otherwise `nil`.
- `Pipeline` feeds each step's output into the next, running every step through `Protect`.
  - On the first error, return `""` and `fmt.Errorf("step %d: %w", i, err)` (index from 0).
  - No steps → `input, nil`.
  - An abort panic passes through `Pipeline` too.
