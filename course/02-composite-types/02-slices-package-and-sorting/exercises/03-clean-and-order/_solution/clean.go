package clean

import (
	"cmp"
	"slices"
)

// RemoveNegatives returns a new slice with the negative numbers of nums removed,
// keeping the order of the rest. nums is not modified.
func RemoveNegatives(nums []int) []int {
	// DeleteFunc works in place, so give it a copy, never the caller's slice.
	return slices.DeleteFunc(slices.Clone(nums), func(n int) bool { return n < 0 })
}

// Reversed returns a reversed copy of words. words is not modified.
func Reversed(words []string) []string {
	out := slices.Clone(words)
	slices.Reverse(out)
	return out
}

// Descending returns a copy of nums sorted from largest to smallest.
// nums is not modified.
func Descending(nums []int) []int {
	out := slices.Clone(nums)
	// Swapping a and b turns ascending into descending.
	// (Alternative: slices.Sort(out); slices.Reverse(out).)
	slices.SortFunc(out, func(a, b int) int { return cmp.Compare(b, a) })
	return out
}
