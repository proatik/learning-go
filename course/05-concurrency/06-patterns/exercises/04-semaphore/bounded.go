package bounded

// Map calls f on every element of in, running at most limit calls at once
// (limit < 1 means 1), and returns the results in input order.
func Map[T, R any](limit int, in []T, f func(T) R) []R {
	// TODO
	return nil
}
