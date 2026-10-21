package reverse

import "strings"

// Reverse returns s with its runes in reverse order.
func Reverse(s string) string {
	r := []rune(s) // one element per rune, so multi-byte characters stay intact
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

// ReverseWords returns the whitespace-separated words of s in reverse
// order, joined by single spaces.
func ReverseWords(s string) string {
	words := strings.Fields(s)

	var b strings.Builder
	for i := len(words) - 1; i >= 0; i-- {
		if b.Len() > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(words[i])
	}
	return b.String()
}
