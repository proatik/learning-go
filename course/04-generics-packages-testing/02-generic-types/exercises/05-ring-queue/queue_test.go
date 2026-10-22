package queue

import (
	"slices"
	"testing"
)

func TestZeroValue(t *testing.T) {
	var q Queue[string]
	if got := q.Len(); got != 0 {
		t.Errorf("empty Queue Len() = %d, want 0", got)
	}
	if v, ok := q.Pop(); v != "" || ok {
		t.Errorf("empty Queue Pop() = (%q, %v), want (\"\", false)", v, ok)
	}
	if v, ok := q.Peek(); v != "" || ok {
		t.Errorf("empty Queue Peek() = (%q, %v), want (\"\", false)", v, ok)
	}
	if got := q.Values(); len(got) != 0 {
		t.Errorf("empty Queue Values() = %q, want []", got)
	}
}

func TestFIFO(t *testing.T) {
	tests := []struct {
		push []int
	}{
		{[]int{1}},
		{[]int{1, 2, 3}},
		{[]int{5, -3, 0, 9, 9, 2, 7, 1, 8, 4}},
	}
	for _, tt := range tests {
		var q Queue[int]
		for _, v := range tt.push {
			q.Push(v)
		}
		if got := q.Values(); !slices.Equal(got, tt.push) {
			t.Errorf("Push %v: Values() = %v, want %v", tt.push, got, tt.push)
		}
		if got := q.Len(); got != len(tt.push) {
			t.Errorf("Push %v: Len() = %d, want %d", tt.push, got, len(tt.push))
		}
		if v, ok := q.Peek(); v != tt.push[0] || !ok {
			t.Errorf("Push %v: Peek() = (%d, %v), want (%d, true)", tt.push, v, ok, tt.push[0])
		}
		var got []int
		for range tt.push {
			v, ok := q.Pop()
			if !ok {
				break
			}
			got = append(got, v)
		}
		if !slices.Equal(got, tt.push) {
			t.Errorf("Push %v then Pop all = %v, want %v", tt.push, got, tt.push)
		}
		if v, ok := q.Pop(); v != 0 || ok {
			t.Errorf("Push %v: extra Pop() = (%d, %v), want (0, false)", tt.push, v, ok)
		}
	}
}

func TestGrowth(t *testing.T) {
	tests := []struct {
		pushes int
		want   int // len(q.buf)
	}{
		{1, 4},
		{4, 4},
		{5, 8},
		{8, 8},
		{9, 16},
		{17, 32},
	}
	for _, tt := range tests {
		var q Queue[int]
		for i := range tt.pushes {
			q.Push(i)
		}
		if got := len(q.buf); got != tt.want {
			t.Errorf("after %d pushes, len(buf) = %d, want %d", tt.pushes, got, tt.want)
		}
	}
}

func TestWrapAround(t *testing.T) {
	var q Queue[int]
	for _, v := range []int{1, 2, 3, 4} {
		q.Push(v)
	}
	q.Pop()
	q.Pop()
	q.Push(5) // wraps to the start of buf
	q.Push(6)
	if got, want := q.Values(), []int{3, 4, 5, 6}; !slices.Equal(got, want) {
		t.Errorf("after wrapping, Values() = %v, want %v", got, want)
	}
	if got := len(q.buf); got != 4 {
		t.Errorf("after wrapping without being full, len(buf) = %d, want 4", got)
	}
	q.Push(7) // full: grows and unwraps
	if got, want := q.Values(), []int{3, 4, 5, 6, 7}; !slices.Equal(got, want) {
		t.Errorf("after growing a wrapped buffer, Values() = %v, want %v", got, want)
	}
	if got := len(q.buf); got != 8 {
		t.Errorf("after growing, len(buf) = %d, want 8", got)
	}
	if q.head != 0 || len(q.buf) < 5 || !slices.Equal(q.buf[:5], []int{3, 4, 5, 6, 7}) {
		t.Errorf("after growing, buf = %v (head %d), want [3 4 5 6 7 0 0 0] (head 0)", q.buf, q.head)
	}
}

func TestSteadyStateDoesNotGrow(t *testing.T) {
	var q Queue[int]
	q.Push(0)
	q.Push(1)
	q.Push(2)
	for i := 3; i < 1000; i++ {
		q.Push(i)
		if v, ok := q.Pop(); v != i-3 || !ok {
			t.Fatalf("cycle %d: Pop() = (%d, %v), want (%d, true)", i, v, ok, i-3)
		}
	}
	if got := len(q.buf); got != 4 {
		t.Errorf("after 1000 push/pop cycles with <= 4 values, len(buf) = %d, want 4", got)
	}
	if got, want := q.Values(), []int{997, 998, 999}; !slices.Equal(got, want) {
		t.Errorf("Values() = %v, want %v", got, want)
	}
}

func TestPopClearsSlot(t *testing.T) {
	var q Queue[*int]
	a, b := new(int), new(int)
	q.Push(a)
	q.Push(b)
	if v, ok := q.Pop(); v != a || !ok {
		t.Fatalf("Pop() = (%p, %v), want (%p, true)", v, ok, a)
	}
	for i, p := range q.buf {
		if p == a {
			t.Errorf("after Pop, buf[%d] still holds the popped pointer, want nil", i)
		}
	}
}

func TestValuesIsCopy(t *testing.T) {
	var q Queue[int]
	q.Push(1)
	q.Push(2)
	vals := q.Values()
	if len(vals) != 2 {
		t.Fatalf("Values() = %v, want [1 2]", vals)
	}
	vals[0] = 99
	if got := q.Values(); !slices.Equal(got, []int{1, 2}) {
		t.Errorf("modifying the Values() result changed the queue: Values() = %v, want [1 2]", got)
	}
}

func TestBounded(t *testing.T) {
	tests := []struct {
		limit   int
		push    []string
		want    []string
		dropped int
	}{
		{2, []string{"a", "b", "c"}, []string{"b", "c"}, 1},
		{3, []string{"a", "b"}, []string{"a", "b"}, 0},
		{3, []string{"a", "b", "c", "d", "e", "f", "g"}, []string{"e", "f", "g"}, 4},
		{1, []string{"x", "y"}, []string{"y"}, 1},
		{0, []string{"x", "y"}, []string{"y"}, 1},
		{-5, []string{"x", "y", "z"}, []string{"z"}, 2},
	}
	for _, tt := range tests {
		b := NewBounded[string](tt.limit)
		for _, v := range tt.push {
			b.Push(v)
		}
		if got := b.Values(); !slices.Equal(got, tt.want) {
			t.Errorf("NewBounded(%d), Push %q: Values() = %q, want %q", tt.limit, tt.push, got, tt.want)
		}
		if got := b.Dropped(); got != tt.dropped {
			t.Errorf("NewBounded(%d), Push %q: Dropped() = %d, want %d", tt.limit, tt.push, got, tt.dropped)
		}
	}
}

func TestBoundedPromotedMethods(t *testing.T) {
	b := NewBounded[int](2)
	b.Push(1)
	b.Push(2)
	b.Push(3)
	if v, ok := b.Peek(); v != 2 || !ok {
		t.Errorf("Bounded Peek() = (%d, %v), want (2, true)", v, ok)
	}
	if v, ok := b.Pop(); v != 2 || !ok {
		t.Errorf("Bounded Pop() = (%d, %v), want (2, true)", v, ok)
	}
	b.Push(4) // room again: nothing dropped
	if got, want := b.Values(), []int{3, 4}; !slices.Equal(got, want) {
		t.Errorf("Values() = %v, want %v", got, want)
	}
	if got := b.Dropped(); got != 1 {
		t.Errorf("Dropped() = %d, want 1", got)
	}
	if got := b.Len(); got != 2 {
		t.Errorf("Len() = %d, want 2", got)
	}
}
