package edit

// RemoveAt removes s[i], keeping the order of the remaining elements,
// and returns the shortened slice. It may overwrite s's backing array.
//
// s[:i] has room for the tail (it shares s's array), so append just shifts
// s[i+1:] one slot to the left without allocating.
func RemoveAt(s []int, i int) []int {
	return append(s[:i], s[i+1:]...)
}

// RemoveUnordered removes s[i] in O(1) by moving the last element into its
// place, and returns the shortened slice.
func RemoveUnordered(s []int, i int) []int {
	last := len(s) - 1
	s[i] = s[last]
	return s[:last]
}

// InsertAt inserts v at index i (0 <= i <= len(s)), shifting later elements
// right, and returns the lengthened slice.
func InsertAt(s []int, i int, v int) []int {
	s = append(s, 0)     // make room for one more (may reallocate)
	copy(s[i+1:], s[i:]) // shift the tail right; copy handles the overlap
	s[i] = v
	return s
}
