# Safe Go

Run tasks concurrently so that a **panicking task can't crash the program**: recover inside each goroutine and report the panic as that task's error.

```go
func SafeRun(tasks []func()) []error
```

## Examples

| Tasks | Result |
|---|---|
| ok, `panic("boom")`, ok | `[<nil> task 1 panicked: boom <nil>]` |
| `panic(42)` | `[task 0 panicked: 42]` |
| `panic(ErrDisk)` with `ErrDisk = errors.New("disk full")` | `[task 0 panicked: disk full]`, and `errors.Is(errs[0], ErrDisk)` is `true` |
| ok, `nil` (a nil `func()`) | `[<nil> task 1 panicked: runtime error: invalid memory address or nil pointer dereference]` |
| three tasks that sleep 1s, one of them then panics | 3 results after **1s** |
| none (`nil` or empty) | an empty slice |

## Rules
- Return a slice with one entry per task, in input order: `nil` if task `i` returned normally, otherwise an error.
- Every task runs in its own goroutine. `SafeRun` returns after all of them have finished or panicked; one panic must not stop the others.
- The error message is exactly `fmt.Sprintf("task %d panicked: %v", i, r)`, where `r` is the value from `recover()`.
- If `r` is itself an `error` (for example `ErrDisk`, or a runtime error), **wrap** it so `errors.Is` / `errors.As` still find it. Other values are formatted with `%v`.
- Calling a nil task panics like any other bug; report it the same way.
