package double

import (
	"slices"
	"testing"
)

func TestDoubleSlice(t *testing.T) {
	tests := []struct {
		in, want []int
	}{
		{[]int{1, 2, 3}, []int{2, 4, 6}},
		{[]int{0, -1, 5}, []int{0, -2, 10}},
		{[]int{7}, []int{14}},
		{[]int{}, []int{}},
	}
	for _, tt := range tests {
		s := slices.Clone(tt.in)
		DoubleSlice(s)
		if !slices.Equal(s, tt.want) {
			t.Errorf("DoubleSlice(%v): slice is now %v, want %v", tt.in, s, tt.want)
		}
	}
}

func TestDoubleSliceNil(t *testing.T) {
	DoubleSlice(nil) // must not panic
}

func TestDoubleSubSlice(t *testing.T) {
	s := []int{1, 2, 3}
	DoubleSlice(s[1:])
	want := []int{1, 4, 6}
	if !slices.Equal(s, want) {
		t.Errorf("DoubleSlice(s[1:]) on [1 2 3]: s is now %v, want %v", s, want)
	}
}

func TestDoubleArray(t *testing.T) {
	q := Quad{1, 2, 3, 4}
	DoubleArray(&q)
	if want := (Quad{2, 4, 6, 8}); q != want {
		t.Errorf("DoubleArray(&[1 2 3 4]): array is now %v, want %v", q, want)
	}
	DoubleArray(&q)
	if want := (Quad{4, 8, 12, 16}); q != want {
		t.Errorf("DoubleArray twice: array is now %v, want %v", q, want)
	}
}

func TestDoubleArrayNil(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("DoubleArray(nil) panicked: %v", r)
		}
	}()
	DoubleArray(nil)
}

func TestDoubled(t *testing.T) {
	q := Quad{1, 2, 3, 4}
	got := Doubled(q)
	if want := (Quad{2, 4, 6, 8}); got != want {
		t.Errorf("Doubled(%v) = %v, want %v", Quad{1, 2, 3, 4}, got, want)
	}
	if want := (Quad{1, 2, 3, 4}); q != want {
		t.Errorf("Doubled changed its argument: q is now %v, want %v", q, want)
	}
	if got := Doubled(Quad{-1, 0, 0, 5}); got != (Quad{-2, 0, 0, 10}) {
		t.Errorf("Doubled([-1 0 0 5]) = %v, want [-2 0 0 10]", got)
	}
}
