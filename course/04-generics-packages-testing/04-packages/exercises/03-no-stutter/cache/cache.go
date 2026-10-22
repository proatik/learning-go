// Package cache is a tiny in-memory string cache.
package cache

// TODO: every exported name below repeats the package name ("cache.CacheStore").
// Rename them as listed in the README, and make the zero value of the store
// usable (Set must create the map when it is nil).

// CacheStats counts cache hits and misses.
type CacheStats struct{ Hits, Misses int }

// CacheStore is an in-memory key/value cache.
type CacheStore struct {
	items map[string]string
	stats CacheStats
}

// NewCacheStore returns an empty cache.
func NewCacheStore() *CacheStore {
	return &CacheStore{items: map[string]string{}}
}

// CacheGet returns the value for key and whether it was present.
func (s *CacheStore) CacheGet(key string) (string, bool) {
	v, ok := s.items[key]
	if ok {
		s.stats.Hits++
	} else {
		s.stats.Misses++
	}
	return v, ok
}

// CacheSet stores value under key.
func (s *CacheStore) CacheSet(key, value string) {
	s.items[key] = value
}

// CacheDelete removes key; deleting a missing key does nothing.
func (s *CacheStore) CacheDelete(key string) {
	delete(s.items, key)
}

// CacheLen returns the number of cached keys.
func (s *CacheStore) CacheLen() int {
	return len(s.items)
}

// CacheStats returns the hit and miss counters.
func (s *CacheStore) CacheStats() CacheStats {
	return s.stats
}
