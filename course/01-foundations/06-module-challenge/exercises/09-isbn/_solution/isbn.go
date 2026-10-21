package isbn

import (
	"fmt"
	"strings"
)

// IsValid reports whether isbn is a valid ISBN-10. Hyphens are ignored.
func IsValid(isbn string) bool {
	s := strings.ReplaceAll(isbn, "-", "")
	if len(s) != 10 {
		return false
	}

	sum := 0
	for i := range len(s) {
		c := s[i]
		var d int
		switch {
		case '0' <= c && c <= '9':
			d = int(c - '0')
		case c == 'X' && i == 9:
			d = 10
		default:
			return false
		}
		sum += d * (10 - i) // weights 10, 9, ..., 1
	}
	return sum%11 == 0
}

// ToISBN13 converts a valid ISBN-10 to its 13-digit ISBN-13 form
// (978 prefix, no hyphens).
func ToISBN13(isbn string) (string, error) {
	if !IsValid(isbn) {
		return "", fmt.Errorf("isbn: %q is not a valid ISBN-10", isbn)
	}

	// Drop the old check character; IsValid guarantees 9 digits remain.
	body := "978" + strings.ReplaceAll(isbn, "-", "")[:9]

	sum := 0
	for i := range len(body) {
		d := int(body[i] - '0')
		if i%2 == 1 {
			d *= 3
		}
		sum += d
	}
	check := (10 - sum%10) % 10
	return body + string(rune('0'+check)), nil
}
