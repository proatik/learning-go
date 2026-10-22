# Every Tick

Call a function on every tick of a `time.Ticker` until a stop channel is closed, then stop the ticker and report how many calls were made.

```go
func Every(interval time.Duration, stop <-chan struct{}, f func(n int)) int
```

## Examples

| Call | `f` is called with | Returns | Returns after |
|---|---|---|---|
| `Every(100*time.Millisecond, stop, f)`, `stop` closed at 1050ms | `1` at 100ms, `2` at 200ms, …, `10` at 1000ms | `10` | 1050ms |
| `Every(time.Second, stop, f)`, `stop` closed at 2500ms | `1` at 1s, `2` at 2s | `2` | 2500ms |
| `Every(time.Second, stop, f)`, `stop` already closed | never | `0` | 0s |
| `Every(0, stop, f)` or a negative interval | never | `0` | 0s |

## Rules
- The first call happens one `interval` after `Every` starts (a ticker doesn't tick at time 0). `n` counts from `1`.
- Call `f` from `Every`'s own goroutine (no new goroutines), so calls never overlap.
- Return as soon as `stop` is closed, with the number of calls made. Stop the ticker before returning (`defer t.Stop()`).
- `interval <= 0` returns `0` immediately: `time.NewTicker` panics on a non-positive duration.
- The tests never close `stop` at the same instant as a tick.
