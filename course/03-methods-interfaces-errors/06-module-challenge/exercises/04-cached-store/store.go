package kv

import "errors"

// Sentinel errors returned (wrapped) by every Store in this package.
var (
	ErrNotFound = errors.New("kv: key not found")
	ErrEmptyKey = errors.New("kv: empty key")
)

// Store is a minimal key-value store.
type Store interface {
	Get(key string) (string, error)
	Set(key, value string) error
	Delete(key string) error
}

// MemoryStore is a Store backed by a map. Create it with NewMemoryStore.
type MemoryStore struct {
	data map[string]string
}

// NewMemoryStore returns an empty MemoryStore.
func NewMemoryStore() *MemoryStore {
	// TODO
	return nil
}

// Get returns the value for key.
func (m *MemoryStore) Get(key string) (string, error) {
	// TODO
	return "", nil
}

// Set stores value under key.
func (m *MemoryStore) Set(key, value string) error {
	// TODO
	return nil
}

// Delete removes key.
func (m *MemoryStore) Delete(key string) error {
	// TODO
	return nil
}

// Len returns the number of stored keys.
func (m *MemoryStore) Len() int {
	// TODO
	return 0
}

// CachedStore is a read-through cache in front of another Store.
type CachedStore struct {
	backend      Store
	cache        map[string]string
	hits, misses int
}

// NewCachedStore wraps backend with an in-memory cache.
func NewCachedStore(backend Store) *CachedStore {
	// TODO
	return nil
}

// Get serves key from the cache, or loads it from the backend and caches it.
func (c *CachedStore) Get(key string) (string, error) {
	// TODO
	return "", nil
}

// Set writes to the backend first and caches the value only on success.
func (c *CachedStore) Set(key, value string) error {
	// TODO
	return nil
}

// Delete drops key from the cache and deletes it from the backend.
func (c *CachedStore) Delete(key string) error {
	// TODO
	return nil
}

// Stats returns the number of cache hits and misses so far.
func (c *CachedStore) Stats() (hits, misses int) {
	// TODO
	return 0, 0
}

// String returns "cache: H hits, M misses, N cached".
func (c *CachedStore) String() string {
	// TODO
	return ""
}
