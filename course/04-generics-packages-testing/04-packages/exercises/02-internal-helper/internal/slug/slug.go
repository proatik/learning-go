// Package slug turns titles into URL-friendly slugs.
//
// It lives under internal/, so only code inside this exercise's directory
// tree may import it. We can change it freely without breaking anyone else.
package slug

// Make returns the lowercase slug of title: letters and digits are kept,
// every run of other characters becomes one "-", and the result never
// starts or ends with "-".
func Make(title string) string {
	// TODO: walk the runes; remember "a dash is pending" instead of writing it
	// immediately, and only write it before the next kept rune.
	return ""
}
