package gen

import (
	"fmt"
	"slices"
	"testing"
	"time"
)

// collect calls start and reads its channel until it is closed, failing if
// the call or the reads take over 2s.
func collect[T any](t *testing.T, call string, start func() <-chan T) []T {
	t.Helper()
	timeout := time.After(2 * time.Second)
	started := make(chan (<-chan T), 1)
	go func() { started <- start() }()
	var ch <-chan T
	select {
	case ch = <-started:
	case <-timeout:
		t.Fatalf("%s did not return within 2s, want it to return at once and send from a goroutine", call)
	}
	if ch == nil {
		t.Fatalf("%s returned a nil channel", call)
	}
	var got []T
	for {
		select {
		case v, ok := <-ch:
			if !ok {
				return got
			}
			got = append(got, v)
			if len(got) > 10000 {
				t.Fatalf("%s sent over 10000 values, want it to stop", call)
			}
		case <-timeout:
			t.Fatalf("%s: channel not closed within 2s (received %v so far)", call, got)
		}
	}
}

func TestCount(t *testing.T) {
	tests := []struct {
		n    int
		want []int
	}{
		{3, []int{1, 2, 3}},
		{1, []int{1}},
		{0, nil},
		{-5, nil},
		{10, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
	}
	for _, tt := range tests {
		call := fmt.Sprintf("Count(%d)", tt.n)
		got := collect(t, call, func() <-chan int { return Count(tt.n) })
		if !slices.Equal(got, tt.want) {
			t.Errorf("Count(%d) = %v, want %v", tt.n, got, tt.want)
		}
	}
}

func TestCountLarge(t *testing.T) {
	got := collect(t, "Count(5000)", func() <-chan int { return Count(5000) })
	if len(got) != 5000 || got[0] != 1 || got[4999] != 5000 {
		t.Errorf("Count(5000) = %d values, want 1..5000", len(got))
	}
}

func TestRunes(t *testing.T) {
	tests := []struct {
		s    string
		want []rune
	}{
		{"go", []rune{'g', 'o'}},
		{"héllo", []rune{'h', 'é', 'l', 'l', 'o'}},
		{"世界", []rune{'世', '界'}},
		{"", nil},
	}
	for _, tt := range tests {
		call := fmt.Sprintf("Runes(%q)", tt.s)
		got := collect(t, call, func() <-chan rune { return Runes(tt.s) })
		if !slices.Equal(got, tt.want) {
			t.Errorf("Runes(%q) = %q, want %q", tt.s, got, tt.want)
		}
	}
}
