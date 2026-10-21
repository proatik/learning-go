# Try

Build a general-purpose **panic-to-error** converter and use it to run a batch of jobs where one crash doesn't stop the others.

```go
func Try(f func()) (err error)
func TryInt(f func() int) (n int, err error)
func TryAll(jobs []func()) []error
```

| Call | Result |
|---|---|
| `Try(func() {})` | `nil` |
| `Try(func() { panic("boom") })` | `recovered: boom` |
| `Try(func() { panic(42) })` | `recovered: 42` |
| `Try(func() { panic(io.EOF) })` | `recovered: EOF`, and `errors.Is(err, io.EOF)` is true |
| `Try(func() { panic(nil) })` | `recovered: runtime error: panic called with nil argument` |
| `Try(nil)` | `recovered: runtime error: invalid memory address or nil pointer dereference` |
| `TryInt(func() int { return 7 })` | `7, nil` |
| `TryInt(func() int { panic("x") })` | `0, recovered: x` |
| `TryAll([]func(){ok, boom, ok})` | `[nil, recovered: boom, nil]` |

## Rules
- Error text is always `recovered: <panic value>`.
- If the panic value is an `error`, wrap it with **`%w`** (so `errors.Is` / `errors.As` see it, including `*runtime.PanicNilError` and other `runtime.Error`s); otherwise format it with `%v`.
- `TryInt` returns `0` on panic.
- `TryAll` runs **every** job in order and returns one entry per job (`nil` for a job that didn't panic). `TryAll(nil)` returns an empty (or nil) slice.
- None of the functions may let a panic escape.
