package swap

// Swap returns its two arguments in reverse order.
func Swap(a, b int) (int, int) {
	// All right-hand values are evaluated before any assignment happens,
	// so no temporary variable is needed.
	a, b = b, a
	return a, b
}

// Rotate shifts its arguments one position to the left: (a, b, c) becomes (b, c, a).
func Rotate(a, b, c int) (int, int, int) {
	a, b, c = b, c, a
	return a, b, c
}

// SwapStrings returns its two arguments in reverse order.
func SwapStrings(a, b string) (string, string) {
	a, b = b, a
	return a, b
}
