package bounds

// MinMax returns the smallest (lo) and largest (hi) of its arguments.
func MinMax(first int, rest ...int) (lo, hi int) {
	lo, hi = first, first
	for _, n := range rest {
		if n < lo {
			lo = n
		}
		if n > hi {
			hi = n
		}
	}
	return lo, hi
}
