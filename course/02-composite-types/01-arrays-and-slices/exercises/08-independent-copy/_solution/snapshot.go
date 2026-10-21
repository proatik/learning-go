package snapshot

// Clone returns a copy of s that shares no memory with it.
// A nil s gives nil; a non-nil empty s gives a non-nil empty slice.
func Clone(s []int) []int {
	if s == nil {
		return nil
	}
	c := make([]int, len(s)) // new backing array; non-nil even when len is 0
	copy(c, s)
	return c
}

// WithAppended returns a new slice holding the elements of s followed by v.
// It never writes to s's backing array, even if s has spare capacity.
//
// A plain append(s, v) would write into s's spare capacity, so two calls on
// the same s would overwrite each other. Allocating exactly len(s)+1 avoids it.
func WithAppended(s []int, v int) []int {
	out := make([]int, len(s), len(s)+1)
	copy(out, s)
	return append(out, v) // fits in cap: no second allocation
}

// Clone2D returns a deep copy of m: the outer slice and every row are new.
// nil rows stay nil; a nil m gives nil.
func Clone2D(m [][]int) [][]int {
	if m == nil {
		return nil
	}
	out := make([][]int, len(m))
	for i, row := range m {
		out[i] = Clone(row) // copying only the outer slice would share the rows
	}
	return out
}
