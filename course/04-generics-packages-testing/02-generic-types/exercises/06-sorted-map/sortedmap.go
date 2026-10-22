package sortedmap

import "cmp"

// Entry is one key/value pair of a Map.
type Entry[K, V any] struct {
	Key K
	Val V
}

// Map is a map whose keys are kept in ascending order. The zero value is an
// empty map.
type Map[K cmp.Ordered, V any] struct {
	keys []K     // always sorted, no duplicates
	vals map[K]V // key → value
}

// Put sets the value for k, adding k if it is new.
func (m *Map[K, V]) Put(k K, v V) {
	// TODO
}

// Get returns the value for k; ok is false if k is missing.
func (m *Map[K, V]) Get(k K) (V, bool) {
	// TODO
	var zero V
	return zero, false
}

// Delete removes k and reports whether it was present.
func (m *Map[K, V]) Delete(k K) bool {
	// TODO
	return false
}

// Len returns the number of keys.
func (m *Map[K, V]) Len() int {
	// TODO
	return 0
}

// Keys returns a copy of the keys in ascending order.
func (m *Map[K, V]) Keys() []K {
	// TODO
	return nil
}

// Entries returns every entry in ascending key order.
func (m *Map[K, V]) Entries() []Entry[K, V] {
	// TODO
	return nil
}

// Range returns the entries with lo <= key < hi, in ascending key order.
func (m *Map[K, V]) Range(lo, hi K) []Entry[K, V] {
	// TODO
	return nil
}

// Min returns the entry with the smallest key; ok is false if m is empty.
func (m *Map[K, V]) Min() (Entry[K, V], bool) {
	// TODO
	return Entry[K, V]{}, false
}
