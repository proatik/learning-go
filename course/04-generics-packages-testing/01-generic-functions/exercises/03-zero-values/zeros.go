package zeros

// Must returns v, or panics with err if err is not nil.
func Must[T any](v T, err error) T {
	// TODO
	var zero T
	return zero
}

// FirstOr returns xs[0], or fallback when xs is empty.
func FirstOr[T any](xs []T, fallback T) T {
	// TODO
	var zero T
	return zero
}

// IsZero reports whether v is the zero value of its type.
func IsZero[T comparable](v T) bool {
	// TODO: compare with `var zero T`
	return false
}

// FirstNonZero returns the first non-zero value in vals,
// or the zero value if there is none.
func FirstNonZero[T comparable](vals ...T) T {
	// TODO
	var zero T
	return zero
}
