package cache

import "sync"

// Cache is a string-to-string cache that is safe for concurrent use. Many
// readers may hold the lock at once; a writer holds it alone. Create one
// with New.
type Cache struct {
	mu    sync.RWMutex // guards items
	items map[string]string
}

// New returns an empty cache.
func New() *Cache {
	return &Cache{items: make(map[string]string)}
}

// Get returns the value stored under key; ok is false if key is missing.
func (c *Cache) Get(key string) (value string, ok bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, ok = c.items[key]
	return value, ok
}

// Set stores value under key, replacing any previous value.
func (c *Cache) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = value
}

// Delete removes key. Deleting a missing key does nothing.
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

// Len returns the number of keys.
func (c *Cache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items)
}
