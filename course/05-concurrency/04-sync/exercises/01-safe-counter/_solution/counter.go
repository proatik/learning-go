package counter

import "sync"

// Counter is an integer counter that is safe for concurrent use. The zero
// value is a counter at 0, ready to use.
type Counter struct {
	mu sync.Mutex // guards n
	n  int
}

// Inc adds 1 to the counter.
func (c *Counter) Inc() {
	c.Add(1)
}

// Add adds delta (which may be negative) to the counter.
func (c *Counter) Add(delta int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.n += delta
}

// Value returns the current count.
func (c *Counter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.n
}
