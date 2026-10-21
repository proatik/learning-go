# LRU Cache

Build a fixed-capacity **least-recently-used** cache with O(1) `Get` and `Put`: a map finds a node instantly, and a doubly linked list keeps nodes in recency order so the oldest can be evicted instantly.

```go
type node struct {
	key        string
	value      int
	prev, next *node
}

type LRU struct {
	capacity int
	items    map[string]*node
	head     *node // most recently used
	tail     *node // least recently used
}

func NewLRU(capacity int) *LRU
func Get(c *LRU, key string) (int, bool)
func Put(c *LRU, key string, value int) (evicted string, ok bool)
func Delete(c *LRU, key string) bool
func Len(c *LRU) int
func Keys(c *LRU) []string
```

## Example

```go
c := NewLRU(2)
Put(c, "a", 1)  // "", false
Put(c, "b", 2)  // "", false
Get(c, "a")     // 1, true   — "a" is now the most recent
Put(c, "c", 3)  // "b", true — cache full: evict the least recent
Get(c, "b")     // 0, false
Keys(c)         // [c a]     — most recent first
```

## Rules
- **NewLRU** returns `nil` if `capacity < 1`; otherwise an empty cache with its map created.
- **Get** returns the value and `true` and marks the key **most recently used**; a missing key returns `0, false` and changes nothing.
- **Put** on an existing key updates the value and marks it most recent (no eviction). On a new key: if the cache is full, first remove the **least** recently used entry and return its key with `true`; then insert the new key as most recent. Otherwise return `"", false`.
- **Delete** removes a key, returning whether it was present.
- **Len** is the number of entries; **Keys** lists them from most to least recently used (empty slice when empty).
- Use the provided fields: every entry is a `*node` in `items` **and** in the list from `head` to `tail`, with `prev`/`next` links consistent in both directions (`head.prev == nil`, `tail.next == nil`). The tests check this.
- Write small helpers such as `unlink(c, n)` and `pushFront(c, n)` — every operation is a combination of them. No loops over the list except in `Keys`.
