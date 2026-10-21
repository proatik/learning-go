package palindrome

import (
	"strings"
	"unicode"
)

// IsPalindrome reports whether s reads the same forwards and backwards,
// ignoring case and every rune that is not a letter or digit.
func IsPalindrome(s string) bool {
	// Step 1: keep only lower-cased letters and digits.
	var b strings.Builder
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(unicode.ToLower(r))
		}
	}

	// Step 2: compare runes from both ends toward the middle.
	r := []rune(b.String())
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		if r[i] != r[j] {
			return false
		}
	}
	return true
}
