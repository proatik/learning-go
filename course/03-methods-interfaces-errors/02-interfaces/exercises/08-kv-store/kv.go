package kv

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

// NewMemStore returns an empty store.
func NewMemStore() *MemStore {
	// TODO
	return &MemStore{}
}

// Get returns the value for key and whether it exists.
func (m *MemStore) Get(key string) (string, bool) {
	// TODO
	return "", false
}

// Set stores value under key, replacing any previous value.
func (m *MemStore) Set(key, value string) {
	// TODO
}

// Delete removes key and reports whether it existed.
func (m *MemStore) Delete(key string) bool {
	// TODO
	return false
}

// Len returns the number of keys.
func (m *MemStore) Len() int {
	// TODO
	return 0
}

// Keys returns all keys in ascending order.
func (m *MemStore) Keys() []string {
	// TODO
	return nil
}

// GetOr returns the value for key in g, or def if the key is missing.
func GetOr(g Getter, key, def string) string {
	// TODO
	return ""
}

// CopyKeys copies every key of keys that exists in src into dst and
// returns how many were copied.
func CopyKeys(dst Setter, src Getter, keys []string) int {
	// TODO
	return 0
}

// ReadOnly returns a view of g that only has the Get method.
func ReadOnly(g Getter) Getter {
	// TODO
	return nil
}
