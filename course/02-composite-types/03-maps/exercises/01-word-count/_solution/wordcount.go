package wordcount

import "strings"

// WordCount returns how many times each whitespace-separated word occurs
// in s. Words are lowercased before counting. The result is never nil.
func WordCount(s string) map[string]int {
	counts := map[string]int{}
	for _, w := range strings.Fields(s) {
		counts[strings.ToLower(w)]++ // a missing key reads as 0
	}
	return counts
}
