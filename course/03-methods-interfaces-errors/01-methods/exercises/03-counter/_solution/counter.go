package counter

// Counter counts up and down but never below zero. The zero value is a
// counter at 0, ready to use.
type Counter struct {
	n int
}

// Inc adds one.
func (c *Counter) Inc() {
	c.n++
}

// Dec subtracts one, but never goes below zero.
func (c *Counter) Dec() {
	c.Add(-1)
}

// Add adds delta (which may be negative), clamping the result at zero.
func (c *Counter) Add(delta int) {
	c.n = max(c.n+delta, 0)
}

// Reset sets the counter back to zero.
func (c *Counter) Reset() {
	c.n = 0
}

// Value returns the current count.
func (c *Counter) Value() int {
	return c.n
}

// CountWords returns one counter per distinct word, counting occurrences.
func CountWords(words []string) map[string]*Counter {
	counts := make(map[string]*Counter)
	for _, w := range words {
		c, ok := counts[w]
		if !ok {
			c = &Counter{}
			counts[w] = c
		}
		c.Inc()
	}
	return counts
}
