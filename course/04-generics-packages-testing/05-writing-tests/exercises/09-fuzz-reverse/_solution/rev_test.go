package rev

import (
	"testing"
	"unicode/utf8"
)

// FuzzReverse checks properties of Reverse that hold for every valid string.
func FuzzReverse(f *testing.F) {
	for _, s := range []string{"", "a", "hello", "hello world", "héllo, 世界", "Go 🚀"} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) {
		if !utf8.ValidString(s) {
			t.Skip("Reverse is defined on valid UTF-8 only")
		}
		r := Reverse(s)
		if back := Reverse(r); back != s {
			t.Errorf("Reverse(Reverse(%q)) = %q, want %q", s, back, s)
		}
		if !utf8.ValidString(r) || len(r) != len(s) {
			t.Errorf("Reverse(%q) = %q, want valid UTF-8 of length %d", s, r, len(s))
		}
		if s == "" {
			return
		}
		first, _ := utf8.DecodeRuneInString(r)
		last, _ := utf8.DecodeLastRuneInString(s)
		if first != last {
			t.Errorf("Reverse(%q) starts with %q, want %q", s, first, last)
		}
	})
}
