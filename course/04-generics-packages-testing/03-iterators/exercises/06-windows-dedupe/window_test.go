package window

import (
	"fmt"
	"iter"
	"slices"
	"testing"
)

// collect ranges over seq and returns its values, breaking after max values
// when max > 0 (and always after 10000, so a runaway sequence can't hang the
// test). A panic during the loop is returned instead of crashing.
func collect[V any](seq iter.Seq[V], max int) (got []V, panicked any) {
	defer func() { panicked = recover() }()
	for v := range seq {
		got = append(got, v)
		if len(got) == max || len(got) == 10000 {
			break
		}
	}
	return got, nil
}

func TestWindows(t *testing.T) {
	tests := []struct {
		s    []int
		k    int
		want string
	}{
		{[]int{1, 2, 3, 4}, 2, "[[1 2] [2 3] [3 4]]"},
		{[]int{1, 2, 3, 4}, 1, "[[1] [2] [3] [4]]"},
		{[]int{1, 2, 3}, 3, "[[1 2 3]]"},
		{[]int{1, 2, 3, 4, 5}, 4, "[[1 2 3 4] [2 3 4 5]]"},
		{[]int{1, 2}, 3, "[]"},
		{[]int{1, 2}, 0, "[]"},
		{[]int{1, 2}, -1, "[]"},
		{nil, 1, "[]"},
		{nil, 0, "[]"},
	}
	for _, tt := range tests {
		seq := Windows(tt.s, tt.k)
		if seq == nil {
			t.Fatalf("Windows(%v, %d) = nil, want a non-nil iterator", tt.s, tt.k)
		}
		got, p := collect(seq, -1)
		if p != nil {
			t.Fatalf("range over Windows(%v, %d) panicked: %v", tt.s, tt.k, p)
		}
		if s := fmt.Sprint(got); s != tt.want && !(tt.want == "[]" && got == nil) {
			t.Errorf("Windows(%v, %d) = %s, want %s", tt.s, tt.k, s, tt.want)
		}
	}
}

func TestWindowsStrings(t *testing.T) {
	s := []string{"日", "本", "語"}
	got, _ := collect(Windows(s, 2), -1)
	if want := "[[日 本] [本 語]]"; fmt.Sprint(got) != want {
		t.Errorf("Windows(%q, 2) = %q, want %s", s, got, want)
	}
}

func TestWindowsShareAndClip(t *testing.T) {
	s := []int{10, 20, 30, 40}
	got, _ := collect(Windows(s, 2), -1)
	if len(got) != 3 {
		t.Fatalf("Windows(%v, 2) yielded %d windows, want 3", s, len(got))
	}
	for i, w := range got {
		if len(w) != 2 || &w[0] != &s[i] {
			t.Errorf("Windows(%v, 2) window %d shares memory with s = false, want true (slice s, don't copy)", s, i)
		}
		if cap(w) != 2 {
			t.Errorf("Windows(%v, 2) window %d cap = %d, want 2 (use s[i : i+k : i+k])", s, i, cap(w))
		}
	}
	_ = append(got[0], 99)
	if s[2] != 30 {
		t.Errorf("append(first window, 99) changed s to %v, want [10 20 30 40]", s)
	}
}

func TestWindowsBreak(t *testing.T) {
	s := make([]int, 1000)
	calls := 0
	Windows(s, 3)(func([]int) bool {
		calls++
		return calls < 2
	})
	if calls != 2 {
		t.Errorf("Windows(1000 items, 3) with yield returning false on call 2: yield calls = %d, want 2", calls)
	}
	got, p := collect(Windows(s, 3), 4)
	if p != nil {
		t.Fatalf("break after 4 windows of Windows(1000 items, 3) panicked: %v (stop when yield returns false)", p)
	}
	if len(got) != 4 {
		t.Errorf("windows received before a break at 4 from Windows(1000 items, 3) = %d, want 4", len(got))
	}
}

func TestDedupe(t *testing.T) {
	tests := []struct {
		in   []int
		want []int
	}{
		{[]int{1, 1, 2, 2, 2, 1, 3, 3}, []int{1, 2, 1, 3}},
		{[]int{0, 0, 1}, []int{0, 1}},
		{[]int{0}, []int{0}},
		{[]int{5, 5, 5}, []int{5}},
		{[]int{1, 2, 3}, []int{1, 2, 3}},
		{[]int{1, 2, 1, 2}, []int{1, 2, 1, 2}},
		{nil, nil},
	}
	for _, tt := range tests {
		seq := Dedupe(slices.Values(tt.in))
		if seq == nil {
			t.Fatalf("Dedupe(%v) = nil, want a non-nil iterator", tt.in)
		}
		got, p := collect(seq, -1)
		if p != nil {
			t.Fatalf("range over Dedupe(%v) panicked: %v", tt.in, p)
		}
		if !slices.Equal(got, tt.want) {
			t.Errorf("Dedupe(%v) = %v, want %v", tt.in, got, tt.want)
		}
	}
	words, _ := collect(Dedupe(slices.Values([]string{"", "", "a", "a", "é", "e"})), -1)
	if want := []string{"", "a", "é", "e"}; !slices.Equal(words, want) {
		t.Errorf("Dedupe(%q) = %q, want %q", []string{"", "", "a", "a", "é", "e"}, words, want)
	}
}

func TestDedupeRunsAgain(t *testing.T) {
	seq := Dedupe(slices.Values([]int{7, 7, 8}))
	first, _ := collect(seq, -1)
	second, _ := collect(seq, -1)
	if !slices.Equal(first, []int{7, 8}) || !slices.Equal(second, []int{7, 8}) {
		t.Errorf("Dedupe([7 7 8]) ranged twice = %v then %v, want [7 8] both times (keep state inside the returned function)", first, second)
	}
}

func TestDedupeBreak(t *testing.T) {
	pulled := 0
	src := func(yield func(int) bool) {
		for i := 0; ; i++ {
			pulled++
			if !yield(i / 3) { // 0 0 0 1 1 1 2 2 2 …
				return
			}
		}
	}
	got, p := collect(Dedupe(src), 2)
	if p != nil {
		t.Fatalf("break after 2 values of Dedupe(0 0 0 1 1 1 …) panicked: %v (stop when yield returns false)", p)
	}
	if !slices.Equal(got, []int{0, 1}) {
		t.Errorf("first 2 values of Dedupe(0 0 0 1 1 1 …) = %v, want [0 1]", got)
	}
	if pulled != 4 {
		t.Errorf("break after 2 values of Dedupe(0 0 0 1 1 1 …): source values produced = %d, want 4", pulled)
	}
}
