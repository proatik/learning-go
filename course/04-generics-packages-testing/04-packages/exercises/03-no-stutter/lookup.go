// Package lookup loads values through a cache.
package lookup

import "learngo/course/04-generics-packages-testing/04-packages/exercises/03-no-stutter/cache"

// TODO: rename to Cached (callers write lookup.Cached) and switch to the new cache names.

// LookupCached returns the cached value for key. On a miss it calls load,
// caches the value when load succeeds, and returns load's result.
func LookupCached(s *cache.CacheStore, key string, load func(key string) (string, error)) (string, error) {
	if v, ok := s.CacheGet(key); ok {
		return v, nil
	}
	v, err := load(key)
	if err != nil {
		return "", err
	}
	s.CacheSet(key, v)
	return v, nil
}
