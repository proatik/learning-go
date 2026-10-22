package dedupe

// Dedupe returns xs without duplicates, keeping the first occurrence of each value.
func Dedupe[T comparable](xs []T) []T {
	return DedupeBy(xs, func(x T) T { return x })
}

// DedupeBy returns xs keeping only the first element for each key(x).
func DedupeBy[T any, K comparable](xs []T, key func(T) K) []T {
	seen := make(map[K]struct{}, len(xs))
	var out []T
	for _, x := range xs {
		k := key(x)
		if _, dup := seen[k]; dup {
			continue
		}
		seen[k] = struct{}{}
		out = append(out, x)
	}
	return out
}

// Intersect returns the distinct values of a that also appear in b, in a's order.
func Intersect[T comparable](a, b []T) []T {
	inB := make(map[T]struct{}, len(b))
	for _, x := range b {
		inB[x] = struct{}{}
	}
	var out []T
	for _, x := range Dedupe(a) {
		if _, ok := inB[x]; ok {
			out = append(out, x)
		}
	}
	return out
}
