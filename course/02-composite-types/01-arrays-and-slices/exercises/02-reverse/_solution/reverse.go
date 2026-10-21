package reverse

// Reverse reverses the elements of s in place.
// Two indices walk toward each other, swapping as they go; for an empty or
// one-element slice the loop body never runs.
func Reverse(s []int) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

// ReversedArray returns a with its elements in reverse order.
// The caller's array is not modified: a is already a private copy.
func ReversedArray(a [5]int) [5]int {
	for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}
	return a
}
