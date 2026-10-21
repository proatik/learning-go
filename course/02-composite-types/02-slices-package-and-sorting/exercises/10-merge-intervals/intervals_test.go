package intervals

import (
	"slices"
	"testing"
)

var tests = []struct {
	in          [][2]int
	wantMerge   [][2]int
	wantCovered int
	wantGaps    [][2]int
}{
	{[][2]int{{1, 3}, {8, 10}, {2, 6}}, [][2]int{{1, 6}, {8, 10}}, 7, [][2]int{{6, 8}}},
	{[][2]int{{1, 3}, {3, 5}}, [][2]int{{1, 5}}, 4, [][2]int{}},
	{[][2]int{{5, 7}, {1, 2}, {3, 4}}, [][2]int{{1, 2}, {3, 4}, {5, 7}}, 4, [][2]int{{2, 3}, {4, 5}}},
	{[][2]int{{1, 10}, {2, 3}, {4, 5}}, [][2]int{{1, 10}}, 9, [][2]int{}},
	{[][2]int{{4, 4}}, [][2]int{{4, 4}}, 0, [][2]int{}},
	{[][2]int{{2, 5}, {2, 5}, {2, 5}}, [][2]int{{2, 5}}, 3, [][2]int{}},
	{[][2]int{{1, 4}, {0, 4}}, [][2]int{{0, 4}}, 4, [][2]int{}},
	{[][2]int{{6, 8}, {1, 9}, {2, 4}, {4, 7}}, [][2]int{{1, 9}}, 8, [][2]int{}},
	{[][2]int{{-5, -1}, {-3, 2}, {10, 12}}, [][2]int{{-5, 2}, {10, 12}}, 9, [][2]int{{2, 10}}},
	{[][2]int{{1, 2}, {2, 3}, {3, 4}, {6, 6}}, [][2]int{{1, 4}, {6, 6}}, 3, [][2]int{{4, 6}}},
	{[][2]int{{0, 1}, {5, 6}, {2, 3}, {3, 4}}, [][2]int{{0, 1}, {2, 4}, {5, 6}}, 4, [][2]int{{1, 2}, {4, 5}}},
	{[][2]int{}, [][2]int{}, 0, [][2]int{}},
	{nil, [][2]int{}, 0, [][2]int{}},
}

func TestMerge(t *testing.T) {
	for _, tt := range tests {
		orig := slices.Clone(tt.in)
		got := Merge(tt.in)
		if !slices.Equal(got, tt.wantMerge) {
			t.Errorf("Merge(%v) = %v, want %v", orig, got, tt.wantMerge)
		}
		if !slices.Equal(tt.in, orig) {
			t.Errorf("Merge modified its input: %v became %v", orig, tt.in)
		}
	}
}

func TestMergeDoesNotAlias(t *testing.T) {
	in := [][2]int{{1, 2}, {5, 6}}
	got := Merge(in)
	if len(got) != 2 {
		t.Fatalf("Merge(%v) = %v, want 2 intervals", in, got)
	}
	got[0][1] = 99
	if in[0] != [2]int{1, 2} {
		t.Errorf("writing to Merge's result changed the input: %v", in)
	}
}

func TestCovered(t *testing.T) {
	for _, tt := range tests {
		orig := slices.Clone(tt.in)
		if got := Covered(tt.in); got != tt.wantCovered {
			t.Errorf("Covered(%v) = %d, want %d", orig, got, tt.wantCovered)
		}
		if !slices.Equal(tt.in, orig) {
			t.Errorf("Covered modified its input: %v became %v", orig, tt.in)
		}
	}
}

func TestGaps(t *testing.T) {
	for _, tt := range tests {
		orig := slices.Clone(tt.in)
		if got := Gaps(tt.in); !slices.Equal(got, tt.wantGaps) {
			t.Errorf("Gaps(%v) = %v, want %v", orig, got, tt.wantGaps)
		}
		if !slices.Equal(tt.in, orig) {
			t.Errorf("Gaps modified its input: %v became %v", orig, tt.in)
		}
	}
}
