// Package lookup loads values through a cache.
package lookup

import "learngo/course/04-generics-packages-testing/04-packages/exercises/03-no-stutter/cache"

// Cached returns the cached value for key. On a miss it calls load,
// caches the value when load succeeds, and returns load's result.
func Cached(s *cache.Store, key string, load func(key string) (string, error)) (string, error) {
	if v, ok := s.Get(key); ok {
		return v, nil
	}
	v, err := load(key)
	if err != nil {
		return "", err
	}
	s.Set(key, v)
	return v, nil
}
