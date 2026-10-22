# TTL Cache

A concurrent cache where every entry has its own **time to live**. Expired entries are invisible at once, and a background **janitor** goroutine deletes them every `sweepEvery` so memory doesn't grow forever. `Close` stops the janitor.

```go
type entry[V any] struct {
	val     V
	expires time.Time
}

type Cache[K comparable, V any] struct {
	mu        sync.Mutex
	items     map[K]entry[V]
	stop      chan struct{} // closed by Close to stop the janitor
	done      chan struct{} // closed by the janitor when it exits
	closeOnce sync.Once
}

func New[K comparable, V any](sweepEvery time.Duration) *Cache[K, V]

func (c *Cache[K, V]) Set(k K, v V, ttl time.Duration)
func (c *Cache[K, V]) Get(k K) (V, bool)
func (c *Cache[K, V]) Delete(k K)
func (c *Cache[K, V]) Len() int
func (c *Cache[K, V]) Close()
```

## Examples

| Code (`c := New[string, int](time.Minute)`) | Result |
|---|---|
| `c.Set("a", 1, 10*time.Second)`; `c.Get("a")` | `1, true` |
| `c.Get("a")` 10s after the `Set` | `0, false` (expired, even before a sweep) |
| `c.Len()` at that moment | `1`: the expired entry is still stored |
| `c.Len()` after the next sweep (at 1m) | `0`: the janitor deleted it |
| `c.Set("a", 2, time.Hour)` on an existing key | replaces the value **and** the expiry |
| `c.Delete("a")`; `c.Get("a")` | `0, false` |
| `c.Close()`; `c.Close()` | janitor stopped; the second call does nothing |

## Rules
- `New` starts the janitor goroutine. Every `sweepEvery` (use a `time.Ticker`) it deletes every expired entry. `sweepEvery <= 0` counts as one second.
- An entry set at time `t` with `ttl` expires at `t + ttl`; it is **expired** when `time.Now()` is at or after that moment (`!now.Before(expires)`). A `ttl <= 0` entry is expired at once.
- `Get` returns the zero value and `false` for a missing or expired key. It doesn't delete anything.
- `Len` counts every stored entry, expired or not (so you can see the janitor's work).
- `Close` stops the janitor and **waits until it has exited**. Calling it again does nothing. After `Close`, `Set`, `Get`, `Delete` and `Len` keep working; there are just no more sweeps.
- Safe for concurrent use.
