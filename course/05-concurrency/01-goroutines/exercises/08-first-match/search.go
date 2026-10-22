package search

// FirstMatch calls pred on every item concurrently and returns the smallest
// index whose item matches, or -1 if none do.
func FirstMatch[T any](items []T, pred func(T) bool) int {
	// TODO
	return 0
}
