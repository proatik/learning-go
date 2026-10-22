package groups

import "cmp"

// GroupBy groups the elements of xs by key(x), keeping their original order.
func GroupBy[T any, K comparable](xs []T, key func(T) K) map[K][]T {
	// TODO
	return nil
}

// Partition splits xs into the elements for which pred is true and the rest.
func Partition[T any](xs []T, pred func(T) bool) (yes, no []T) {
	// TODO
	return nil, nil
}

// MaxBy returns the first element with the largest key(x).
// ok is false when xs is empty.
func MaxBy[T any, K cmp.Ordered](xs []T, key func(T) K) (T, bool) {
	// TODO: remember both the best element and its key
	var zero T
	return zero, false
}
