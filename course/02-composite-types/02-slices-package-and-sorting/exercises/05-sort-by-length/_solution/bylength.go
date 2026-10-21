package bylength

import (
	"cmp"
	"slices"
	"strings"
	"unicode/utf8"
)

// byLength orders by rune count, then alphabetically.
func byLength(a, b string) int {
	return cmp.Or(
		cmp.Compare(utf8.RuneCountInString(a), utf8.RuneCountInString(b)),
		strings.Compare(a, b), // tie-breaker: only decides when lengths are equal
	)
}

// SortByLength sorts words in place by rune count, shortest first.
// Words of equal length are sorted alphabetically.
func SortByLength(words []string) {
	slices.SortFunc(words, byLength) // a named function works as a comparator too
}

// SortByLengthDesc sorts words in place by rune count, longest first.
// Words of equal length are sorted alphabetically (ascending).
func SortByLengthDesc(words []string) {
	slices.SortFunc(words, func(a, b string) int {
		return cmp.Or(
			cmp.Compare(utf8.RuneCountInString(b), utf8.RuneCountInString(a)), // b, a: descending
			strings.Compare(a, b), // a, b: still ascending
		)
	})
}

// Shortest returns the first n words in SortByLength order, as a new slice.
// If n exceeds len(words) it returns all of them; if n <= 0 an empty slice.
// words is not modified.
func Shortest(words []string, n int) []string {
	sorted := slices.Clone(words)
	SortByLength(sorted)
	n = max(0, min(n, len(sorted))) // clamp n into [0, len]
	return sorted[:n]
}
