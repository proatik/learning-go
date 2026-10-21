package acronym

import (
	"strings"
	"unicode"
)

// Acronym returns the upper-cased first letter of each word in phrase.
// Words are separated by whitespace or hyphens.
func Acronym(phrase string) string {
	// Treat hyphens like spaces, then let Fields split on any whitespace run.
	words := strings.Fields(strings.ReplaceAll(phrase, "-", " "))

	var b strings.Builder
	for _, word := range words {
		for _, r := range word {
			if unicode.IsLetter(r) {
				b.WriteRune(unicode.ToUpper(r))
				break // only the first letter of each word
			}
		}
	}
	return b.String()
}
