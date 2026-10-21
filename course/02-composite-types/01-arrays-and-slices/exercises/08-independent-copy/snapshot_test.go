package snapshot

import (
	"reflect"
	"slices"
	"testing"
)

func TestClone(t *testing.T) {
	tests := [][]int{
		{1, 2, 3},
		{42},
		{0, 0},
		{-1, 5, -1, 5},
	}
	for _, in := range tests {
		got := Clone(in)
		if !slices.Equal(got, in) {
			t.Errorf("Clone(%v) = %v, want %v", in, got, in)
		}
	}
}

func TestCloneNilAndEmpty(t *testing.T) {
	if got := Clone(nil); got != nil {
		t.Errorf("Clone(nil) = %#v, want nil", got)
	}
	if got := Clone([]int{}); got == nil || len(got) != 0 {
		t.Errorf("Clone([]int{}) = %#v, want a non-nil empty slice", got)
	}
}

func TestCloneIndependent(t *testing.T) {
	s := []int{1, 2, 3}
	c := Clone(s)
	if len(c) != 3 {
		t.Fatalf("Clone(%v) = %v, want 3 elements", s, c)
	}
	c[0] = 100
	if s[0] != 1 {
		t.Errorf("writing to the clone changed the original: s = %v", s)
	}
	s[2] = 300
	if c[2] != 3 {
		t.Errorf("writing to the original changed the clone: clone = %v", c)
	}
}

func TestCloneCapacityIsIndependent(t *testing.T) {
	backing := []int{1, 2, 3, 4, 5}
	s := backing[:2] // len 2, cap 5: appending to s would write backing[2]
	c := Clone(s)
	if len(c) != 2 {
		t.Fatalf("Clone(%v) = %v, want 2 elements", s, c)
	}
	_ = append(c, 99)
	if backing[2] != 3 {
		t.Errorf("appending to the clone overwrote the original's array: %v", backing)
	}
}

func TestWithAppended(t *testing.T) {
	tests := []struct {
		in   []int
		v    int
		want []int
	}{
		{nil, 1, []int{1}},
		{[]int{}, 1, []int{1}},
		{[]int{1, 2}, 3, []int{1, 2, 3}},
		{[]int{5}, 5, []int{5, 5}},
	}
	for _, tt := range tests {
		in := slices.Clone(tt.in)
		got := WithAppended(in, tt.v)
		if !slices.Equal(got, tt.want) {
			t.Errorf("WithAppended(%v, %d) = %v, want %v", tt.in, tt.v, got, tt.want)
		}
		if !slices.Equal(in, tt.in) {
			t.Errorf("WithAppended modified its input: %v became %v", tt.in, in)
		}
	}
}

func TestWithAppendedSpareCapacity(t *testing.T) {
	base := make([]int, 2, 10)
	base[0], base[1] = 7, 8
	a := WithAppended(base, 1)
	b := WithAppended(base, 2)
	if !slices.Equal(a, []int{7, 8, 1}) || !slices.Equal(b, []int{7, 8, 2}) {
		t.Fatalf("WithAppended(base, 1), WithAppended(base, 2) = %v, %v, want [7 8 1], [7 8 2] (did the second call overwrite the first?)", a, b)
	}
	if spare := base[:3]; spare[2] != 0 {
		t.Errorf("WithAppended wrote into base's spare capacity: base[:3] = %v", spare)
	}
	a[0] = -1
	if base[0] != 7 || b[0] != 7 {
		t.Errorf("results of WithAppended share memory: base = %v, b = %v", base, b)
	}
}

func TestClone2D(t *testing.T) {
	tests := [][][]int{
		{{1, 2}, {3}},
		{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}},
		{{}},
		{{1}, nil, {2, 3}},
		{},
	}
	for _, in := range tests {
		got := Clone2D(in)
		if !reflect.DeepEqual(got, in) {
			t.Errorf("Clone2D(%v) = %v, want %v", in, got, in)
		}
	}
	if got := Clone2D(nil); got != nil {
		t.Errorf("Clone2D(nil) = %#v, want nil", got)
	}
}

func TestClone2DDeep(t *testing.T) {
	m := [][]int{{1, 2}, {3, 4}}
	c := Clone2D(m)
	if len(c) != 2 || len(c[0]) != 2 || len(c[1]) != 2 {
		t.Fatalf("Clone2D(%v) = %v, want [[1 2] [3 4]]", m, c)
	}
	c[0][0] = 100
	c[1] = []int{9, 9}
	if m[0][0] != 1 {
		t.Errorf("writing c[0][0] changed the original: m = %v (rows must be copied, not shared)", m)
	}
	if !slices.Equal(m[1], []int{3, 4}) {
		t.Errorf("replacing c[1] changed the original: m = %v", m)
	}
	m[1][1] = 400
	if c[1][1] != 9 {
		t.Errorf("writing the original changed the clone: c = %v", c)
	}
}
