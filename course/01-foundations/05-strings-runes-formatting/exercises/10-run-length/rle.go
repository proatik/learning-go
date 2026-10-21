package rle

// Encode replaces each run of identical runes in s with the rune followed
// by the run length, e.g. "aaabcc" becomes "a3b1c2".
func Encode(s string) string {
	// TODO
	return ""
}

// Decode reverses Encode. It returns an error if s is not a valid
// encoding (a symbol without a count, a count without a symbol, or a
// count of zero).
func Decode(s string) (string, error) {
	// TODO
	return "", nil
}
