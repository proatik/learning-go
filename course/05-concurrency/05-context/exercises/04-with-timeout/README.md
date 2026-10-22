# Per-call Timeout

Wrap a call to another service with its **own** timeout, the way every outgoing call in a microservice should be wrapped.

```go
func Fetch(ctx context.Context, d time.Duration, f func(context.Context) (string, error)) (string, error)
```

## Examples

| Situation | Returns | After |
|---|---|---|
| `d = 1s`, `f` answers `"ok"` after 100ms | `"ok", nil` | 100ms |
| `d = 1s`, `f` would take 5s (it watches its ctx) | `"", err`: `errors.Is(err, context.DeadlineExceeded)`, `err.Error()` = `"fetch: context deadline exceeded"` | 1s |
| `d = 1s`, but `ctx` itself times out after 500ms | `"", err` with `context.DeadlineExceeded` | 500ms |
| `d = 1s`, `ctx` is canceled after 200ms | `"", err` with `context.Canceled` | 200ms |
| `f` returns `ErrNotFound` | `"", err`: `errors.Is(err, ErrNotFound)`, `"fetch: not found"` | 0s |
| `d = 0`, `f` answers `"ok"` after 1h | `"ok", nil` (no timeout of its own) | 1h |

## Rules
- Call `f` exactly once, with a **child** of `ctx` that times out after `d` (`context.WithTimeout`). Values and the parent's own deadline or cancellation still reach `f`.
- `d <= 0` means "no extra timeout": give `f` a child from `context.WithCancel` instead.
- **Always** cancel the child before `Fetch` returns, on every path (`defer cancel()`). The tests check that the context `f` received is canceled afterwards.
- On success return `f`'s result and `nil`. On failure return `""` and `f`'s error wrapped as `fmt.Errorf("fetch: %w", err)`.
