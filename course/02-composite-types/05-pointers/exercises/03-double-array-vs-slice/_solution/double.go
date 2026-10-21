package double

// Quad is an array of four ints. Arrays are values: assigning or passing one
// copies all its elements.
type Quad [4]int

// DoubleSlice doubles every element of s in place.
func DoubleSlice(s []int) {
	// s is a copy of the slice header, but it points to the caller's
	// backing array, so writes through s[i] are visible to the caller.
	for i := range s {
		s[i] *= 2
	}
}

// DoubleArray doubles every element of the array a points to.
// It does nothing if a is nil.
func DoubleArray(a *Quad) {
	if a == nil {
		return
	}
	// Indexing and ranging over a pointer to an array dereference it
	// automatically: a[i] means (*a)[i].
	for i := range a {
		a[i] *= 2
	}
}

// Doubled returns a copy of a with every element doubled.
// The caller's array is not changed.
func Doubled(a Quad) Quad {
	// a is already a private copy, so we may modify and return it.
	for i := range a {
		a[i] *= 2
	}
	return a
}
