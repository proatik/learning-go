package initials

import "testing"

// assertEqual reports a test failure when got != want.
func assertEqual[T comparable](t *testing.T, got, want T) {
	// TODO: mark this function as a helper, so failures point at the caller.
	// TODO: if got != want, report "got %v, want %v" with t.Errorf.
}

// TestInitials checks Initials using the assertEqual helper.
func TestInitials(t *testing.T) {
	// TODO: call assertEqual once per case, e.g.
	//       assertEqual(t, Initials("ada lovelace"), "AL")
	// Cover: two words, one word, extra spaces, lowercase, non-ASCII, empty.
	_ = Initials
}
