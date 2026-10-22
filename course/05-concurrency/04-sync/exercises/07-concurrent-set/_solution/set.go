package syncset

import "sync"

// Set is a set of comparable values that is safe for concurrent use. The
// zero value is an empty set, ready to use.
type Set[T comparable] struct {
	mu    sync.RWMutex   // guards items
	items map[T]struct{} // created on first Add
}

// Add inserts v and reports whether it was newly added (false if it was
// already there).
func (s *Set[T]) Add(v T) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.items[v]; ok {
		return false
	}
	if s.items == nil {
		s.items = make(map[T]struct{})
	}
	s.items[v] = struct{}{}
	return true
}

// Remove deletes v and reports whether it was there.
func (s *Set[T]) Remove(v T) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.items[v]; !ok {
		return false
	}
	delete(s.items, v)
	return true
}

// Has reports whether v is in the set.
func (s *Set[T]) Has(v T) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.items[v]
	return ok
}

// Len returns the number of values in the set.
func (s *Set[T]) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.items)
}

// Values returns the values in the set as a new slice, in no particular
// order.
func (s *Set[T]) Values() []T {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]T, 0, len(s.items))
	for v := range s.items {
		out = append(out, v)
	}
	return out
}
