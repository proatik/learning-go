# No Stutter

The `cache` package works, but every call site stutters: `cache.NewCacheStore()`, `c.CacheGet(k)`. Rename the API so it reads well with the package name, and make the zero value usable.

```go
// package cache  (cache/cache.go): the API after your changes
type Stats struct{ Hits, Misses int }
type Store struct{ /* unexported fields */ }

func New() *Store
func (s *Store) Get(key string) (string, bool)
func (s *Store) Set(key, value string)
func (s *Store) Delete(key string)
func (s *Store) Len() int
func (s *Store) Stats() Stats

// package lookup  (lookup.go)
func Cached(s *cache.Store, key string, load func(key string) (string, error)) (string, error)
```

## Renames

| Before (stutters) | After |
|---|---|
| `cache.CacheStore` | `cache.Store` |
| `cache.NewCacheStore()` | `cache.New()` |
| `cache.CacheStats` | `cache.Stats` |
| `CacheGet` / `CacheSet` / `CacheDelete` / `CacheLen` / `CacheStats` methods | `Get` / `Set` / `Delete` / `Len` / `Stats` |
| `lookup.LookupCached` | `lookup.Cached` |

## Examples

| Code | Result |
|---|---|
| `s := cache.New(); s.Set("a", "1"); s.Get("a")` | `"1", true` |
| `s.Get("zzz")` | `"", false` |
| `var s cache.Store; s.Set("k", "v"); s.Len()` | `1` (zero value works) |
| after 1 hit and 2 misses: `s.Stats()` | `Stats{Hits: 1, Misses: 2}` |
| `lookup.Cached(s, "k", load)` twice | calls `load` only the first time |

## Rules
- The tests use the **new** names, so they won't compile until every rename is done.
- Behaviour stays the same:
  - `Get` counts a hit or a miss.
  - `Delete` of a missing key does nothing.
  - `Len` is the number of keys.
- **New:** the zero value `cache.Store{}` must be ready to use. No method may panic on it, so create the map lazily in `Set`.
- `Cached`:
  - on a hit, returns the cached value;
  - on a miss, calls `load(key)`, stores the value only when `err == nil`, and returns `load`'s result.
- Update `lookup.go` to the new names too, and the doc comments: each must start with the new name.
