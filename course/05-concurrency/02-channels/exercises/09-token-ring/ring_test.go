package ring

import (
	"fmt"
	"runtime"
	"slices"
	"testing"
	"time"
)

// ring calls Ring and fails if it doesn't return within 2s.
func ring(t *testing.T, names []string, turns int) []string {
	t.Helper()
	done := make(chan []string, 1)
	go func() { done <- Ring(names, turns) }()
	select {
	case got := <-done:
		return got
	case <-time.After(2 * time.Second):
		t.Fatalf("Ring(%q, %d) did not return within 2s (is a player blocked, e.g. sending to itself?)", names, turns)
		return nil
	}
}

// settle waits up to 1s for the goroutine count to drop back to before and
// returns how many extra goroutines are left.
func settle(before int) int {
	deadline := time.Now().Add(time.Second)
	for runtime.NumGoroutine() > before && time.Now().Before(deadline) {
		time.Sleep(5 * time.Millisecond)
	}
	return runtime.NumGoroutine() - before
}

func TestRing(t *testing.T) {
	tests := []struct {
		names []string
		turns int
		want  []string
	}{
		{[]string{"ping", "pong"}, 4, []string{"ping", "pong", "ping", "pong"}},
		{[]string{"a", "b", "c"}, 7, []string{"a", "b", "c", "a", "b", "c", "a"}},
		{[]string{"a", "b", "c"}, 2, []string{"a", "b"}},
		{[]string{"a", "b", "c"}, 1, []string{"a"}},
		{[]string{"x", "y"}, 1, []string{"x"}},
		{[]string{"é", "世界"}, 3, []string{"é", "世界", "é"}},
	}
	for _, tt := range tests {
		got := ring(t, tt.names, tt.turns)
		if !slices.Equal(got, tt.want) {
			t.Errorf("Ring(%q, %d) = %q, want %q", tt.names, tt.turns, got, tt.want)
		}
	}
}

func TestRingSolo(t *testing.T) {
	got := ring(t, []string{"solo"}, 3)
	if want := []string{"solo", "solo", "solo"}; !slices.Equal(got, want) {
		t.Errorf("Ring([solo], 3) = %q, want %q", got, want)
	}
}

func TestRingNoTurns(t *testing.T) {
	tests := []struct {
		names []string
		turns int
	}{
		{[]string{"a", "b"}, 0},
		{[]string{"a", "b"}, -2},
		{[]string{"solo"}, 0},
		{nil, 5},
		{[]string{}, 0},
	}
	for _, tt := range tests {
		if got := ring(t, tt.names, tt.turns); len(got) != 0 {
			t.Errorf("Ring(%q, %d) = %q, want []", tt.names, tt.turns, got)
		}
	}
}

func TestRingLarge(t *testing.T) {
	names := make([]string, 100)
	for i := range names {
		names[i] = fmt.Sprint("p", i)
	}
	got := ring(t, names, 10_000)
	if len(got) != 10_000 {
		t.Fatalf("Ring(100 players, 10000) = %d turns, want 10000", len(got))
	}
	for i, name := range got {
		if name != names[i%100] {
			t.Fatalf("Ring(100 players, 10000) turn %d = %q, want %q", i, name, names[i%100])
		}
	}
}

func TestRingNoLeaks(t *testing.T) {
	tests := []struct {
		players, turns int
		want           int
	}{
		{5, 12, 12},
		{1, 4, 4},
		{3, 0, 0},
		{50, 75, 75},
	}
	for _, tt := range tests {
		names := make([]string, tt.players)
		for i := range names {
			names[i] = fmt.Sprint("p", i)
		}
		before := runtime.NumGoroutine()
		got := ring(t, names, tt.turns)
		if len(got) != tt.want {
			t.Fatalf("Ring(%d players, %d) = %d turns, want %d", tt.players, tt.turns, len(got), tt.want)
		}
		if n := settle(before); n > 0 {
			t.Fatalf("Ring(%d players, %d): %d goroutine(s) still running after it returned, want 0 (close the links so every player exits)", tt.players, tt.turns, n)
		}
	}
}
