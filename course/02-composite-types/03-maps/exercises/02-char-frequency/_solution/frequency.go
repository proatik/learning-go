package frequency

import "unicode"

// CharFrequency counts each non-whitespace rune in s (case-sensitive).
// The result is never nil.
func CharFrequency(s string) map[rune]int {
	freq := map[rune]int{}
	for _, r := range s { // range over a string decodes runes
		if unicode.IsSpace(r) {
			continue
		}
		freq[r]++
	}
	return freq
}

// FirstUnique returns the first rune of s that occurs exactly once in s.
// ok is false if every rune repeats (or s is empty).
func FirstUnique(s string) (r rune, ok bool) {
	counts := map[rune]int{}
	for _, r := range s {
		counts[r]++
	}
	// Walk s again: the string has an order, the map does not.
	for _, r := range s {
		if counts[r] == 1 {
			return r, true
		}
	}
	return 0, false
}
