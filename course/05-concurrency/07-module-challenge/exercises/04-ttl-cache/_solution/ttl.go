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
	if sweepEvery <= 0 {
		sweepEvery = time.Second
	}
	c := &Cache[K, V]{
		items: make(map[K]entry[V]),
		stop:  make(chan struct{}),
		done:  make(chan struct{}),
	}
	go c.janitor(sweepEvery)
	return c
}

// janitor deletes expired entries on every tick until stop is closed.
func (c *Cache[K, V]) janitor(every time.Duration) {
	defer close(c.done)
	tick := time.NewTicker(every)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			c.sweep()
		case <-c.stop:
			return
		}
	}
}

// sweep deletes every expired entry.
func (c *Cache[K, V]) sweep() {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now()
	for k, e := range c.items {
		if !now.Before(e.expires) {
			delete(c.items, k)
		}
	}
}

// Set stores v under k for ttl, replacing any previous value and expiry.
func (c *Cache[K, V]) Set(k K, v V, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[k] = entry[V]{val: v, expires: time.Now().Add(ttl)}
}

// Get returns the value for k; ok is false if k is missing or expired.
func (c *Cache[K, V]) Get(k K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	e, ok := c.items[k]
	if !ok || !time.Now().Before(e.expires) {
		var zero V
		return zero, false
	}
	return e.val, true
}

// Delete removes k.
func (c *Cache[K, V]) Delete(k K) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, k)
}

// Len returns the number of stored entries, including expired ones the
// janitor hasn't deleted yet.
func (c *Cache[K, V]) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.items)
}

// Close stops the janitor and waits for it to exit. It is safe to call more
// than once.
func (c *Cache[K, V]) Close() {
	c.closeOnce.Do(func() { close(c.stop) })
	<-c.done
}
