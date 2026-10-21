package cache

// node is one cache entry in the doubly linked recency list.
type node struct {
	key        string
	value      int
	prev, next *node
}

// LRU is a fixed-capacity least-recently-used cache. Create one with NewLRU.
type LRU struct {
	// OnEvict, if non-nil, is called with every entry removed because the
	// cache was full (not for Delete).
	OnEvict func(key string, value int)

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
func (c *LRU) Get(key string) (int, bool) {
	// TODO
	return 0, false
}

// Peek returns the value for key without changing the recency order.
func (c *LRU) Peek(key string) (int, bool) {
	// TODO
	return 0, false
}

// Put stores value under key as the most recently used entry, evicting the
// least recently used entry if the cache is full. It returns the evicted key.
func (c *LRU) Put(key string, value int) (evicted string, ok bool) {
	// TODO
	return "", false
}

// Delete removes key and reports whether it was present.
func (c *LRU) Delete(key string) bool {
	// TODO
	return false
}

// Len returns the number of entries. It is safe on a nil *LRU.
func (c *LRU) Len() int {
	// TODO
	return 0
}

// Keys returns the keys from most to least recently used. It is safe on a
// nil *LRU.
func (c *LRU) Keys() []string {
	// TODO
	return nil
}

// String lists the entries from most to least recently used, like
// "[c:3 a:1]". It is safe on a nil *LRU.
func (c *LRU) String() string {
	// TODO
	return ""
}
