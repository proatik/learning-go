# Retry with Backoff

Retry a failing operation with exponential backoff, but respect the context: stop the moment it's canceled, and don't even start a wait that would end past the deadline.

```go
var ErrNoTime = errors.New("retry: not enough time for another attempt")

func Retry(ctx context.Context, attempts int, base time.Duration, op func(context.Context) error) error
```

## Examples

`boom` is the error `op` returns. Times are relative to the call.

| Situation | Attempts at | Returns | After |
|---|---|---|---|
| `op` succeeds at once | 0s | `nil` | 0s |
| 5 attempts, base 1s, `op` fails twice then succeeds | 0s, 1s, 3s | `nil` | 3s |
| 3 attempts, base 100ms, always fails | 0, 100ms, 300ms | `"retry: giving up after attempt 3: boom"` | 300ms |
| base 1s, `cancel()` at 1.5s | 0s, 1s | `"retry: stopped after attempt 2: context canceled: boom"` | 1.5s |
| `ctx` already canceled | none | `context.Canceled` (exactly `ctx.Err()`) | 0s |
| 10 attempts, base 1s, `ctx` deadline 5s | 0s, 1s, 3s | `"retry: not enough time for another attempt (after attempt 3): boom"` | 3s |
| same, deadline 7s + 1ns | 0s, 1s, 3s, 7s | `ErrNoTime` error `(after attempt 4)` | 7s |

## Rules
- `attempts < 1` counts as 1. Before the first attempt, if `ctx` is already done return `ctx.Err()` without calling `op`.
- Call `op(ctx)`. Return `nil` as soon as it succeeds.
- After failed attempt `n`, in this order:
  1. `n == attempts` → return `fmt.Errorf("retry: giving up after attempt %d: %w", n, err)`.
  2. `ctx` is done → return `fmt.Errorf("retry: stopped after attempt %d: %w: %w", n, context.Cause(ctx), err)`.
  3. The wait is `base << (n-1)`: `base`, `2*base`, `4*base`, … If `ctx` has a deadline and `now + wait` is **not before** it, there's no time for another attempt: return `fmt.Errorf("%w (after attempt %d): %w", ErrNoTime, n, err)` at once, without waiting.
  4. Wait, but stop early if `ctx` is done: return the rule-2 error the moment that happens.
- `err` above is always the **last** error from `op`. The errors wrap both causes, so `errors.Is` finds `boom` and `context.Canceled` / `ErrNoTime`.
- Use a timer and `select` for the wait; never `time.Sleep`.
