# Rate Limiter

Allow at most `n` events per `interval`. Tokens live in a buffered channel; a goroutine driven by a `time.Ticker` refills it at the start of every interval.

```go
type Limiter struct {
	tokens chan struct{} // buffer of n: one value = one allowed event
	stop   chan struct{} // closed by Stop
}

func New(n int, interval time.Duration) *Limiter

func (l *Limiter) Allow() bool // take a token if one is free; never blocks
func (l *Limiter) Wait()       // block until a token is free, then take it
func (l *Limiter) Stop()       // stop the ticker and the refill goroutine
```

## Examples

| Code | Result |
|---|---|
| `l := New(3, time.Second)`; `Allow()` ×4 at 0s | `true true true false` |
| …then at 1s: `Allow()` ×4 | `true true true false` (refilled to 3) |
| …nothing used until 4s: `Allow()` ×4 | `true true true false` (unused tokens don't pile up) |
| `New(3, time.Second)`; `Wait()` ×7 in a row | they return at `0 0 0 1s 1s 1s 2s` |
| `New(4, time.Second)`; 5 goroutines × 2 `Wait()` | 4 return at 0s, 4 at 1s, 2 at 2s |
| `New(0, time.Second)` or `New(1, 0)` | panics: `ratelimit: n and interval must be positive` |

## Rules
- `New` panics with exactly `"ratelimit: n and interval must be positive"` if `n < 1` or `interval <= 0`.
- A new limiter starts **full** (`n` tokens). Every `interval` after `New`, the refill goroutine tops the bucket back up to **exactly** `n`: send `n - len(tokens)` tokens, each with a non-blocking send. (Don't "send until the buffer is full": a send to a channel with goroutines blocked in `Wait` hands the token straight to one of them, so the buffer never fills and every waiter gets through.)
- `Allow` takes a token without blocking (`select` + `default`). `Wait` blocks until it gets one. Both are safe to call from many goroutines: the channel does the synchronisation.
- `Stop` closes `stop`; the refill goroutine then stops its ticker and returns. Each test calls `Stop` once at the end; the `synctest` bubble fails with `blocked goroutines remain` if the goroutine is still running.
