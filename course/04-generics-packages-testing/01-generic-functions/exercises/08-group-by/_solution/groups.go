package groups

import "cmp"

// GroupBy groups the elements of xs by key(x), keeping their original order.
func GroupBy[T any, K comparable](xs []T, key func(T) K) map[K][]T {
	groups := make(map[K][]T)
	for _, x := range xs {
		k := key(x)
		groups[k] = append(groups[k], x)
	}
	return groups
}

// Partition splits xs into the elements for which pred is true and the rest.
func Partition[T any](xs []T, pred func(T) bool) (yes, no []T) {
	for _, x := range xs {
		if pred(x) {
			yes = append(yes, x)
		} else {
			no = append(no, x)
		}
	}
	return yes, no
}

// MaxBy returns the first element with the largest key(x).
// ok is false when xs is empty.
func MaxBy[T any, K cmp.Ordered](xs []T, key func(T) K) (T, bool) {
	if len(xs) == 0 {
		var zero T
		return zero, false
	}
	best, bestKey := xs[0], key(xs[0])
	for _, x := range xs[1:] {
		if k := key(x); k > bestKey {
			best, bestKey = x, k
		}
	}
	return best, true
}
