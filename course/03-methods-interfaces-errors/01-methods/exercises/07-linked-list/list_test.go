package list

import (
	"slices"
	"testing"
)

// build makes a list without relying on FromSlice.
func build(vals []int) *List {
	var head *List
	for i := len(vals) - 1; i >= 0; i-- {
		head = &List{Val: vals[i], Next: head}
	}
	return head
}

// collect reads a list without relying on Values (stops after 100 nodes in
// case of a cycle).
func collect(l *List) []int {
	var out []int
	for ; l != nil && len(out) < 100; l = l.Next {
		out = append(out, l.Val)
	}
	return out
}

func TestNilReceiver(t *testing.T) {
	var l *List
	if got := l.Len(); got != 0 {
		t.Errorf("nil.Len() = %d, want 0", got)
	}
	if got := l.Sum(); got != 0 {
		t.Errorf("nil.Sum() = %d, want 0", got)
	}
	if l.Contains(0) {
		t.Errorf("nil.Contains(0) = true, want false")
	}
	if got := l.Values(); got != nil {
		t.Errorf("nil.Values() = %v, want nil", got)
	}
	if got := l.String(); got != "[]" {
		t.Errorf("nil.String() = %q, want %q", got, "[]")
	}
	if got := l.Reverse(); got != nil {
		t.Errorf("nil.Reverse() = %v, want nil", collect(got))
	}
	if got := l.RemoveAll(1); got != nil {
		t.Errorf("nil.RemoveAll(1) = %v, want nil", collect(got))
	}
}

func TestFromSlice(t *testing.T) {
	tests := [][]int{nil, {}, {1}, {1, 2, 3}, {-5, 0, 5, 5}}
	for _, vals := range tests {
		got := FromSlice(vals)
		if len(vals) == 0 {
			if got != nil {
				t.Errorf("FromSlice(%v) = %v, want nil", vals, collect(got))
			}
			continue
		}
		if c := collect(got); !slices.Equal(c, vals) {
			t.Errorf("FromSlice(%v) = %v, want %v", vals, c, vals)
		}
	}
}

func TestQueries(t *testing.T) {
	tests := []struct {
		vals []int
		len  int
		sum  int
		str  string
	}{
		{[]int{1}, 1, 1, "[1]"},
		{[]int{1, 2, 3}, 3, 6, "[1 -> 2 -> 3]"},
		{[]int{-4, 0, 4, 10}, 4, 10, "[-4 -> 0 -> 4 -> 10]"},
	}
	for _, tt := range tests {
		l := build(tt.vals)
		if got := l.Len(); got != tt.len {
			t.Errorf("%v.Len() = %d, want %d", tt.vals, got, tt.len)
		}
		if got := l.Sum(); got != tt.sum {
			t.Errorf("%v.Sum() = %d, want %d", tt.vals, got, tt.sum)
		}
		if got := l.String(); got != tt.str {
			t.Errorf("%v.String() = %q, want %q", tt.vals, got, tt.str)
		}
		if got := l.Values(); !slices.Equal(got, tt.vals) {
			t.Errorf("%v.Values() = %v, want %v", tt.vals, got, tt.vals)
		}
	}
}

func TestContains(t *testing.T) {
	l := build([]int{3, 1, 4})
	for _, tt := range []struct {
		v    int
		want bool
	}{{3, true}, {4, true}, {1, true}, {2, false}, {0, false}} {
		if got := l.Contains(tt.v); got != tt.want {
			t.Errorf("[3 1 4].Contains(%d) = %v, want %v", tt.v, got, tt.want)
		}
	}
}

func TestPush(t *testing.T) {
	var l *List
	l = l.Push(2)
	l = l.Push(1)
	l = l.Push(0)
	if got, want := collect(l), []int{0, 1, 2}; !slices.Equal(got, want) {
		t.Errorf("nil.Push(2).Push(1).Push(0) = %v, want %v", got, want)
	}
}

func TestPushReusesNodes(t *testing.T) {
	old := build([]int{1, 2})
	l := old.Push(0)
	if l == nil || l.Next != old {
		t.Errorf("Push(0) on [1 2]: new head's Next should be the old head")
	}
}

func TestAppend(t *testing.T) {
	tests := []struct {
		vals []int
		v    int
		want []int
	}{
		{nil, 5, []int{5}},
		{[]int{1}, 2, []int{1, 2}},
		{[]int{1, 2, 3}, 4, []int{1, 2, 3, 4}},
	}
	for _, tt := range tests {
		head := build(tt.vals)
		got := head.Append(tt.v)
		if c := collect(got); !slices.Equal(c, tt.want) {
			t.Errorf("%v.Append(%d) = %v, want %v", tt.vals, tt.v, c, tt.want)
		}
		if head != nil && got != head {
			t.Errorf("%v.Append(%d) returned a different head", tt.vals, tt.v)
		}
	}
}

func TestReverse(t *testing.T) {
	tests := []struct {
		vals, want []int
	}{
		{[]int{1}, []int{1}},
		{[]int{1, 2}, []int{2, 1}},
		{[]int{1, 2, 3, 4, 5}, []int{5, 4, 3, 2, 1}},
	}
	for _, tt := range tests {
		head := build(tt.vals)
		got := head.Reverse()
		if c := collect(got); !slices.Equal(c, tt.want) {
			t.Errorf("%v.Reverse() = %v, want %v", tt.vals, c, tt.want)
		}
		if got != nil && len(tt.vals) > 1 && head.Next != nil {
			t.Errorf("%v.Reverse(): old head should now be the last node", tt.vals)
		}
	}
}

func TestRemoveAll(t *testing.T) {
	tests := []struct {
		vals []int
		v    int
		want []int
	}{
		{[]int{2, 1, 2, 3, 2}, 2, []int{1, 3}},
		{[]int{1, 2, 3}, 4, []int{1, 2, 3}},
		{[]int{1, 2, 3}, 1, []int{2, 3}},
		{[]int{1, 2, 3}, 3, []int{1, 2}},
		{[]int{7, 7, 7}, 7, nil},
		{[]int{1, 7, 7, 2}, 7, []int{1, 2}},
	}
	for _, tt := range tests {
		got := build(tt.vals).RemoveAll(tt.v)
		if c := collect(got); !slices.Equal(c, tt.want) {
			t.Errorf("%v.RemoveAll(%d) = %v, want %v", tt.vals, tt.v, c, tt.want)
		}
	}
}
