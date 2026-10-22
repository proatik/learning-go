// Package slug turns titles into URL-friendly slugs.
//
// It lives under internal/, so only code inside this exercise's directory
// tree may import it. We can change it freely without breaking anyone else.
package slug

import (
	"strings"
	"unicode"
)

// Make returns the lowercase slug of title: letters and digits are kept,
// every run of other characters becomes one "-", and the result never
// starts or ends with "-".
func Make(title string) string {
	var b strings.Builder
	dash := false
	for _, r := range title {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			dash = true
			continue
		}
		if dash && b.Len() > 0 {
			b.WriteByte('-')
		}
		dash = false
		b.WriteRune(unicode.ToLower(r))
	}
	return b.String()
}
