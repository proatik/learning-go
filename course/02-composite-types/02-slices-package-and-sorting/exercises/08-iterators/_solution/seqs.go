package seqs

import (
	"slices"
	"strings"
)

// Evens returns the even numbers of nums, in their original order.
func Evens(nums []int) []int {
	var out []int
	for v := range slices.Values(nums) { // iter.Seq[int]: one value per step
		if v%2 == 0 {
			out = append(out, v)
		}
	}
	return out
}

// LastIndex returns the index of the last occurrence of w in words, or -1.
func LastIndex(words []string, w string) int {
	for i, v := range slices.Backward(words) { // iter.Seq2[int, string], from the end
		if v == w {
			return i // leaving the loop stops the iterator
		}
	}
	return -1
}

// ChunkSums returns the sum of each consecutive group of n numbers in nums.
// The last group may be shorter. If n < 1 it returns an empty slice.
func ChunkSums(nums []int, n int) []int {
	if n < 1 {
		return nil // slices.Chunk would panic
	}
	var sums []int
	for chunk := range slices.Chunk(nums, n) { // each chunk is a []int sub-slice
		total := 0
		for _, v := range chunk {
			total += v
		}
		sums = append(sums, total)
	}
	return sums
}

// SortedFields returns the whitespace-separated words of s in sorted order.
func SortedFields(s string) []string {
	// FieldsSeq yields the words lazily; Sorted collects them into a sorted slice.
	return slices.Sorted(strings.FieldsSeq(s))
}
