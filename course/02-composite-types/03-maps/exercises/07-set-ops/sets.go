package sets

// Set is a set of strings. Only the keys matter; the empty-struct values
// take no memory.
type Set map[string]struct{}

// NewSet returns a non-nil set containing items (duplicates collapse).
func NewSet(items ...string) Set {
	// TODO
	return nil
}

// Union returns the elements in a or b, sorted.
func Union(a, b Set) []string {
	// TODO
	return nil
}

// Intersection returns the elements in both a and b, sorted.
func Intersection(a, b Set) []string {
	// TODO
	return nil
}

// Difference returns the elements in a that are not in b, sorted.
func Difference(a, b Set) []string {
	// TODO
	return nil
}

// IsSubset reports whether every element of a is also in b.
func IsSubset(a, b Set) bool {
	// TODO
	return false
}
