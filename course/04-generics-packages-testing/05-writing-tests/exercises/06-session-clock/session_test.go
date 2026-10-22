package session

import "testing"

// TestSession drives a Session with a fake clock that the test moves by hand.
func TestSession(t *testing.T) {
	// TODO: make a fake clock:
	//       now := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
	//       clock := func() time.Time { return now }
	//       s := New(time.Hour, clock)
	// TODO: move time by assigning to now (now = now.Add(...)) and check
	//       Expired and Remaining just before, exactly at, and after the deadline.
	// TODO: check that Touch restarts the TTL from the current fake time.
	_ = New
}
