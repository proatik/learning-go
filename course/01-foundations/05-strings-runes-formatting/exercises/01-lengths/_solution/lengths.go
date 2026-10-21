package lengths

import "unicode/utf8"

// Lengths returns the number of bytes and the number of runes in s.
func Lengths(s string) (bytes, runes int) {
	// len counts bytes; RuneCountInString decodes the UTF-8 and counts runes.
	return len(s), utf8.RuneCountInString(s)
}

// CountVowels returns how many of the ASCII vowels a, e, i, o, u
// (in either case) appear in s.
func CountVowels(s string) int {
	count := 0
	for _, r := range s { // r is a rune, so 'é' or '世' is one step
		switch r {
		case 'a', 'e', 'i', 'o', 'u', 'A', 'E', 'I', 'O', 'U':
			count++
		}
	}
	return count
}
