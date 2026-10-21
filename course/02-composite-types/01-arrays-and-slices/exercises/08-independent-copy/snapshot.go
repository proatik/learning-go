package snapshot

// Clone returns a copy of s that shares no memory with it.
// A nil s gives nil; a non-nil empty s gives a non-nil empty slice.
func Clone(s []int) []int {
	// TODO
	return nil
}

// WithAppended returns a new slice holding the elements of s followed by v.
// It never writes to s's backing array, even if s has spare capacity.
func WithAppended(s []int, v int) []int {
	// TODO
	return nil
}

// Clone2D returns a deep copy of m: the outer slice and every row are new.
// nil rows stay nil; a nil m gives nil.
func Clone2D(m [][]int) [][]int {
	// TODO
	return nil
}
