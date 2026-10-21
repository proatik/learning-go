package kv

import (
	"errors"
	"fmt"
)

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

// Compile-time checks.
var (
	_ Store = (*MemoryStore)(nil)
	_ Store = (*CachedStore)(nil)
)

// MemoryStore is a Store backed by a map. Create it with NewMemoryStore.
type MemoryStore struct {
	data map[string]string
}

// NewMemoryStore returns an empty MemoryStore.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{data: make(map[string]string)}
}

// Get returns the value for key.
func (m *MemoryStore) Get(key string) (string, error) {
	if key == "" {
		return "", ErrEmptyKey
	}
	v, ok := m.data[key]
	if !ok {
		return "", fmt.Errorf("get %q: %w", key, ErrNotFound)
	}
	return v, nil
}

// Set stores value under key.
func (m *MemoryStore) Set(key, value string) error {
	if key == "" {
		return ErrEmptyKey
	}
	m.data[key] = value
	return nil
}

// Delete removes key.
func (m *MemoryStore) Delete(key string) error {
	if key == "" {
		return ErrEmptyKey
	}
	if _, ok := m.data[key]; !ok {
		return fmt.Errorf("delete %q: %w", key, ErrNotFound)
	}
	delete(m.data, key)
	return nil
}

// Len returns the number of stored keys.
func (m *MemoryStore) Len() int { return len(m.data) }

// CachedStore is a read-through cache in front of another Store.
type CachedStore struct {
	backend      Store
	cache        map[string]string
	hits, misses int
}

// NewCachedStore wraps backend with an in-memory cache.
func NewCachedStore(backend Store) *CachedStore {
	return &CachedStore{backend: backend, cache: make(map[string]string)}
}

// Get serves key from the cache, or loads it from the backend and caches it.
func (c *CachedStore) Get(key string) (string, error) {
	if key == "" {
		return "", ErrEmptyKey
	}
	if v, ok := c.cache[key]; ok {
		c.hits++
		return v, nil
	}
	c.misses++
	v, err := c.backend.Get(key)
	if err != nil {
		return "", fmt.Errorf("cache: %w", err)
	}
	c.cache[key] = v
	return v, nil
}

// Set writes to the backend first and caches the value only on success.
func (c *CachedStore) Set(key, value string) error {
	if key == "" {
		return ErrEmptyKey
	}
	if err := c.backend.Set(key, value); err != nil {
		return fmt.Errorf("cache: %w", err)
	}
	c.cache[key] = value
	return nil
}

// Delete drops key from the cache and deletes it from the backend.
func (c *CachedStore) Delete(key string) error {
	if key == "" {
		return ErrEmptyKey
	}
	delete(c.cache, key)
	if err := c.backend.Delete(key); err != nil {
		return fmt.Errorf("cache: %w", err)
	}
	return nil
}

// Stats returns the number of cache hits and misses so far.
func (c *CachedStore) Stats() (hits, misses int) { return c.hits, c.misses }

// String returns "cache: H hits, M misses, N cached".
func (c *CachedStore) String() string {
	return fmt.Sprintf("cache: %d hits, %d misses, %d cached", c.hits, c.misses, len(c.cache))
}
