package wrap

import (
	"strings"
	"unicode/utf8"
)

// Wrap re-flows text into lines of at most width runes, breaking only
// between words. Words longer than width get a line of their own.
func Wrap(text string, width int) string {
	var b strings.Builder
	lineLen := 0 // runes on the line currently being filled

	for _, word := range strings.Fields(text) {
		n := utf8.RuneCountInString(word)
		switch {
		case lineLen == 0:
			// First word of the output: nothing to separate.
		case lineLen+1+n <= width:
			// Fits on the current line after a space.
			b.WriteByte(' ')
			lineLen++
		default:
			// Doesn't fit: start a new line.
			b.WriteByte('\n')
			lineLen = 0
		}
		b.WriteString(word)
		lineLen += n
	}
	return b.String()
}
