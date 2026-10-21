package cache

import (
	"fmt"
	"strings"
)

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
	if capacity < 1 {
		return nil
	}
	return &LRU{capacity: capacity, items: make(map[string]*node, capacity)}
}

// unlink removes n from the list, fixing its neighbours and the ends.
func (c *LRU) unlink(n *node) {
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
func (c *LRU) pushFront(n *node) {
	n.next = c.head
	if c.head != nil {
		c.head.prev = n
	}
	c.head = n
	if c.tail == nil {
		c.tail = n
	}
}

// Get returns the value for key and marks it most recently used.
func (c *LRU) Get(key string) (int, bool) {
	n, ok := c.items[key]
	if !ok {
		return 0, false
	}
	c.unlink(n)
	c.pushFront(n)
	return n.value, true
}

// Peek returns the value for key without changing the recency order.
func (c *LRU) Peek(key string) (int, bool) {
	n, ok := c.items[key]
	if !ok {
		return 0, false
	}
	return n.value, true
}

// Put stores value under key as the most recently used entry, evicting the
// least recently used entry if the cache is full. It returns the evicted key.
func (c *LRU) Put(key string, value int) (evicted string, ok bool) {
	if n, found := c.items[key]; found {
		n.value = value
		c.unlink(n)
		c.pushFront(n)
		return "", false
	}
	if len(c.items) == c.capacity {
		old := c.tail
		c.unlink(old)
		delete(c.items, old.key)
		if c.OnEvict != nil {
			c.OnEvict(old.key, old.value)
		}
		evicted, ok = old.key, true
	}
	n := &node{key: key, value: value}
	c.items[key] = n
	c.pushFront(n)
	return evicted, ok
}

// Delete removes key and reports whether it was present.
func (c *LRU) Delete(key string) bool {
	n, ok := c.items[key]
	if !ok {
		return false
	}
	c.unlink(n)
	delete(c.items, key)
	return true
}

// Len returns the number of entries. It is safe on a nil *LRU.
func (c *LRU) Len() int {
	if c == nil {
		return 0
	}
	return len(c.items)
}

// Keys returns the keys from most to least recently used. It is safe on a
// nil *LRU.
func (c *LRU) Keys() []string {
	if c == nil {
		return []string{}
	}
	keys := make([]string, 0, len(c.items))
	for n := c.head; n != nil; n = n.next {
		keys = append(keys, n.key)
	}
	return keys
}

// String lists the entries from most to least recently used, like
// "[c:3 a:1]". It is safe on a nil *LRU.
func (c *LRU) String() string {
	if c == nil {
		return "[]"
	}
	var b strings.Builder
	b.WriteByte('[')
	for n := c.head; n != nil; n = n.next {
		if n != c.head {
			b.WriteByte(' ')
		}
		fmt.Fprintf(&b, "%s:%d", n.key, n.value)
	}
	b.WriteByte(']')
	return b.String()
}
