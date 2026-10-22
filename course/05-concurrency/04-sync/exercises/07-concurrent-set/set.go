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
	// TODO
	return false
}

// Remove deletes v and reports whether it was there.
func (s *Set[T]) Remove(v T) bool {
	// TODO
	return false
}

// Has reports whether v is in the set.
func (s *Set[T]) Has(v T) bool {
	// TODO
	return false
}

// Len returns the number of values in the set.
func (s *Set[T]) Len() int {
	// TODO
	return 0
}

// Values returns the values in the set as a new slice, in no particular
// order.
func (s *Set[T]) Values() []T {
	// TODO
	return nil
}
