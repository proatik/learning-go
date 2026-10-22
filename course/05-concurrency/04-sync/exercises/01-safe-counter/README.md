# Safe Counter

Make a counter that many goroutines can update at once without losing a single increment, using a `sync.Mutex`.

```go
type Counter struct {
	mu sync.Mutex // guards n
	n  int
}

func (c *Counter) Inc()
func (c *Counter) Add(delta int)
func (c *Counter) Value() int
```

## Examples

| Code | Result |
|---|---|
| `var c Counter; c.Value()` | `0` (the zero value is ready to use) |
| `c.Inc(); c.Inc(); c.Value()` | `2` |
| `c.Add(10); c.Add(-15)` | `Value()` goes `12`, then `-3` |
| 50 goroutines × 1000 `c.Inc()` | `Value()` = `50000`, every time |

## Rules
- Every method locks `c.mu` while it touches `c.n`, **including `Value`** (a read racing with a write is still a data race).
- Unlock with `defer` right after `Lock`.
- `Add` accepts any `delta`: negative and zero too.
- The tests lock `c.mu` themselves and check that each method **waits** for it, so the mutex must really be used.
