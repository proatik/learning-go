# Atomic Max

Track a running maximum from many goroutines **without a mutex**, using a `CompareAndSwap` loop. Then use it in a `Gauge` that records the peak number of tasks in progress (the "max concurrency" counter you'll use to test worker pools).

```go
func StoreMax(m *atomic.Int64, v int64) bool

type Gauge struct {
	current atomic.Int64
	peak    atomic.Int64
}

func (g *Gauge) Inc() int64
func (g *Gauge) Dec() int64
func (g *Gauge) Current() int64
func (g *Gauge) Peak() int64
```

## Examples

| Code | Result |
|---|---|
| `m` holds `5`: `StoreMax(&m, 3)` | `false`, `m` stays `5` |
| `m` holds `5`: `StoreMax(&m, 5)` | `false` (not greater) |
| `m` holds `5`: `StoreMax(&m, 6)` | `true`, `m` is `6` |
| `m` holds `-10`: `StoreMax(&m, -3)` | `true`, `m` is `-3` |
| `var g Gauge; g.Inc(); g.Inc(); g.Inc()` | returns `1`, `2`, `3`; `Peak()` = `3` |
| `g.Dec()` | `2`; `Peak()` still `3` |
| 50 goroutines `Inc()`, all wait, then `Dec()` | `Peak()` = `50`, `Current()` = `0` |

## Rules
- `StoreMax` stores `v` only if it is **strictly greater** than the current value, and returns whether it stored.
- It must be correct under contention: the stored value **never goes down**, and ends at the largest value any goroutine passed. `Load` then `Store` is not enough (another goroutine can store a bigger value in between); use a `Load` + `CompareAndSwap` loop.
- `Inc` returns the value **its own** increment produced (each of 50 overlapping `Inc` calls returns a different number from 1 to 50), and raises the peak with `StoreMax`.
- `Dec` returns the new current value. It may go below zero; the peak is never lowered.
- No mutex: only the atomic fields.
