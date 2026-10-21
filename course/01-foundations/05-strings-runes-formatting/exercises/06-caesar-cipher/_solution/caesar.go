package caesar

import "strings"

// Encrypt shifts every ASCII letter in s by shift places, wrapping around
// the alphabet and preserving case. Other runes are left unchanged.
func Encrypt(s string, shift int) string {
	// Bring any shift (negative or > 26) into the range 0–25.
	// Go's % keeps the sign of the dividend, so add 26 before the second %.
	k := rune((shift%26 + 26) % 26)

	var b strings.Builder
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z':
			r = 'a' + (r-'a'+k)%26
		case r >= 'A' && r <= 'Z':
			r = 'A' + (r-'A'+k)%26
		}
		b.WriteRune(r)
	}
	return b.String()
}

// Decrypt reverses Encrypt with the same shift.
func Decrypt(s string, shift int) string {
	return Encrypt(s, -shift)
}
