package extremes

import "cmp"

// MaxOf returns the largest element of xs and true,
// or the zero value of T and false when xs is empty.
func MaxOf[T cmp.Ordered](xs []T) (T, bool) {
	// TODO: declare `var zero T` for the empty case
	var zero T
	return zero, false
}

// MinOf returns the smallest element of xs and true,
// or the zero value of T and false when xs is empty.
func MinOf[T cmp.Ordered](xs []T) (T, bool) {
	// TODO
	var zero T
	return zero, false
}

// Clamp limits v to the range [lo, hi]. If lo > hi they are swapped first.
func Clamp[T cmp.Ordered](v, lo, hi T) T {
	// TODO
	var zero T
	return zero
}
