package counts

// MergeCounts returns a new map holding, for every key of a or b, the sum
// of its counts. Neither a nor b is modified. The result is never nil.
func MergeCounts(a, b map[string]int) map[string]int {
	// TODO
	return nil
}

// Prune deletes, in place, every entry of m whose count is less than minCount.
func Prune(m map[string]int, minCount int) {
	// TODO
}

// PrintCounts prints "key: value" lines sorted by key, or "(none)" if m is
// empty.
func PrintCounts(m map[string]int) {
	// TODO
}
