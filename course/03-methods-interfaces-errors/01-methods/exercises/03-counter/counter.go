package counter

// Counter counts up and down but never below zero. The zero value is a
// counter at 0, ready to use.
type Counter struct {
	n int
}

// Inc adds one.
func (c *Counter) Inc() {
	// TODO
}

// Dec subtracts one, but never goes below zero.
func (c *Counter) Dec() {
	// TODO
}

// Add adds delta (which may be negative), clamping the result at zero.
func (c *Counter) Add(delta int) {
	// TODO
}

// Reset sets the counter back to zero.
func (c *Counter) Reset() {
	// TODO
}

// Value returns the current count.
func (c *Counter) Value() int {
	// TODO
	return 0
}

// CountWords returns one counter per distinct word, counting occurrences.
func CountWords(words []string) map[string]*Counter {
	// TODO
	return nil
}
