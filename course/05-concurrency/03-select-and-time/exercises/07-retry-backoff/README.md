# Retry with Backoff

Retry a failing call with **exponential backoff** (wait `base`, `2·base`, `4·base`, … capped at `maxDelay`), and give up early if `stop` is closed while waiting.

```go
var (
	ErrNoAttempts = errors.New("retry: no attempts")
	ErrStopped    = errors.New("retry: stopped")
)

func Do(stop <-chan struct{}, attempts int, base, maxDelay time.Duration, f func() error) error
```

## Examples

`base = 100ms`, `maxDelay = 1s` unless noted; `f` returns `fail 1`, `fail 2`, … on its 1st, 2nd, … call.

| Call | `f` called at | Returns | Returns after |
|---|---|---|---|
| `attempts = 4`, `f` always fails | 0, 100ms, 300ms, 700ms | `retry: 4 attempts failed: fail 4` | 700ms |
| `attempts = 5`, `f` succeeds on call 3 | 0, 100ms, 300ms | `nil` | 300ms |
| `attempts = 6`, `maxDelay = 300ms`, always fails | 0, 100, 300, 600, 900, 1200ms | `retry: 6 attempts failed: fail 6` | 1200ms |
| `attempts = 1`, always fails | 0 | `retry: 1 attempts failed: fail 1` | 0s |
| `attempts = 0` | never | `ErrNoAttempts` | 0s |
| `attempts = 5`, `stop` closed at 250ms | 0, 100ms | `retry: stopped after 2 attempts: fail 2` | 250ms |

## Rules
- Call `f` up to `attempts` times; return `nil` at the first success.
- The wait after failed call `k` (1, 2, 3, …) is `base·2^(k-1)`, but never more than `maxDelay`. There's **no** wait after the last attempt.
- All attempts failed: return `fmt.Errorf("retry: %d attempts failed: %w", attempts, lastErr)`.
- `attempts < 1`: don't call `f`; return `ErrNoAttempts`.
- Wait with a `select` on a timer and `stop`. If `stop` closes during a wait, return at once with `fmt.Errorf("%w after %d attempts: %w", ErrStopped, calls, lastErr)`. Stop the timer in that case. A `nil` `stop` never stops.
- The tests use `errors.Is` for the sentinels and for `lastErr`, and compare the full message.
