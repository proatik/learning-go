package initials

import "testing"

// assertEqual reports a test failure when got != want.
func assertEqual[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

// TestInitials checks Initials using the assertEqual helper.
func TestInitials(t *testing.T) {
	assertEqual(t, Initials("ada lovelace"), "AL")
	assertEqual(t, Initials("Grace Brewster Hopper"), "GBH")
	assertEqual(t, Initials("linus"), "L")
	assertEqual(t, Initials("  ken   thompson  "), "KT")
	assertEqual(t, Initials("émile zola"), "ÉZ")
	assertEqual(t, Initials(""), "")
}
