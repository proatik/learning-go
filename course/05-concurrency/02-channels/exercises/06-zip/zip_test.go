package zip

import (
	"fmt"
	"slices"
	"testing"
	"time"
)

// closed returns a buffered channel holding vals, already closed. Its len
// shows how many values were not read.
func closed[T any](vals []T) chan T {
	ch := make(chan T, len(vals))
	for _, v := range vals {
		ch <- v
	}
	close(ch)
	return ch
}

// zipAll calls Zip(a, b) and reads the result until it is closed, failing if
// anything takes over 2s.
func zipAll[X, Y any](t *testing.T, call string, a <-chan X, b <-chan Y) []Pair[X, Y] {
	t.Helper()
	timeout := time.After(2 * time.Second)
	started := make(chan (<-chan Pair[X, Y]), 1)
	go func() { started <- Zip(a, b) }()
	var ch <-chan Pair[X, Y]
	select {
	case ch = <-started:
	case <-timeout:
		t.Fatalf("%s did not return within 2s, want it to return at once and work in a goroutine", call)
	}
	if ch == nil {
		t.Fatalf("%s returned a nil channel", call)
	}
	var got []Pair[X, Y]
	for {
		select {
		case p, ok := <-ch:
			if !ok {
				return got
			}
			got = append(got, p)
			if len(got) > 10000 {
				t.Fatalf("%s sent over 10000 pairs, want it to stop when an input is closed", call)
			}
		case <-timeout:
			t.Fatalf("%s: output not closed within 2s (received %v so far)", call, got)
		}
	}
}

type ps = Pair[int, string]

func TestZip(t *testing.T) {
	tests := []struct {
		a            []int
		b            []string
		want         []ps
		leftA, leftB int // values Zip must leave unread
	}{
		{[]int{1, 2, 3}, []string{"x", "y", "z"}, []ps{{1, "x"}, {2, "y"}, {3, "z"}}, 0, 0},
		{[]int{1, 2, 3}, []string{"x"}, []ps{{1, "x"}}, 1, 0},
		{[]int{1}, []string{"x", "y", "z"}, []ps{{1, "x"}}, 0, 2},
		{[]int{0, 0}, []string{"", ""}, []ps{{0, ""}, {0, ""}}, 0, 0},
		{nil, []string{"x"}, nil, 0, 1},
		{[]int{7}, nil, nil, 0, 0},
		{nil, nil, nil, 0, 0},
		{[]int{1, 2}, []string{"é", "世界", "!"}, []ps{{1, "é"}, {2, "世界"}}, 0, 1},
	}
	for _, tt := range tests {
		a, b := closed(tt.a), closed(tt.b)
		call := fmt.Sprintf("Zip(%v, %q)", tt.a, tt.b)
		got := zipAll[int, string](t, call, a, b)
		if !slices.Equal(got, tt.want) {
			t.Errorf("%s = %v, want %v", call, got, tt.want)
			continue
		}
		if len(a) != tt.leftA || len(b) != tt.leftB {
			t.Errorf("%s left (%d in a, %d in b) unread, want (%d, %d): read a first, then b, and stop at the first closed one", call, len(a), len(b), tt.leftA, tt.leftB)
		}
	}
}

func TestZipStreams(t *testing.T) {
	a := make(chan int)
	b := make(chan bool)
	go func() {
		defer close(a)
		for i := range 300 {
			a <- i
		}
	}()
	go func() {
		defer close(b)
		for i := range 200 {
			b <- i%2 == 0
		}
	}()
	got := zipAll[int, bool](t, "Zip(0..299, 200 bools)", a, b)
	if len(got) != 200 {
		t.Fatalf("Zip(0..299, 200 bools) = %d pairs, want 200", len(got))
	}
	for i, p := range got {
		if want := (Pair[int, bool]{i, i%2 == 0}); p != want {
			t.Fatalf("Zip(0..299, 200 bools)[%d] = %v, want %v", i, p, want)
		}
	}
}
