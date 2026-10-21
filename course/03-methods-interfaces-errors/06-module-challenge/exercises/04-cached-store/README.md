# Cached Key-Value Store

Design a tiny storage API: a `Store` interface, an in-memory implementation, and a **decorator** that puts a cache in front of *any* `Store` — all speaking the same sentinel errors.

```go
var (
	ErrNotFound = errors.New("kv: key not found")
	ErrEmptyKey = errors.New("kv: empty key")
)

type Store interface {
	Get(key string) (string, error)
	Set(key, value string) error
	Delete(key string) error
}

type MemoryStore struct{ data map[string]string }

func NewMemoryStore() *MemoryStore
func (m *MemoryStore) Get(key string) (string, error)
func (m *MemoryStore) Set(key, value string) error
func (m *MemoryStore) Delete(key string) error
func (m *MemoryStore) Len() int

type CachedStore struct {
	backend      Store
	cache        map[string]string
	hits, misses int
}

func NewCachedStore(backend Store) *CachedStore
func (c *CachedStore) Get(key string) (string, error)
func (c *CachedStore) Set(key, value string) error
func (c *CachedStore) Delete(key string) error
func (c *CachedStore) Stats() (hits, misses int)
func (c *CachedStore) String() string
```

## Examples

| Code | Result |
|---|---|
| `m := NewMemoryStore(); m.Set("a", "1")` | `nil` |
| `m.Get("zz")` | `"", get "zz": kv: key not found` |
| `c := NewCachedStore(m); c.Get("a"); c.Get("a")` | `"1"` twice — 1 miss, then 1 hit |
| `c.Get("zz")` | `cache: get "zz": kv: key not found` |
| `fmt.Sprint(c)` | `"cache: 1 hits, 2 misses, 1 cached"` |
| `NewCachedStore(NewCachedStore(m))` | caches stack: a `*CachedStore` is a `Store` |

## Rules
- Every method of both types returns **bare** `ErrEmptyKey` for `key == ""`, and does nothing else (no backend call, no stats).
- **MemoryStore**
  - `Get` on a missing key → `fmt.Errorf("get %q: %w", key, ErrNotFound)`.
  - `Delete` on a missing key → `fmt.Errorf("delete %q: %w", key, ErrNotFound)`.
  - `Set` overwrites; empty **values** are allowed. `Len` is the number of keys.
- **CachedStore** uses only the `Store` interface of its backend — it must work with any implementation.
  - **Get**: a key in the cache is a **hit** (no backend call). Otherwise it is a **miss**: call `backend.Get`; on success cache the value; on error return `fmt.Errorf("cache: %w", err)` and cache nothing.
  - **Set**: write-through — call `backend.Set` first. Only on success update the cache; on error return it wrapped the same way.
  - **Delete**: always remove the key from the cache, then return the backend's `Delete` error (wrapped) or `nil`.
  - **Stats** returns the hit and miss counters; **String** returns `"cache: %d hits, %d misses, %d cached"` (cached = entries in the cache map).
- Wrap with `%w` so callers can use `errors.Is(err, ErrNotFound)` through any number of layers.
