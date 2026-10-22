package lru

// Entry is one key/value pair stored in a Cache.
type Entry[K comparable, V any] struct {
	Key K
	Val V
}

// node is an element of the cache's doubly linked recency list. It embeds
// Entry, so n.Key and n.Val are promoted fields.
type node[K comparable, V any] struct {
	Entry[K, V]
	prev, next *node[K, V] // prev is more recent, next less recent
}

// Cache is a fixed-size least-recently-used cache. Create one with New.
type Cache[K comparable, V any] struct {
	capacity int
	items    map[K]*node[K, V]
	head     *node[K, V] // most recently used
	tail     *node[K, V] // least recently used
}

// New returns an empty cache holding at most capacity entries. It panics if
// capacity < 1.
func New[K comparable, V any](capacity int) *Cache[K, V] {
	// TODO
	return &Cache[K, V]{}
}

// Get returns the value for k and marks it most recently used; ok is false
// if k is missing.
func (c *Cache[K, V]) Get(k K) (V, bool) {
	// TODO
	var zero V
	return zero, false
}

// Peek returns the value for k without changing the recency order.
func (c *Cache[K, V]) Peek(k K) (V, bool) {
	// TODO
	var zero V
	return zero, false
}

// Put stores v under k and marks it most recently used. If a new key doesn't
// fit, the least recently used entry is evicted and returned with ok true.
func (c *Cache[K, V]) Put(k K, v V) (evicted Entry[K, V], ok bool) {
	// TODO
	return Entry[K, V]{}, false
}

// Remove deletes k and reports whether it was present.
func (c *Cache[K, V]) Remove(k K) bool {
	// TODO
	return false
}

// Len returns the number of entries.
func (c *Cache[K, V]) Len() int {
	// TODO
	return 0
}

// Entries returns every entry from most to least recently used.
func (c *Cache[K, V]) Entries() []Entry[K, V] {
	// TODO
	return nil
}
