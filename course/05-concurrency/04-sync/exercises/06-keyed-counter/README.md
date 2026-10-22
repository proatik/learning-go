# Keyed Counter

Count words, status codes or user IDs from many goroutines: a `map[string]int` guarded by a mutex, with a `Snapshot` that hands out a **copy**.

```go
type Counter struct {
	mu     sync.Mutex     // guards counts
	counts map[string]int // created on first write
}

func (c *Counter) Add(key string, n int) int
func (c *Counter) Inc(key string) int
func (c *Counter) Get(key string) int
func (c *Counter) Snapshot() map[string]int
func (c *Counter) Reset() map[string]int
```

## Examples

| Code | Result |
|---|---|
| `var c Counter; c.Get("go")` | `0` |
| `c.Inc("go")`, `c.Inc("go")` | `1`, then `2` (the new count) |
| `c.Add("go", 10)` | `12` |
| `c.Add("rust", -3)` | `-3` |
| `c.Snapshot()` | `map[go:12 rust:-3]` |
| `s := c.Snapshot(); s["go"] = 100; c.Get("go")` | `12` (the copy is independent) |
| `c.Reset()` | `map[go:12 rust:-3]`; afterwards `c.Snapshot()` = `map[]` |

## Rules
- The zero value works: create the map on the first write (reading a nil map is fine, writing panics).
- `Add` and `Inc` return the key's count **right after** their own change. With 50 goroutines incrementing one key 10000 times in total, the return values are exactly 1 to 10000, each once.
- `Snapshot` returns a **new** map: changing it must not change the counter, and later counting must not change it. Never `nil`: an empty counter gives `map[string]int{}`.
- `Reset` returns the counts and empties the counter in one step: a count is never lost or returned twice, even while other goroutines are counting. Never `nil`.
- Every method holds `c.mu` while it touches `c.counts`. Returning `c.counts` itself would let callers read the map without the lock.
