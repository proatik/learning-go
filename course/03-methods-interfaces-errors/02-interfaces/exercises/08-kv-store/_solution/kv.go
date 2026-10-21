package kv

import (
	"maps"
	"slices"
)

// Getter reads values by key.
type Getter interface {
	Get(key string) (string, bool)
}

// Setter stores values by key.
type Setter interface {
	Set(key, value string)
}

// Deleter removes keys.
type Deleter interface {
	Delete(key string) bool
}

// Store is a full key-value store.
type Store interface {
	Getter
	Setter
	Deleter
	Len() int
}

// MemStore is an in-memory Store; its zero value is an empty store.
type MemStore struct {
	data map[string]string
}

var _ Store = (*MemStore)(nil)

// NewMemStore returns an empty store.
func NewMemStore() *MemStore {
	return &MemStore{data: map[string]string{}}
}

// Get returns the value for key and whether it exists.
func (m *MemStore) Get(key string) (string, bool) {
	v, ok := m.data[key] // reading a nil map is fine
	return v, ok
}

// Set stores value under key, replacing any previous value.
func (m *MemStore) Set(key, value string) {
	if m.data == nil {
		m.data = map[string]string{}
	}
	m.data[key] = value
}

// Delete removes key and reports whether it existed.
func (m *MemStore) Delete(key string) bool {
	if _, ok := m.data[key]; !ok {
		return false
	}
	delete(m.data, key)
	return true
}

// Len returns the number of keys.
func (m *MemStore) Len() int { return len(m.data) }

// Keys returns all keys in ascending order.
func (m *MemStore) Keys() []string {
	return slices.Sorted(maps.Keys(m.data))
}

// GetOr returns the value for key in g, or def if the key is missing.
func GetOr(g Getter, key, def string) string {
	if v, ok := g.Get(key); ok {
		return v
	}
	return def
}

// CopyKeys copies every key of keys that exists in src into dst and
// returns how many were copied.
func CopyKeys(dst Setter, src Getter, keys []string) int {
	n := 0
	for _, k := range keys {
		if v, ok := src.Get(k); ok {
			dst.Set(k, v)
			n++
		}
	}
	return n
}

// readOnly wraps a Getter and exposes nothing else.
type readOnly struct {
	g Getter
}

// Get reads through to the wrapped Getter.
func (r readOnly) Get(key string) (string, bool) { return r.g.Get(key) }

// ReadOnly returns a view of g that only has the Get method.
func ReadOnly(g Getter) Getter {
	return readOnly{g: g}
}
