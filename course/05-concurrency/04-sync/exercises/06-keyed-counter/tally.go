package tally

import "sync"

// Counter counts occurrences of string keys. It is safe for concurrent use.
// The zero value is an empty counter, ready to use.
type Counter struct {
	mu     sync.Mutex     // guards counts
	counts map[string]int // created on first write
}

// Add adds n to key's count and returns the new count.
func (c *Counter) Add(key string, n int) int {
	// TODO
	return 0
}

// Inc adds 1 to key's count and returns the new count.
func (c *Counter) Inc(key string) int {
	// TODO
	return 0
}

// Get returns key's count, or 0 if key was never counted.
func (c *Counter) Get(key string) int {
	// TODO
	return 0
}

// Snapshot returns a copy of every count. Changing the copy doesn't affect
// the counter, and later counting doesn't change the copy.
func (c *Counter) Snapshot() map[string]int {
	// TODO
	return nil
}

// Reset removes every count and returns what they were.
func (c *Counter) Reset() map[string]int {
	// TODO
	return nil
}
