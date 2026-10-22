# Generic LRU Cache

A fixed-size **least-recently-used** cache for any key and value types: a map for O(1) lookup plus a doubly linked list that keeps entries in recency order. The list node **embeds** the generic `Entry[K, V]`.

```go
type Entry[K comparable, V any] struct {
	Key K
	Val V
}

type node[K comparable, V any] struct {
	Entry[K, V]            // promotes Key and Val
	prev, next *node[K, V] // prev is more recent, next less recent
}

type Cache[K comparable, V any] struct {
	capacity int
	items    map[K]*node[K, V]
	head     *node[K, V] // most recently used
	tail     *node[K, V] // least recently used
}

func New[K comparable, V any](capacity int) *Cache[K, V]

func (c *Cache[K, V]) Get(k K) (V, bool)
func (c *Cache[K, V]) Peek(k K) (V, bool)
func (c *Cache[K, V]) Put(k K, v V) (evicted Entry[K, V], ok bool)
func (c *Cache[K, V]) Remove(k K) bool
func (c *Cache[K, V]) Len() int
func (c *Cache[K, V]) Entries() []Entry[K, V]
```

## Examples

| Code (capacity 2) | Result |
|---|---|
| `c.Put("a", 1); c.Put("b", 2)` | `Entries()` = `[{b 2} {a 1}]` (most recent first) |
| `c.Get("a")` | `1, true`; now `Entries()` = `[{a 1} {b 2}]` |
| `c.Put("c", 3)` | evicts `{b 2}`: returns `{b 2}, true` |
| `c.Peek("a")` | `1, true`, order unchanged |
| `c.Put("a", 10)` | `{"" 0}, false` (update, no eviction); `a` becomes most recent |
| `c.Get("zz")` | `0, false` |
| `c.Remove("a")` / `c.Remove("a")` | `true` / `false` |
| `New[string, int](0)` | panics: `lru: capacity must be at least 1` |

## Rules
- `New` panics with exactly `"lru: capacity must be at least 1"` if `capacity < 1`.
- `Get` and `Put` make the key the **most recently used**; `Peek`, `Len` and `Entries` don't change the order.
- `Put` on an existing key updates the value (no eviction). `Put` of a new key into a full cache first evicts the **least** recently used entry and returns it with `true`; otherwise it returns the zero `Entry` and `false`.
- A missing key: `Get` / `Peek` return the zero value of `V` and `false`; `Remove` returns `false`.
- `Entries` lists every entry from most to least recently used (empty or nil when empty).
- Every operation is O(1) except `Entries`. Keep `head`, `tail`, the `prev` / `next` links **and** the map consistent (the tests walk your list in both directions).
