# Run All

Run every function in a slice in its **own goroutine** and return only when all of them have finished.

```go
func RunAll(tasks []func())
```

## Examples

| Call | Result |
|---|---|
| `RunAll(three tasks that each sleep 1s)` | returns after **1s** (not 3s), all three ran |
| `RunAll(tasks sleeping 1s, 3s, 2s)` | returns after 3s, all three ran |
| `RunAll(nil)` / `RunAll([]func(){})` | returns immediately |

## Rules
- Start one goroutine per task with `go`, and use a `sync.WaitGroup` with `Add` / `Done` / `Wait`.
- `RunAll` must not return until **every** task has returned.
- Each task runs exactly once.
