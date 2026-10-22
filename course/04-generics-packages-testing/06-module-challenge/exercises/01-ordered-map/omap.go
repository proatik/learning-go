// Package omap provides a map that remembers insertion order.
package omap

import "iter"

// entry is one key/value pair in the insertion-order list.
type entry[K comparable, V any] struct {
	key        K
	value      V
	prev, next *entry[K, V]
}

// OrderedMap is a map that iterates in insertion order. Create one with New.
type OrderedMap[K comparable, V any] struct {
	index map[K]*entry[K, V]
	head  *entry[K, V] // oldest
	tail  *entry[K, V] // newest
}

// New returns an empty OrderedMap.
func New[K comparable, V any]() *OrderedMap[K, V] {
	// TODO
	return nil
}

// Collect builds an OrderedMap from seq. A repeated key keeps its first
// position and its last value.
func Collect[K comparable, V any](seq iter.Seq2[K, V]) *OrderedMap[K, V] {
	// TODO
	return nil
}

// Set stores value under key. A new key goes to the end; an existing key
// keeps its position.
func (m *OrderedMap[K, V]) Set(key K, value V) {
	// TODO
}

// Get returns the value for key and whether it was present.
func (m *OrderedMap[K, V]) Get(key K) (V, bool) {
	// TODO
	var zero V
	return zero, false
}

// Delete removes key and reports whether it was present.
func (m *OrderedMap[K, V]) Delete(key K) bool {
	// TODO
	return false
}

// Len returns the number of entries.
func (m *OrderedMap[K, V]) Len() int {
	// TODO
	return 0
}

// All yields the pairs from oldest to newest. Deleting the current key
// while ranging is allowed.
func (m *OrderedMap[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		// TODO
	}
}

// Keys yields the keys from oldest to newest.
func (m *OrderedMap[K, V]) Keys() iter.Seq[K] {
	return func(yield func(K) bool) {
		// TODO
	}
}

// Backward yields the pairs from newest to oldest.
func (m *OrderedMap[K, V]) Backward() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		// TODO
	}
}
