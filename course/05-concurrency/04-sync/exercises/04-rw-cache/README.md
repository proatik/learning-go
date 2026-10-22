# Read-Mostly Cache

A cache that is read far more often than written. Use a `sync.RWMutex` so readers run **in parallel** and only writers get the lock to themselves.

```go
type Cache struct {
	mu    sync.RWMutex // guards items
	items map[string]string
}

func New() *Cache

func (c *Cache) Get(key string) (value string, ok bool)
func (c *Cache) Set(key, value string)
func (c *Cache) Delete(key string)
func (c *Cache) Len() int
```

## Examples

| Code | Result |
|---|---|
| `c := New(); c.Len()` | `0` |
| `c.Get("a")` | `"", false` |
| `c.Set("a", "1"); c.Get("a")` | `"1", true` |
| `c.Set("a", "one"); c.Get("a")` | `"one", true` (overwrite; `Len` unchanged) |
| `c.Delete("a"); c.Get("a")` | `"", false` |
| `c.Delete("missing")` | nothing happens |

## Rules
- `New` returns a cache with a non-nil `items` map.
- `Get` and `Len` only read: use `RLock` / `RUnlock`.
- `Set` and `Delete` change the map: use `Lock` / `Unlock`.
- The tests hold `c.mu` themselves to check this:
  - while a **reader** holds `RLock`, `Get` and `Len` must still return, and `Set` / `Delete` must wait;
  - while a **writer** holds `Lock`, `Get` and `Len` must wait.
- The empty string is a valid key and a valid value.
