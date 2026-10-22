# Heartbeat Watchdog

A worker sends heartbeats; `Watch` returns a context that is canceled **with a cause** as soon as the heartbeats stall or stop. Code using the context can then report *why* it was stopped.

```go
var (
	ErrStalled = errors.New("watchdog: no heartbeat in time")
	ErrStopped = errors.New("watchdog: heartbeats stopped")
)

func Watch(parent context.Context, beats <-chan struct{}, timeout time.Duration) (context.Context, context.CancelFunc)
```

## Examples

| Situation (`timeout` = 1s) | `ctx.Err()` | `context.Cause(ctx)` |
|---|---|---|
| a beat every 500ms | `nil` | `nil` |
| beats at 0.5s and 1.2s, then nothing | `context.Canceled` at exactly 2.2s | `ErrStalled` |
| no beat at all | `context.Canceled` at exactly 1s | `ErrStalled` |
| `beats` is closed at 300ms | `context.Canceled` at 300ms | `ErrStopped` |
| the returned `cancel()` is called | `context.Canceled` | `context.Canceled` |
| `parent` canceled with `cancel(ErrShutdown)` | `context.Canceled` | `ErrShutdown` |
| `parent` times out | `context.DeadlineExceeded` | `context.DeadlineExceeded` |

## Rules
- Return a child of `parent` created with `context.WithCancelCause`, plus a `CancelFunc` that cancels it with cause `context.Canceled` (call the cause function with `nil`).
- Start **one** goroutine that watches `beats`. Each beat restarts the `timeout` countdown.
- No beat for `timeout` → cancel with `ErrStalled`. `beats` closed → cancel with `ErrStopped`.
- The goroutine must **exit** once the context is done for any reason (parent canceled, `cancel()` called, or one of the above). The tests fail if it leaks.
- Only the first cancellation counts: a later stall or close doesn't change the cause.
