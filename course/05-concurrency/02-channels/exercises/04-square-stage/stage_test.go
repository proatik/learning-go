package stage

import (
	"fmt"
	"slices"
	"testing"
	"time"
)

// start calls f and fails unless it returns a non-nil channel within 2s.
func start(t *testing.T, call string, f func() <-chan int) <-chan int {
	t.Helper()
	got := make(chan (<-chan int), 1)
	go func() { got <- f() }()
	select {
	case ch := <-got:
		if ch == nil {
			t.Fatalf("%s returned a nil channel", call)
		}
		return ch
	case <-time.After(2 * time.Second):
		t.Fatalf("%s did not return within 2s, want it to return at once and work in a goroutine", call)
		return nil
	}
}

// recv receives one value from ch, failing after 2s.
func recv(t *testing.T, call string, ch <-chan int) (int, bool) {
	t.Helper()
	select {
	case v, ok := <-ch:
		return v, ok
	case <-time.After(2 * time.Second):
		t.Fatalf("%s: no value or close within 2s", call)
		return 0, false
	}
}

// collect reads ch until it is closed.
func collect(t *testing.T, call string, ch <-chan int) []int {
	t.Helper()
	var got []int
	for {
		v, ok := recv(t, call, ch)
		if !ok {
			return got
		}
		got = append(got, v)
		if len(got) > 10000 {
			t.Fatalf("%s sent over 10000 values, want it to stop", call)
		}
	}
}

func TestGenerate(t *testing.T) {
	tests := [][]int{{1, 2, 3}, {0}, {-5, 5}, nil}
	for _, nums := range tests {
		call := fmt.Sprintf("Generate(%v)", nums)
		got := collect(t, call, start(t, call, func() <-chan int { return Generate(nums...) }))
		if !slices.Equal(got, nums) {
			t.Errorf("%s = %v, want %v", call, got, nums)
		}
	}
}

func TestSquare(t *testing.T) {
	tests := []struct {
		in   []int
		want []int
	}{
		{[]int{1, 2, 3}, []int{1, 4, 9}},
		{[]int{-3, 0}, []int{9, 0}},
		{[]int{100}, []int{10000}},
		{nil, nil},
	}
	for _, tt := range tests {
		// A buffered, closed input, so the test doesn't depend on Generate.
		in := make(chan int, len(tt.in))
		for _, v := range tt.in {
			in <- v
		}
		close(in)
		call := fmt.Sprintf("Square(%v)", tt.in)
		got := collect(t, call, start(t, call, func() <-chan int { return Square(in) }))
		if !slices.Equal(got, tt.want) {
			t.Errorf("%s = %v, want %v", call, got, tt.want)
		}
	}
}

func TestPipeline(t *testing.T) {
	call := "Square(Square(Generate(1, 2, 3)))"
	got := collect(t, call, start(t, call, func() <-chan int { return Square(Square(Generate(1, 2, 3))) }))
	if want := []int{1, 16, 81}; !slices.Equal(got, want) {
		t.Errorf("%s = %v, want %v", call, got, want)
	}

	nums := make([]int, 500)
	for i := range nums {
		nums[i] = i
	}
	call = "Square(Generate(0..499))"
	got = collect(t, call, start(t, call, func() <-chan int { return Square(Generate(nums...)) }))
	if len(got) != 500 || got[499] != 499*499 {
		t.Errorf("%s = %d values ending in %v, want 500 values ending in %d", call, len(got), got[max(len(got)-1, 0):], 499*499)
	}
}

func TestSquareStreams(t *testing.T) {
	in := make(chan int) // unbuffered and left open: Square must answer each value as it comes
	out := start(t, "Square(in)", func() <-chan int { return Square(in) })
	for _, v := range []int{3, -4, 7} {
		select {
		case in <- v:
		case <-time.After(2 * time.Second):
			t.Fatalf("Square(in): nobody received %d from in within 2s", v)
		}
		got, ok := recv(t, fmt.Sprintf("Square(in) after sending %d", v), out)
		if !ok || got != v*v {
			t.Fatalf("Square(in) after sending %d (in still open) = (%d, %v), want (%d, true)", v, got, ok, v*v)
		}
	}
	close(in)
	if v, ok := recv(t, "Square(in) after close(in)", out); ok {
		t.Errorf("Square(in) after close(in) = (%d, true), want the output closed", v)
	}
}
