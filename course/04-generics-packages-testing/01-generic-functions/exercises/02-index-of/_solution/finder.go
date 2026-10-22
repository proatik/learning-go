package finder

// IndexOf returns the index of the first element equal to target, or -1.
func IndexOf[T comparable](xs []T, target T) int {
	for i, x := range xs {
		if x == target {
			return i
		}
	}
	return -1
}

// LastIndexOf returns the index of the last element equal to target, or -1.
func LastIndexOf[T comparable](xs []T, target T) int {
	for i := len(xs) - 1; i >= 0; i-- {
		if xs[i] == target {
			return i
		}
	}
	return -1
}

// Contains reports whether target is in xs.
func Contains[T comparable](xs []T, target T) bool {
	return IndexOf(xs, target) >= 0
}

// Count returns how many elements of xs equal target.
func Count[T comparable](xs []T, target T) int {
	n := 0
	for _, x := range xs {
		if x == target {
			n++
		}
	}
	return n
}
