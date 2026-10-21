package sets

import (
	"maps"
	"slices"
)

// Set is a set of strings. Only the keys matter; the empty-struct values
// take no memory.
type Set map[string]struct{}

// NewSet returns a non-nil set containing items (duplicates collapse).
func NewSet(items ...string) Set {
	s := make(Set, len(items))
	for _, it := range items {
		s[it] = struct{}{}
	}
	return s
}

// Union returns the elements in a or b, sorted.
func Union(a, b Set) []string {
	u := make(Set, len(a)+len(b)) // a fresh set: never write into a or b
	maps.Copy(u, a)               // copying from a nil map copies nothing
	maps.Copy(u, b)
	return slices.Sorted(maps.Keys(u))
}

// Intersection returns the elements in both a and b, sorted.
func Intersection(a, b Set) []string {
	var out []string
	for x := range a {
		if _, ok := b[x]; ok { // lookups in a nil b are fine: ok is false
			out = append(out, x)
		}
	}
	slices.Sort(out)
	return out
}

// Difference returns the elements in a that are not in b, sorted.
func Difference(a, b Set) []string {
	var out []string
	for x := range a {
		if _, ok := b[x]; !ok {
			out = append(out, x)
		}
	}
	slices.Sort(out)
	return out
}

// IsSubset reports whether every element of a is also in b.
func IsSubset(a, b Set) bool {
	for x := range a {
		if _, ok := b[x]; !ok {
			return false
		}
	}
	return true
}
