package rev

import "testing"

// FuzzReverse checks properties of Reverse that hold for every valid string.
func FuzzReverse(f *testing.F) {
	// TODO: seed the corpus with f.Add(...): include "", ASCII words,
	//       a sentence with spaces, and non-ASCII text such as "héllo, 世界".
	f.Fuzz(func(t *testing.T, s string) {
		// TODO: skip input that is not valid UTF-8 (utf8.ValidString + t.Skip).
		// TODO: check the properties from the README:
		//   1. Reverse(Reverse(s)) == s
		//   2. Reverse(s) is valid UTF-8 and has the same length as s
		//   3. the first rune of Reverse(s) is the last rune of s
		_ = Reverse(s)
	})
}
