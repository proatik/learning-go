package dedupe

// HasDuplicate reports whether any value occurs more than once in nums.
func HasDuplicate(nums []int) bool {
	seen := make(map[int]bool, len(nums)) // size hint: at most len(nums) keys
	for _, n := range nums {
		if seen[n] { // missing key → false
			return true
		}
		seen[n] = true
	}
	return false
}

// Dedupe returns words without duplicates, keeping the first occurrence of
// each word in its original order. It does not modify words.
func Dedupe(words []string) []string {
	seen := make(map[string]struct{}, len(words))
	var out []string
	for _, w := range words {
		if _, ok := seen[w]; ok {
			continue
		}
		seen[w] = struct{}{}
		out = append(out, w)
	}
	return out
}
