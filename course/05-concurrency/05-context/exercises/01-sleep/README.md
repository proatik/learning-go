# Cancelable Sleep

Write `Sleep`: like `time.Sleep`, but it wakes up early when the context is done.

```go
func Sleep(ctx context.Context, d time.Duration) error
```

## Examples

| Call | Returns | After |
|---|---|---|
| `Sleep(context.Background(), time.Second)` | `nil` | 1s |
| `Sleep(ctx, time.Hour)` with `ctx` from `WithTimeout(…, 2*time.Second)` | `context.DeadlineExceeded` | 2s |
| `Sleep(ctx, time.Hour)`, `cancel()` called after 300ms | `context.Canceled` | 300ms |
| `Sleep(ctx, time.Second)` with `ctx` already canceled | `context.Canceled` | 0s |
| `Sleep(context.Background(), 0)` | `nil` | 0s |
| `Sleep(ctx, -time.Second)` with `ctx` already canceled | `context.Canceled` | 0s |

## Rules
- Return `nil` once `d` has passed, or `ctx.Err()` as soon as `ctx` is done, whichever comes first.
- If `ctx` is already done when `Sleep` is called, return `ctx.Err()` at once, even when `d <= 0`.
- With a live `ctx` and `d <= 0`, return `nil` at once.
- Don't call `time.Sleep`: it can't be interrupted. Use a timer and `select`, and stop the timer before returning.
