package stack

import (
	"slices"
	"testing"
)

// Compile-time check: *Stack[T] implements Container[T].
var _ Container[int] = (*Stack[int])(nil)

func TestZeroValue(t *testing.T) {
	var s Stack[string]
	if got := s.Len(); got != 0 {
		t.Errorf("empty Stack[string] Len() = %d, want 0", got)
	}
	if v, ok := s.Pop(); v != "" || ok {
		t.Errorf("empty Stack[string] Pop() = (%q, %v), want (\"\", false)", v, ok)
	}
	if v, ok := s.Peek(); v != "" || ok {
		t.Errorf("empty Stack[string] Peek() = (%q, %v), want (\"\", false)", v, ok)
	}
	var p Stack[*int]
	if v, ok := p.Pop(); v != nil || ok {
		t.Errorf("empty Stack[*int] Pop() = (%v, %v), want (<nil>, false)", v, ok)
	}
}

func TestPushPopInts(t *testing.T) {
	tests := []struct {
		push []int
		want []int // pop order
	}{
		{[]int{1}, []int{1}},
		{[]int{1, 2, 3}, []int{3, 2, 1}},
		{[]int{0, -5, 0, 9}, []int{9, 0, -5, 0}},
	}
	for _, tt := range tests {
		var s Stack[int]
		for _, v := range tt.push {
			s.Push(v)
		}
		if got := s.Len(); got != len(tt.push) {
			t.Errorf("Push %v: Len() = %d, want %d", tt.push, got, len(tt.push))
		}
		var got []int
		for range tt.push {
			v, ok := s.Pop()
			if !ok {
				break
			}
			got = append(got, v)
		}
		if !slices.Equal(got, tt.want) {
			t.Errorf("Push %v then Pop all = %v, want %v", tt.push, got, tt.want)
		}
		if v, ok := s.Pop(); v != 0 || ok {
			t.Errorf("Push %v: extra Pop() = (%d, %v), want (0, false)", tt.push, v, ok)
		}
	}
}

func TestPeekStrings(t *testing.T) {
	var s Stack[string]
	s.Push("héllo")
	s.Push("世界")
	if v, ok := s.Peek(); v != "世界" || !ok {
		t.Errorf("Peek() = (%q, %v), want (%q, true)", v, ok, "世界")
	}
	if got := s.Len(); got != 2 {
		t.Errorf("Peek() changed the stack: Len() = %d, want 2", got)
	}
	if v, ok := s.Pop(); v != "世界" || !ok {
		t.Errorf("Pop() = (%q, %v), want (%q, true)", v, ok, "世界")
	}
	if v, ok := s.Peek(); v != "héllo" || !ok {
		t.Errorf("Peek() after Pop = (%q, %v), want (%q, true)", v, ok, "héllo")
	}
}

func TestStructsAndSlices(t *testing.T) {
	type point struct{ X, Y int }
	var ps Stack[point]
	ps.Push(point{1, 2})
	if v, ok := ps.Pop(); v != (point{1, 2}) || !ok {
		t.Errorf("Stack[point] Pop() = (%v, %v), want ({1 2}, true)", v, ok)
	}
	var ss Stack[[]int]
	ss.Push([]int{1, 2})
	ss.Push(nil)
	if v, ok := ss.Pop(); v != nil || !ok {
		t.Errorf("Stack[[]int] Pop() = (%v, %v), want ([], true)", v, ok)
	}
	if v, ok := ss.Pop(); !slices.Equal(v, []int{1, 2}) || !ok {
		t.Errorf("Stack[[]int] Pop() = (%v, %v), want ([1 2], true)", v, ok)
	}
}

func TestReuseAfterEmpty(t *testing.T) {
	var s Stack[int]
	s.Push(1)
	s.Pop()
	s.Pop()
	s.Push(7)
	if v, ok := s.Pop(); v != 7 || !ok {
		t.Errorf("Push(7) after emptying, Pop() = (%d, %v), want (7, true)", v, ok)
	}
}

func TestDrainStack(t *testing.T) {
	s := &Stack[int]{}
	for _, v := range []int{1, 2, 3} {
		s.Push(v)
	}
	if got, want := Drain[int](s), []int{3, 2, 1}; !slices.Equal(got, want) {
		t.Errorf("Drain(stack of [1 2 3]) = %v, want %v", got, want)
	}
	if got := s.Len(); got != 0 {
		t.Errorf("after Drain, Len() = %d, want 0", got)
	}
	if got := Drain[int](s); len(got) != 0 {
		t.Errorf("Drain(empty stack) = %v, want []", got)
	}
}

// fifo is a different Container: a first-in, first-out queue.
type fifo struct{ items []string }

func (q *fifo) Push(v string) { q.items = append(q.items, v) }
func (q *fifo) Len() int      { return len(q.items) }
func (q *fifo) Pop() (string, bool) {
	if len(q.items) == 0 {
		return "", false
	}
	v := q.items[0]
	q.items = q.items[1:]
	return v, true
}

func TestDrainAnyContainer(t *testing.T) {
	q := &fifo{items: []string{"a", "b", "c"}}
	if got, want := Drain[string](q), []string{"a", "b", "c"}; !slices.Equal(got, want) {
		t.Errorf("Drain(fifo of [a b c]) = %q, want %q", got, want)
	}
	if got := q.Len(); got != 0 {
		t.Errorf("after Drain, fifo Len() = %d, want 0", got)
	}
}
