package prefix

import (
	"slices"
	"testing"
)

func TestPrefixSums(t *testing.T) {
	tests := []struct {
		in   []int
		want []int
	}{
		{[]int{7}, []int{7}},
		{[]int{1, 2, 3, 4}, []int{1, 3, 6, 10}},
		{[]int{5, -5, 5}, []int{5, 0, 5}},
		{[]int{0, 0, 0}, []int{0, 0, 0}},
		{[]int{-1, -2, -3}, []int{-1, -3, -6}},
		{[]int{10, 20, 30, 40, 50}, []int{10, 30, 60, 100, 150}},
	}
	for _, tt := range tests {
		in := slices.Clone(tt.in)
		got := PrefixSums(in)
		if !slices.Equal(got, tt.want) {
			t.Errorf("PrefixSums(%v) = %v, want %v", tt.in, got, tt.want)
		}
		if !slices.Equal(in, tt.in) {
			t.Errorf("PrefixSums modified its input: %v became %v", tt.in, in)
		}
	}
}

func TestPrefixSumsEmpty(t *testing.T) {
	for _, in := range [][]int{nil, {}} {
		if got := PrefixSums(in); len(got) != 0 {
			t.Errorf("PrefixSums(%v) = %v, want an empty slice", in, got)
		}
	}
}

func TestPrefixSumsNewSlice(t *testing.T) {
	in := []int{1, 2, 3}
	got := PrefixSums(in)
	if len(got) != 3 {
		t.Fatalf("PrefixSums(%v) = %v, want 3 elements", in, got)
	}
	got[0] = 99
	if in[0] != 1 {
		t.Errorf("PrefixSums result shares memory with its input: writing to the result changed the input to %v", in)
	}
}

func TestRangeSum(t *testing.T) {
	s := []int{3, 1, 4, 1, 5, 9, 2, 6}
	p := []int{3, 4, 8, 9, 14, 23, 25, 31} // PrefixSums(s), written out so this test stands alone
	for lo := 0; lo <= len(s); lo++ {
		for hi := lo; hi <= len(s); hi++ {
			want := 0
			for _, v := range s[lo:hi] {
				want += v
			}
			if got := RangeSum(p, lo, hi); got != want {
				t.Errorf("RangeSum(%v, %d, %d) = %d, want %d (sum of %v)", p, lo, hi, got, want, s[lo:hi])
			}
		}
	}
}

func TestRangeSumSmall(t *testing.T) {
	p := []int{1, 3, 6, 10}
	tests := []struct {
		lo, hi, want int
	}{
		{1, 3, 5},
		{0, 4, 10},
		{2, 2, 0},
		{0, 0, 0},
		{4, 4, 0},
		{3, 4, 4},
		{0, 1, 1},
	}
	for _, tt := range tests {
		if got := RangeSum(p, tt.lo, tt.hi); got != tt.want {
			t.Errorf("RangeSum(%v, %d, %d) = %d, want %d", p, tt.lo, tt.hi, got, tt.want)
		}
	}
}
