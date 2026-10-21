# LRU Cache, Method Edition

Refactor Module 2's LRU cache into a proper type: the same O(1) map + doubly linked list, but behaviour lives in **methods**, internals are unexported, it prints itself, and it reports evictions through a callback.

```go
type node struct {
	key        string
	value      int
	prev, next *node
}

type LRU struct {
	OnEvict func(key string, value int) // optional eviction callback

	capacity int
	items    map[string]*node
	head     *node // most recently used
	tail     *node // least recently used
}

func NewLRU(capacity int) *LRU
func (c *LRU) Get(key string) (int, bool)
func (c *LRU) Peek(key string) (int, bool)
func (c *LRU) Put(key string, value int) (evicted string, ok bool)
func (c *LRU) Delete(key string) bool
func (c *LRU) Len() int
func (c *LRU) Keys() []string
func (c *LRU) String() string
```

## Examples

| Code | Result |
|---|---|
| `c := NewLRU(2); c.Put("a", 1); c.Put("b", 2)` | `"", false` each |
| `c.Get("a")` | `1, true` — `a` is now most recent |
| `c.Put("c", 3)` | `"b", true` — full: evict the least recent |
| `c.Keys()` / `fmt.Sprint(c)` | `[c a]` / `"[c:3 a:1]"` |
| `c.Peek("a")` | `1, true` — order unchanged |
| `var n *LRU; n.Len()` | `0` (no panic) |

## Rules
- **NewLRU** returns `nil` if `capacity < 1`; otherwise an empty cache whose `items` map is created.
- **Get** returns the value and marks the key most recently used; a missing key returns `0, false`.
- **Peek** is like `Get` but **never** changes the order.
- **Put** on an existing key updates the value and makes it most recent (no eviction). On a new key with a full cache, first remove the **least** recently used entry, call `OnEvict(key, value)` with it if `OnEvict != nil`, and return its key with `true`. Otherwise return `"", false`.
- **Delete** removes a key and reports whether it was present. It does **not** call `OnEvict`.
- **Keys** lists keys from most to least recent — an empty, non-nil slice when empty.
- **String** returns `"[k1:v1 k2:v2 ...]"`, most recent first; `"[]"` when empty.
- **Len**, **Keys** and **String** must work on a **nil** `*LRU` (`0`, `[]string{}`, `"[]"`).
- Use **pointer receivers** everywhere, and unexported helper methods such as `c.unlink(n)` and `c.pushFront(n)`. The tests check that the list links stay consistent.
