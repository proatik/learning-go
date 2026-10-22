package port

import "testing"

// TestParsePort checks valid ports and every error path of ParsePort.
func TestParsePort(t *testing.T) {
	// TODO: a table of {in string; want int; wantErr error} cases.
	//       wantErr is nil for valid ports, or ErrEmpty / ErrNotNumber / ErrOutOfRange.
	// TODO: for each case, check the error with errors.Is (the errors are wrapped,
	//       so == will not work), and check the number too.
	// Report:
	//   "ParsePort(%q) error = %v, want %v"
	//   "ParsePort(%q) = %d, want %d"
	_ = ParsePort
}
