// Package cache is a tiny in-memory string cache.
package cache

// Stats counts cache hits and misses.
type Stats struct{ Hits, Misses int }

// Store is an in-memory key/value cache. The zero value is an empty cache
// ready to use.
type Store struct {
	items map[string]string
	stats Stats
}

// New returns an empty cache.
func New() *Store {
	return &Store{items: map[string]string{}}
}

// Get returns the value for key and whether it was present.
func (s *Store) Get(key string) (string, bool) {
	v, ok := s.items[key] // reading a nil map is fine
	if ok {
		s.stats.Hits++
	} else {
		s.stats.Misses++
	}
	return v, ok
}

// Set stores value under key.
func (s *Store) Set(key, value string) {
	if s.items == nil {
		s.items = map[string]string{}
	}
	s.items[key] = value
}

// Delete removes key; deleting a missing key does nothing.
func (s *Store) Delete(key string) {
	delete(s.items, key)
}

// Len returns the number of cached keys.
func (s *Store) Len() int {
	return len(s.items)
}

// Stats returns the hit and miss counters.
func (s *Store) Stats() Stats {
	return s.stats
}
