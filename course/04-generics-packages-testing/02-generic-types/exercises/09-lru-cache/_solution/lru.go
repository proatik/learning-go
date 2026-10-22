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
	if capacity < 1 {
		panic("lru: capacity must be at least 1")
	}
	return &Cache[K, V]{capacity: capacity, items: make(map[K]*node[K, V], capacity)}
}

// Get returns the value for k and marks it most recently used; ok is false
// if k is missing.
func (c *Cache[K, V]) Get(k K) (V, bool) {
	n, ok := c.items[k]
	if !ok {
		var zero V
		return zero, false
	}
	c.moveToFront(n)
	return n.Val, true
}

// Peek returns the value for k without changing the recency order.
func (c *Cache[K, V]) Peek(k K) (V, bool) {
	n, ok := c.items[k]
	if !ok {
		var zero V
		return zero, false
	}
	return n.Val, true
}

// Put stores v under k and marks it most recently used. If a new key doesn't
// fit, the least recently used entry is evicted and returned with ok true.
func (c *Cache[K, V]) Put(k K, v V) (evicted Entry[K, V], ok bool) {
	if n, found := c.items[k]; found {
		n.Val = v
		c.moveToFront(n)
		return Entry[K, V]{}, false
	}
	if len(c.items) == c.capacity {
		lru := c.tail
		c.unlink(lru)
		delete(c.items, lru.Key)
		evicted, ok = lru.Entry, true
	}
	n := &node[K, V]{Entry: Entry[K, V]{Key: k, Val: v}}
	c.pushFront(n)
	c.items[k] = n
	return evicted, ok
}

// Remove deletes k and reports whether it was present.
func (c *Cache[K, V]) Remove(k K) bool {
	n, ok := c.items[k]
	if !ok {
		return false
	}
	c.unlink(n)
	delete(c.items, k)
	return true
}

// Len returns the number of entries.
func (c *Cache[K, V]) Len() int {
	return len(c.items)
}

// Entries returns every entry from most to least recently used.
func (c *Cache[K, V]) Entries() []Entry[K, V] {
	out := make([]Entry[K, V], 0, len(c.items))
	for n := c.head; n != nil; n = n.next {
		out = append(out, n.Entry)
	}
	return out
}

// pushFront links n in as the most recently used node.
func (c *Cache[K, V]) pushFront(n *node[K, V]) {
	n.prev, n.next = nil, c.head
	if c.head != nil {
		c.head.prev = n
	} else {
		c.tail = n
	}
	c.head = n
}

// unlink takes n out of the list, fixing its neighbours, head and tail.
func (c *Cache[K, V]) unlink(n *node[K, V]) {
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

// moveToFront marks n as the most recently used node.
func (c *Cache[K, V]) moveToFront(n *node[K, V]) {
	if c.head == n {
		return
	}
	c.unlink(n)
	c.pushFront(n)
}
