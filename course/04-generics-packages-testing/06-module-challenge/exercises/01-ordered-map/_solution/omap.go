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
	return &OrderedMap[K, V]{index: make(map[K]*entry[K, V])}
}

// Collect builds an OrderedMap from seq. A repeated key keeps its first
// position and its last value.
func Collect[K comparable, V any](seq iter.Seq2[K, V]) *OrderedMap[K, V] {
	m := New[K, V]()
	for k, v := range seq {
		m.Set(k, v)
	}
	return m
}

// Set stores value under key. A new key goes to the end; an existing key
// keeps its position.
func (m *OrderedMap[K, V]) Set(key K, value V) {
	if e, ok := m.index[key]; ok {
		e.value = value
		return
	}
	e := &entry[K, V]{key: key, value: value, prev: m.tail}
	if m.tail != nil {
		m.tail.next = e
	} else {
		m.head = e
	}
	m.tail = e
	m.index[key] = e
}

// Get returns the value for key and whether it was present.
func (m *OrderedMap[K, V]) Get(key K) (V, bool) {
	if e, ok := m.index[key]; ok {
		return e.value, true
	}
	var zero V
	return zero, false
}

// Delete removes key and reports whether it was present.
func (m *OrderedMap[K, V]) Delete(key K) bool {
	e, ok := m.index[key]
	if !ok {
		return false
	}
	if e.prev != nil {
		e.prev.next = e.next
	} else {
		m.head = e.next
	}
	if e.next != nil {
		e.next.prev = e.prev
	} else {
		m.tail = e.prev
	}
	delete(m.index, key)
	return true
}

// Len returns the number of entries.
func (m *OrderedMap[K, V]) Len() int {
	return len(m.index)
}

// All yields the pairs from oldest to newest. Deleting the current key
// while ranging is allowed.
func (m *OrderedMap[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for e := m.head; e != nil; {
			next := e.next
			if !yield(e.key, e.value) {
				return
			}
			e = next
		}
	}
}

// Keys yields the keys from oldest to newest.
func (m *OrderedMap[K, V]) Keys() iter.Seq[K] {
	return func(yield func(K) bool) {
		for k := range m.All() {
			if !yield(k) {
				return
			}
		}
	}
}

// Backward yields the pairs from newest to oldest.
func (m *OrderedMap[K, V]) Backward() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for e := m.tail; e != nil; {
			prev := e.prev
			if !yield(e.key, e.value) {
				return
			}
			e = prev
		}
	}
}
