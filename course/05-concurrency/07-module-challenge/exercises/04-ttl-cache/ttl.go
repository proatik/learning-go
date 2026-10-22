package ttl

import (
	"sync"
	"time"
)

// entry is a stored value with its expiry time.
type entry[V any] struct {
	val     V
	expires time.Time
}

// Cache is a concurrent map whose entries expire. Create one with New and
// stop its janitor with Close.
type Cache[K comparable, V any] struct {
	mu        sync.Mutex
	items     map[K]entry[V]
	stop      chan struct{} // closed by Close to stop the janitor
	done      chan struct{} // closed by the janitor when it exits
	closeOnce sync.Once
}

// New returns an empty cache and starts a janitor that deletes expired
// entries every sweepEvery.
func New[K comparable, V any](sweepEvery time.Duration) *Cache[K, V] {
	// TODO
	return &Cache[K, V]{}
}

// Set stores v under k for ttl, replacing any previous value and expiry.
func (c *Cache[K, V]) Set(k K, v V, ttl time.Duration) {
	// TODO
}

// Get returns the value for k; ok is false if k is missing or expired.
func (c *Cache[K, V]) Get(k K) (V, bool) {
	// TODO
	var zero V
	return zero, false
}

// Delete removes k.
func (c *Cache[K, V]) Delete(k K) {
	// TODO
}

// Len returns the number of stored entries, including expired ones the
// janitor hasn't deleted yet.
func (c *Cache[K, V]) Len() int {
	// TODO
	return 0
}

// Close stops the janitor and waits for it to exit. It is safe to call more
// than once.
func (c *Cache[K, V]) Close() {
	// TODO
}
