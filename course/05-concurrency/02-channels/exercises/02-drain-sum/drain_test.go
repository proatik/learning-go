package drain

import (
	"slices"
	"testing"
	"time"
)

// closed returns a buffered channel holding vals, already closed.
func closed[T any](vals ...T) <-chan T {
	ch := make(chan T, len(vals))
	for _, v := range vals {
		ch <- v
	}
	close(ch)
	return ch
}

// slowly returns an unbuffered channel that a goroutine feeds with vals, one
// at a time, and then closes.
func slowly[T any](vals ...T) <-chan T {
	ch := make(chan T)
	go func() {
		defer close(ch)
		for _, v := range vals {
			time.Sleep(time.Millisecond)
			ch <- v
		}
	}()
	return ch
}

// within runs f and fails if it doesn't return within 2s.
func within[R any](t *testing.T, call string, f func() R) R {
	t.Helper()
	done := make(chan R, 1)
	go func() { done <- f() }()
	select {
	case r := <-done:
		return r
	case <-time.After(2 * time.Second):
		t.Fatalf("%s did not return within 2s, want it to stop once the channel is closed", call)
		panic("unreachable")
	}
}

func TestSum(t *testing.T) {
	tests := []struct {
		in   []int
		want int
	}{
		{[]int{1, 2, 3}, 6},
		{[]int{-4, 4, 10}, 10},
		{[]int{0, 0, 0}, 0},
		{[]int{42}, 42},
		{nil, 0},
	}
	for _, tt := range tests {
		got := within(t, "Sum", func() int { return Sum(closed(tt.in...)) })
		if got != tt.want {
			t.Errorf("Sum(%v) = %d, want %d", tt.in, got, tt.want)
		}
	}
}

func TestSumWaitsForSender(t *testing.T) {
	in := []int{5, 0, 7, 0, 3}
	got := within(t, "Sum", func() int { return Sum(slowly(in...)) })
	if got != 15 {
		t.Errorf("Sum(%v sent one at a time) = %d, want 15", in, got)
	}
}

func TestSumMany(t *testing.T) {
	vals := make([]int, 1000)
	for i := range vals {
		vals[i] = i + 1
	}
	got := within(t, "Sum", func() int { return Sum(closed(vals...)) })
	if got != 500500 {
		t.Errorf("Sum(1..1000) = %d, want 500500", got)
	}
}

func TestDrain(t *testing.T) {
	strs := []string{"a", "b", "a", "héllo", ""}
	if got := within(t, "Drain", func() []string { return Drain(closed(strs...)) }); !slices.Equal(got, strs) {
		t.Errorf("Drain(%q) = %q, want %q", strs, got, strs)
	}
	zeros := []int{0, 0}
	if got := within(t, "Drain", func() []int { return Drain(closed(zeros...)) }); !slices.Equal(got, zeros) {
		t.Errorf("Drain(%v) = %v, want %v", zeros, got, zeros)
	}
	if got := within(t, "Drain", func() []float64 { return Drain(closed[float64]()) }); len(got) != 0 {
		t.Errorf("Drain([]) = %v, want []", got)
	}
}

func TestDrainWaitsForSender(t *testing.T) {
	in := []bool{false, true, false}
	got := within(t, "Drain", func() []bool { return Drain(slowly(in...)) })
	if !slices.Equal(got, in) {
		t.Errorf("Drain(%v sent one at a time) = %v, want %v", in, got, in)
	}
}
