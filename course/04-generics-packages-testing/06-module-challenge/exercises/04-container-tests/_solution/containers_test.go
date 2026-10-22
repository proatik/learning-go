package containers

import (
	"slices"
	"testing"
)

// op is one step of a scripted test: push a value, or pop and expect a
// result.
type op struct {
	push   int  // value to push when pop is false
	pop    bool // pop instead of push
	want   int  // expected popped value
	wantOK bool // expected ok from Pop
}

func push(v int) op   { return op{push: v} }
func pop(want int) op { return op{pop: true, want: want, wantOK: true} }
func popEmpty() op    { return op{pop: true} }
func pushAll(vs ...int) []op {
	ops := make([]op, len(vs))
	for i, v := range vs {
		ops[i] = push(v)
	}
	return ops
}

// popper is what both containers share, so one helper can drive either.
type popper interface {
	Push(int)
	Pop() (int, bool)
	Len() int
}

// runOps applies ops to c, checking every Pop and the final Len.
func runOps(t *testing.T, c popper, ops []op, wantLen int) {
	t.Helper()
	for i, o := range ops {
		if !o.pop {
			c.Push(o.push)
			continue
		}
		got, ok := c.Pop()
		if got != o.want || ok != o.wantOK {
			t.Fatalf("step %d: Pop() = %d, %v, want %d, %v", i, got, ok, o.want, o.wantOK)
		}
	}
	if got := c.Len(); got != wantLen {
		t.Errorf("Len() = %d, want %d", got, wantLen)
	}
}

func concat(parts ...[]op) []op { return slices.Concat(parts...) }

func TestStack(t *testing.T) {
	tests := []struct {
		name    string
		ops     []op
		wantLen int
		wantAll []int
	}{
		{"empty", []op{popEmpty()}, 0, nil},
		{"last in first out", concat(pushAll(1, 2, 3), []op{pop(3), pop(2)}), 1, []int{1}},
		{"pop until empty", concat(pushAll(1, 2), []op{pop(2), pop(1), popEmpty()}), 0, nil},
		{"all is top to bottom", pushAll(1, 2, 3), 3, []int{3, 2, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewStack[int]()
			runOps(t, s, tt.ops, tt.wantLen)
			if got := slices.Collect(s.All()); !slices.Equal(got, tt.wantAll) {
				t.Errorf("All() = %v, want %v", got, tt.wantAll)
			}
		})
	}

	t.Run("peek", func(t *testing.T) {
		s := NewStack[string]()
		if v, ok := s.Peek(); v != "" || ok {
			t.Errorf("Peek() on empty = %q, %v, want %q, false", v, ok, "")
		}
		s.Push("a")
		s.Push("b")
		for range 2 {
			if v, ok := s.Peek(); v != "b" || !ok {
				t.Errorf("Peek() = %q, %v, want %q, true", v, ok, "b")
			}
		}
		if n := s.Len(); n != 2 {
			t.Errorf("Len() after Peek = %d, want 2", n)
		}
	})

	t.Run("all stops on break", func(t *testing.T) {
		s := NewStack[int]()
		for _, v := range []int{1, 2, 3} {
			s.Push(v)
		}
		var got []int
		for v := range s.All() {
			got = append(got, v)
			break
		}
		if want := []int{3}; !slices.Equal(got, want) {
			t.Errorf("All() with break = %v, want %v", got, want)
		}
	})
}

func TestQueue(t *testing.T) {
	tests := []struct {
		name    string
		ops     []op
		wantLen int
		wantAll []int
	}{
		{"empty", []op{popEmpty()}, 0, nil},
		{"first in first out", concat(pushAll(1, 2, 3), []op{pop(1), pop(2)}), 1, []int{3}},
		{"empty after pops returns zero", concat(pushAll(7, 8), []op{pop(7), pop(8), popEmpty()}), 0, nil},
		{"wraps around", concat(pushAll(1, 2, 3), []op{pop(1), pop(2)}, pushAll(4, 5, 6)), 4, []int{3, 4, 5, 6}},
		{
			"grows after wrapping",
			concat(pushAll(1, 2, 3, 4), []op{pop(1), pop(2)}, pushAll(5, 6, 7, 8, 9)),
			7, []int{3, 4, 5, 6, 7, 8, 9},
		},
		{"grows from the start", pushAll(1, 2, 3, 4, 5, 6, 7, 8, 9), 9, []int{1, 2, 3, 4, 5, 6, 7, 8, 9}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := NewQueue[int]()
			runOps(t, q, tt.ops, tt.wantLen)
			if got := slices.Collect(q.All()); !slices.Equal(got, tt.wantAll) {
				t.Errorf("All() = %v, want %v", got, tt.wantAll)
			}
		})
	}
}
