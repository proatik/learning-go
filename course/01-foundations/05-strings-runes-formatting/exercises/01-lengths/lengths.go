package lengths

// Lengths returns the number of bytes and the number of runes in s.
func Lengths(s string) (bytes, runes int) {
	// TODO: len counts bytes; unicode/utf8 can count runes.
	return 0, 0
}

// CountVowels returns how many of the ASCII vowels a, e, i, o, u
// (in either case) appear in s.
func CountVowels(s string) int {
	// TODO: range over s and count the vowel runes.
	return 0
}
