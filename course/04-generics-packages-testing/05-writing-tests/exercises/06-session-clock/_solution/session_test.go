package session

import (
	"testing"
	"time"
)

// TestSession drives a Session with a fake clock that the test moves by hand.
func TestSession(t *testing.T) {
	start := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
	now := start
	clock := func() time.Time { return now }

	t.Run("expiry boundary", func(t *testing.T) {
		now = start
		s := New(time.Hour, clock)
		steps := []struct {
			at          time.Duration // time since start
			wantExpired bool
			wantLeft    time.Duration
		}{
			{0, false, time.Hour},
			{time.Hour - time.Nanosecond, false, time.Nanosecond},
			{time.Hour, true, 0},
			{2 * time.Hour, true, 0},
		}
		for _, st := range steps {
			now = start.Add(st.at)
			if got := s.Expired(); got != st.wantExpired {
				t.Errorf("Expired() at +%v = %t, want %t", st.at, got, st.wantExpired)
			}
			if got := s.Remaining(); got != st.wantLeft {
				t.Errorf("Remaining() at +%v = %v, want %v", st.at, got, st.wantLeft)
			}
		}
	})

	t.Run("touch restarts the TTL", func(t *testing.T) {
		now = start
		s := New(time.Hour, clock)
		now = start.Add(40 * time.Minute)
		s.Touch()
		now = start.Add(time.Hour + 39*time.Minute)
		if s.Expired() {
			t.Errorf("Expired() 59m after Touch = true, want false")
		}
		if got, want := s.Remaining(), time.Minute; got != want {
			t.Errorf("Remaining() 59m after Touch = %v, want %v", got, want)
		}
		now = start.Add(time.Hour + 40*time.Minute)
		if !s.Expired() {
			t.Errorf("Expired() 1h after Touch = false, want true")
		}
	})
}
