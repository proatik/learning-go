package extremes

import "cmp"

// MaxOf returns the largest element of xs and true,
// or the zero value of T and false when xs is empty.
func MaxOf[T cmp.Ordered](xs []T) (T, bool) {
	if len(xs) == 0 {
		var zero T
		return zero, false
	}
	best := xs[0]
	for _, x := range xs[1:] {
		if x > best {
			best = x
		}
	}
	return best, true
}

// MinOf returns the smallest element of xs and true,
// or the zero value of T and false when xs is empty.
func MinOf[T cmp.Ordered](xs []T) (T, bool) {
	if len(xs) == 0 {
		var zero T
		return zero, false
	}
	best := xs[0]
	for _, x := range xs[1:] {
		if x < best {
			best = x
		}
	}
	return best, true
}

// Clamp limits v to the range [lo, hi]. If lo > hi they are swapped first.
func Clamp[T cmp.Ordered](v, lo, hi T) T {
	if lo > hi {
		lo, hi = hi, lo
	}
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
