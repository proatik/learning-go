package num

// Sign returns -1 if n is negative, 0 if n is zero and 1 if n is positive.
func Sign(n int) int {
	if n < 0 {
		return -1
	} else if n > 0 {
		return 1
	} else {
		return 0
	}
}

// Abs returns the absolute value of n.
func Abs(n int) int {
	if n < 0 {
		return -n
	}
	return n // early return above means no else is needed
}

// Clamp returns x limited to the inclusive range [lo, hi]. It assumes lo <= hi.
func Clamp(x, lo, hi int) int {
	// min caps x at hi; max then raises the result to at least lo.
	return max(lo, min(x, hi))
}
