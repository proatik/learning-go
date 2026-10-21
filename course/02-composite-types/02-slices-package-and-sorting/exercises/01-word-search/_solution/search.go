package search

import "slices"

// HasWord reports whether w is one of words.
func HasWord(words []string, w string) bool {
	return slices.Contains(words, w)
}

// Position returns the index of the first occurrence of w in words, or -1.
func Position(words []string, w string) int {
	return slices.Index(words, w)
}

// FirstLonger returns the index of the first word longer than n bytes, or -1.
func FirstLonger(words []string, n int) int {
	// The closure captures n from the enclosing function.
	return slices.IndexFunc(words, func(w string) bool { return len(w) > n })
}

// AnyEmpty reports whether words contains an empty string.
func AnyEmpty(words []string) bool {
	return slices.ContainsFunc(words, func(w string) bool { return w == "" })
}

// SameWords reports whether a and b contain the same words in the same order.
func SameWords(a, b []string) bool {
	// == does not compile for slices; Equal compares length and elements.
	return slices.Equal(a, b)
}
