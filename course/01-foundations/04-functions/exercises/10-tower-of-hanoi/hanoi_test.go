package hanoi

import (
	"fmt"
	"testing"
)

func TestHanoiMoveLog(t *testing.T) {
	tests := []struct {
		n    int
		want string
	}{
		{1, "A->C "},
		{2, "A->B A->C B->C "},
		{3, "A->C A->B C->B A->C B->A B->C A->C "},
	}
	for _, tt := range tests {
		log := ""
		count := Hanoi(tt.n, "A", "C", "B", func(from, to string) {
			log += from + "->" + to + " "
		})
		if log != tt.want {
			t.Errorf("Hanoi(%d, \"A\", \"C\", \"B\", f) moves = %q, want %q", tt.n, log, tt.want)
		}
		if want := (1 << tt.n) - 1; count != want {
			t.Errorf("Hanoi(%d, ...) = %d, want %d", tt.n, count, want)
		}
	}
}

// simulate replays the moves on real pegs and reports the first illegal one.
func simulate(t *testing.T, n int, from, to, via string) {
	t.Helper()
	pegs := map[string][]int{from: nil, to: nil, via: nil}
	for d := n; d >= 1; d-- {
		pegs[from] = append(pegs[from], d) // largest disk at the bottom
	}
	calls, bad := 0, ""
	count := Hanoi(n, from, to, via, func(src, dst string) {
		calls++
		if bad != "" {
			return
		}
		s, ok1 := pegs[src]
		_, ok2 := pegs[dst]
		switch {
		case !ok1 || !ok2:
			bad = fmt.Sprintf("move %d: unknown peg in %s->%s", calls, src, dst)
		case len(s) == 0:
			bad = fmt.Sprintf("move %d: %s->%s from an empty peg", calls, src, dst)
		default:
			disk := s[len(s)-1]
			if d := pegs[dst]; len(d) > 0 && d[len(d)-1] < disk {
				bad = fmt.Sprintf("move %d: disk %d placed on smaller disk %d (%s->%s)", calls, disk, d[len(d)-1], src, dst)
				return
			}
			pegs[src] = s[:len(s)-1]
			pegs[dst] = append(pegs[dst], disk)
		}
	})
	if bad != "" {
		t.Fatalf("Hanoi(%d, %q, %q, %q): illegal %s", n, from, to, via, bad)
	}
	want := (1 << n) - 1
	if count != want || calls != want {
		t.Errorf("Hanoi(%d, ...) returned %d and called move %d times, want %d", n, count, calls, want)
	}
	if len(pegs[to]) != n {
		t.Errorf("Hanoi(%d, %q, %q, %q): %d disks ended on %q, want %d", n, from, to, via, len(pegs[to]), to, n)
	}
}

func TestHanoiLegal(t *testing.T) {
	for n := 1; n <= 10; n++ {
		simulate(t, n, "A", "C", "B")
	}
	simulate(t, 5, "left", "right", "middle")
	simulate(t, 4, "X", "Y", "Z")
}

func TestHanoiZero(t *testing.T) {
	for _, n := range []int{0, -3} {
		calls := 0
		count := Hanoi(n, "A", "C", "B", func(string, string) { calls++ })
		if count != 0 || calls != 0 {
			t.Errorf("Hanoi(%d, ...) returned %d with %d moves, want 0 and 0", n, count, calls)
		}
	}
}

func TestHanoiNilMove(t *testing.T) {
	tests := []struct{ n, want int }{
		{1, 1},
		{3, 7},
		{10, 1023},
		{20, 1048575},
	}
	for _, tt := range tests {
		if got := Hanoi(tt.n, "A", "C", "B", nil); got != tt.want {
			t.Errorf("Hanoi(%d, \"A\", \"C\", \"B\", nil) = %d, want %d", tt.n, got, tt.want)
		}
	}
}
