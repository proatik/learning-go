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
	// TODO
	return &Cache{}
}

// Get returns the value stored under key; ok is false if key is missing.
func (c *Cache) Get(key string) (value string, ok bool) {
	// TODO
	return "", false
}

// Set stores value under key, replacing any previous value.
func (c *Cache) Set(key, value string) {
	// TODO
}

// Delete removes key. Deleting a missing key does nothing.
func (c *Cache) Delete(key string) {
	// TODO
}

// Len returns the number of keys.
func (c *Cache) Len() int {
	// TODO
	return 0
}
