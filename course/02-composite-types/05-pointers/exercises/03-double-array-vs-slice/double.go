package double

// Quad is an array of four ints. Arrays are values: assigning or passing one
// copies all its elements.
type Quad [4]int

// DoubleSlice doubles every element of s in place.
func DoubleSlice(s []int) {
	// TODO
}

// DoubleArray doubles every element of the array a points to.
// It does nothing if a is nil.
func DoubleArray(a *Quad) {
	// TODO
}

// Doubled returns a copy of a with every element doubled.
// The caller's array is not changed.
func Doubled(a Quad) Quad {
	// TODO
	return Quad{}
}
