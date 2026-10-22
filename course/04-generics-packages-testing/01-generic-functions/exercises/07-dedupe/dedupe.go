package dedupe

// Dedupe returns xs without duplicates, keeping the first occurrence of each value.
func Dedupe[T comparable](xs []T) []T {
	// TODO: a "seen" map[T]struct{}
	return nil
}

// DedupeBy returns xs keeping only the first element for each key(x).
func DedupeBy[T any, K comparable](xs []T, key func(T) K) []T {
	// TODO
	return nil
}

// Intersect returns the distinct values of a that also appear in b, in a's order.
func Intersect[T comparable](a, b []T) []T {
	// TODO
	return nil
}
