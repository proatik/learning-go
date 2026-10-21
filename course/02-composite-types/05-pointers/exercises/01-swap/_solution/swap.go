package swap

// Swap exchanges the values that a and b point to.
func Swap(a, b *int) {
	// The right side is fully evaluated before anything is assigned.
	*a, *b = *b, *a
}

// SwapIfGreater swaps *a and *b when *a > *b and reports whether it did.
func SwapIfGreater(a, b *int) bool {
	if *a <= *b {
		return false
	}
	Swap(a, b) // a and b are already pointers: pass them on as they are
	return true
}
