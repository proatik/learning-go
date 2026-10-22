package leap

import "testing"

// TestIsLeap checks IsLeap with a table of years and the expected answers.
func TestIsLeap(t *testing.T) {
	// TODO: declare a table: a slice of structs with a year and the wanted bool.
	// TODO: loop over it, call IsLeap, and report mismatches with t.Errorf:
	//       "IsLeap(%d) = %t, want %t"
	// Cover each rule: a normal year, a leap year, a century, and a 400th year.
	_ = IsLeap
}
