package sortedmap

import (
	"cmp"
	"slices"
)

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
	if m.vals == nil {
		m.vals = make(map[K]V)
	}
	if _, ok := m.vals[k]; !ok {
		i, _ := slices.BinarySearch(m.keys, k)
		m.keys = slices.Insert(m.keys, i, k)
	}
	m.vals[k] = v
}

// Get returns the value for k; ok is false if k is missing.
func (m *Map[K, V]) Get(k K) (V, bool) {
	v, ok := m.vals[k]
	return v, ok
}

// Delete removes k and reports whether it was present.
func (m *Map[K, V]) Delete(k K) bool {
	i, found := slices.BinarySearch(m.keys, k)
	if !found {
		return false
	}
	m.keys = slices.Delete(m.keys, i, i+1)
	delete(m.vals, k)
	return true
}

// Len returns the number of keys.
func (m *Map[K, V]) Len() int {
	return len(m.keys)
}

// Keys returns a copy of the keys in ascending order.
func (m *Map[K, V]) Keys() []K {
	return slices.Clone(m.keys)
}

// Entries returns every entry in ascending key order.
func (m *Map[K, V]) Entries() []Entry[K, V] {
	return m.entries(m.keys)
}

// Range returns the entries with lo <= key < hi, in ascending key order.
func (m *Map[K, V]) Range(lo, hi K) []Entry[K, V] {
	if cmp.Compare(lo, hi) >= 0 {
		return nil
	}
	from, _ := slices.BinarySearch(m.keys, lo)
	to, _ := slices.BinarySearch(m.keys, hi)
	return m.entries(m.keys[from:to])
}

// Min returns the entry with the smallest key; ok is false if m is empty.
func (m *Map[K, V]) Min() (Entry[K, V], bool) {
	if len(m.keys) == 0 {
		return Entry[K, V]{}, false
	}
	k := m.keys[0]
	return Entry[K, V]{Key: k, Val: m.vals[k]}, true
}

// entries pairs each of keys with its value.
func (m *Map[K, V]) entries(keys []K) []Entry[K, V] {
	out := make([]Entry[K, V], len(keys))
	for i, k := range keys {
		out[i] = Entry[K, V]{Key: k, Val: m.vals[k]}
	}
	return out
}
