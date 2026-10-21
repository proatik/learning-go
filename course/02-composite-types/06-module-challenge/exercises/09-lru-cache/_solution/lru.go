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
	if capacity < 1 {
		return nil
	}
	return &LRU{capacity: capacity, items: make(map[string]*node, capacity)}
}

// unlink removes n from the list, fixing its neighbours and the ends.
func unlink(c *LRU, n *node) {
	if n.prev != nil {
		n.prev.next = n.next
	} else {
		c.head = n.next
	}
	if n.next != nil {
		n.next.prev = n.prev
	} else {
		c.tail = n.prev
	}
	n.prev, n.next = nil, nil
}

// pushFront inserts a detached n as the most recently used node.
func pushFront(c *LRU, n *node) {
	n.next = c.head
	if c.head != nil {
		c.head.prev = n
	}
	c.head = n
	if c.tail == nil { // the list was empty
		c.tail = n
	}
}

// Get returns the value for key and marks it most recently used.
func Get(c *LRU, key string) (int, bool) {
	n, ok := c.items[key]
	if !ok {
		return 0, false
	}
	unlink(c, n)
	pushFront(c, n)
	return n.value, true
}

// Put stores value under key, evicting the least recently used entry if
// the cache is full. It returns the evicted key, if any.
func Put(c *LRU, key string, value int) (evicted string, ok bool) {
	if n, found := c.items[key]; found {
		n.value = value
		unlink(c, n)
		pushFront(c, n)
		return "", false
	}
	if len(c.items) == c.capacity {
		old := c.tail
		unlink(c, old)
		delete(c.items, old.key)
		evicted, ok = old.key, true
	}
	n := &node{key: key, value: value}
	c.items[key] = n
	pushFront(c, n)
	return evicted, ok
}

// Delete removes key and reports whether it was present.
func Delete(c *LRU, key string) bool {
	n, ok := c.items[key]
	if !ok {
		return false
	}
	unlink(c, n)
	delete(c.items, key)
	return true
}

// Len returns the number of entries in the cache.
func Len(c *LRU) int {
	return len(c.items)
}

// Keys returns the keys from most to least recently used.
func Keys(c *LRU) []string {
	keys := make([]string, 0, len(c.items))
	for n := c.head; n != nil; n = n.next {
		keys = append(keys, n.key)
	}
	return keys
}
