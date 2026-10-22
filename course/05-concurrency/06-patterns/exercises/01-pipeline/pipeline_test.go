package pipeline

import (
	"slices"
	"testing"
	"time"
)

// drain reads ch until it is closed, failing if that takes over 2s.
func drain(t *testing.T, name string, ch <-chan int) []int {
	t.Helper()
	if ch == nil {
		t.Fatalf("%s returned a nil channel", name)
	}
	var got []int
	timeout := time.After(2 * time.Second)
	for {
		select {
		case v, ok := <-ch:
			if !ok {
				return got
			}
			got = append(got, v)
		case <-timeout:
			t.Fatalf("%s: channel not closed within 2s (received %v so far)", name, got)
		}
	}
}

// closed returns a closed channel holding vals.
func closed(vals ...int) <-chan int {
	ch := make(chan int, len(vals))
	for _, v := range vals {
		ch <- v
	}
	close(ch)
	return ch
}

// sumWithin runs Sum(in) and fails if it doesn't return within 2s.
func sumWithin(t *testing.T, name string, in <-chan int) int {
	t.Helper()
	done := make(chan int, 1)
	go func() { done <- Sum(in) }()
	select {
	case n := <-done:
		return n
	case <-time.After(2 * time.Second):
		t.Fatalf("%s did not return within 2s", name)
		return 0
	}
}

func TestGen(t *testing.T) {
	tests := []struct {
		name string
		in   []int
	}{
		{"Gen(1, 2, 3)", []int{1, 2, 3}},
		{"Gen(-5, 0, 5, 0)", []int{-5, 0, 5, 0}},
		{"Gen(7)", []int{7}},
		{"Gen()", nil},
	}
	for _, tt := range tests {
		got := drain(t, tt.name, Gen(tt.in...))
		if !slices.Equal(got, tt.in) {
			t.Errorf("%s = %v, want %v", tt.name, got, tt.in)
		}
	}
}

func TestGenDoesNotBlock(t *testing.T) {
	nums := make([]int, 1000)
	done := make(chan (<-chan int), 1)
	go func() { done <- Gen(nums...) }()
	select {
	case ch := <-done:
		if got := drain(t, "Gen(1000 values)", ch); len(got) != 1000 {
			t.Errorf("Gen(1000 values) yielded %d values, want 1000", len(got))
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("Gen(1000 values) did not return within 2s: send from a goroutine, not before returning")
	}
}

func TestSquare(t *testing.T) {
	tests := []struct {
		in   []int
		want []int
	}{
		{[]int{1, 2, 3}, []int{1, 4, 9}},
		{[]int{-4, 0, 10}, []int{16, 0, 100}},
		{nil, nil},
	}
	for _, tt := range tests {
		got := drain(t, "Square", Square(closed(tt.in...)))
		if !slices.Equal(got, tt.want) {
			t.Errorf("Square(%v) = %v, want %v", tt.in, got, tt.want)
		}
	}
}

// TestSquareStreams feeds Square one value at a time and expects each result
// before sending the next: a stage must not wait for its whole input.
func TestSquareStreams(t *testing.T) {
	in := make(chan int)
	out := Square(in)
	if out == nil {
		t.Fatalf("Square returned a nil channel")
	}
	for _, n := range []int{2, 3, 4} {
		select {
		case in <- n:
		case <-time.After(2 * time.Second):
			t.Fatalf("Square did not receive %d within 2s", n)
		}
		select {
		case got := <-out:
			if got != n*n {
				t.Errorf("Square(%d) = %d, want %d", n, got, n*n)
			}
		case <-time.After(2 * time.Second):
			t.Fatalf("Square(%d): no result within 2s (a stage must emit as it goes)", n)
		}
	}
	close(in)
	if got := drain(t, "Square", out); len(got) != 0 {
		t.Errorf("Square after close(in) yielded %v, want nothing", got)
	}
}

func TestSum(t *testing.T) {
	tests := []struct {
		in   []int
		want int
	}{
		{[]int{1, 2, 3}, 6},
		{[]int{-1, 1, -10}, -10},
		{[]int{42}, 42},
		{nil, 0},
	}
	for _, tt := range tests {
		if got := sumWithin(t, "Sum", closed(tt.in...)); got != tt.want {
			t.Errorf("Sum(%v) = %d, want %d", tt.in, got, tt.want)
		}
	}
}

func TestPipeline(t *testing.T) {
	tests := []struct {
		in   []int
		want int
	}{
		{[]int{1, 2, 3}, 14},
		{[]int{-3, 3}, 18},
		{nil, 0},
	}
	for _, tt := range tests {
		sq := Square(Gen(tt.in...))
		if sq == nil {
			t.Fatalf("Square(Gen(%v)) returned a nil channel", tt.in)
		}
		if got := sumWithin(t, "Sum(Square(Gen(...)))", sq); got != tt.want {
			t.Errorf("Sum(Square(Gen(%v))) = %d, want %d", tt.in, got, tt.want)
		}
	}
	nums := make([]int, 100)
	for i := range nums {
		nums[i] = i + 1
	}
	if got := sumWithin(t, "Sum(Square(Gen(1..100)))", Square(Gen(nums...))); got != 338350 {
		t.Errorf("Sum(Square(Gen(1..100))) = %d, want 338350", got)
	}
}
