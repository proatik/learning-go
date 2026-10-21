package canon

import "slices"

// Unique returns the distinct values of nums in ascending order.
// nums is not modified.
func Unique(nums []int) []int {
	out := slices.Clone(nums)
	slices.Sort(out)           // equal values become neighbours...
	return slices.Compact(out) // ...so Compact can collapse them
}

// CountDistinct returns the number of different words in words.
// words is not modified.
func CountDistinct(words []string) int {
	sorted := slices.Clone(words)
	slices.Sort(sorted)
	return len(slices.Compact(sorted))
}

// SortedLetters returns the runes of s sorted in ascending order.
func SortedLetters(s string) string {
	r := []rune(s) // a new slice: sorting it doesn't touch s (strings are immutable anyway)
	slices.Sort(r)
	return string(r)
}

// IsAnagram reports whether a and b consist of the same runes,
// each appearing the same number of times.
func IsAnagram(a, b string) bool {
	return SortedLetters(a) == SortedLetters(b)
}
