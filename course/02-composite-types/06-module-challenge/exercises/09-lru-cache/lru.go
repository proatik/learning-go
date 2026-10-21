package lru

// node is one cache entry in the doubly linked recency list.
type node struct {
	key        string
	value      int
	prev, next *node
}

// LRU is a fixed-capacity least-recently-used cache.
type LRU struct {
	capacity int
	items    map[string]*node
	head     *node // most recently used
	tail     *node // least recently used
}

// NewLRU returns an empty cache holding at most capacity entries, or nil
// if capacity < 1.
func NewLRU(capacity int) *LRU {
	// TODO
	return nil
}

// Get returns the value for key and marks it most recently used.
func Get(c *LRU, key string) (int, bool) {
	// TODO
	return 0, false
}

// Put stores value under key, evicting the least recently used entry if
// the cache is full. It returns the evicted key, if any.
func Put(c *LRU, key string, value int) (evicted string, ok bool) {
	// TODO
	return "", false
}

// Delete removes key and reports whether it was present.
func Delete(c *LRU, key string) bool {
	// TODO
	return false
}

// Len returns the number of entries in the cache.
func Len(c *LRU) int {
	// TODO
	return 0
}

// Keys returns the keys from most to least recently used.
func Keys(c *LRU) []string {
	// TODO
	return nil
}
