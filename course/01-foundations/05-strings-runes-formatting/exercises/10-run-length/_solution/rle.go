package rle

import (
	"fmt"
	"strconv"
	"strings"
)

// Encode replaces each run of identical runes in s with the rune followed
// by the run length, e.g. "aaabcc" becomes "a3b1c2".
func Encode(s string) string {
	var b strings.Builder
	var prev rune
	count := 0
	for _, r := range s {
		if count > 0 && r == prev {
			count++
			continue
		}
		if count > 0 { // a new rune ends the previous run
			b.WriteRune(prev)
			b.WriteString(strconv.Itoa(count))
		}
		prev, count = r, 1
	}
	if count > 0 { // flush the final run
		b.WriteRune(prev)
		b.WriteString(strconv.Itoa(count))
	}
	return b.String()
}

// isDigit reports whether r is an ASCII digit 0–9.
func isDigit(r rune) bool { return r >= '0' && r <= '9' }

// Decode reverses Encode. It returns an error if s is not a valid
// encoding (a symbol without a count, a count without a symbol, or a
// count of zero).
func Decode(s string) (string, error) {
	r := []rune(s)
	var b strings.Builder
	i := 0
	for i < len(r) {
		sym := r[i]
		if isDigit(sym) {
			return "", fmt.Errorf("count without a symbol at rune %d in %q", i, s)
		}
		i++

		// Read the count: one or more digits.
		n, digits := 0, 0
		for i < len(r) && isDigit(r[i]) {
			n = n*10 + int(r[i]-'0')
			i++
			digits++
		}
		if digits == 0 {
			return "", fmt.Errorf("symbol %q has no count in %q", sym, s)
		}
		if n == 0 {
			return "", fmt.Errorf("symbol %q has count 0 in %q", sym, s)
		}
		b.WriteString(strings.Repeat(string(sym), n))
	}
	return b.String(), nil
}
