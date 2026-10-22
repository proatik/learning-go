package tally

import (
	"maps"
	"sync"
)

// Counter counts occurrences of string keys. It is safe for concurrent use.
// The zero value is an empty counter, ready to use.
type Counter struct {
	mu     sync.Mutex     // guards counts
	counts map[string]int // created on first write
}

// Add adds n to key's count and returns the new count.
func (c *Counter) Add(key string, n int) int {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.counts == nil {
		c.counts = make(map[string]int)
	}
	c.counts[key] += n
	return c.counts[key]
}

// Inc adds 1 to key's count and returns the new count.
func (c *Counter) Inc(key string) int {
	return c.Add(key, 1)
}

// Get returns key's count, or 0 if key was never counted.
func (c *Counter) Get(key string) int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.counts[key]
}

// Snapshot returns a copy of every count. Changing the copy doesn't affect
// the counter, and later counting doesn't change the copy.
func (c *Counter) Snapshot() map[string]int {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := maps.Clone(c.counts)
	if out == nil {
		out = map[string]int{}
	}
	return out
}

// Reset removes every count and returns what they were.
func (c *Counter) Reset() map[string]int {
	c.mu.Lock()
	defer c.mu.Unlock()
	old := c.counts
	c.counts = nil
	if old == nil {
		old = map[string]int{}
	}
	return old
}
