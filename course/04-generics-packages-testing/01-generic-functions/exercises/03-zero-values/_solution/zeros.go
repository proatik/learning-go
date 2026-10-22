package zeros

// Must returns v, or panics with err if err is not nil.
func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

// FirstOr returns xs[0], or fallback when xs is empty.
func FirstOr[T any](xs []T, fallback T) T {
	if len(xs) == 0 {
		return fallback
	}
	return xs[0]
}

// IsZero reports whether v is the zero value of its type.
func IsZero[T comparable](v T) bool {
	var zero T
	return v == zero
}

// FirstNonZero returns the first non-zero value in vals,
// or the zero value if there is none.
func FirstNonZero[T comparable](vals ...T) T {
	for _, v := range vals {
		if !IsZero(v) {
			return v
		}
	}
	var zero T
	return zero
}
