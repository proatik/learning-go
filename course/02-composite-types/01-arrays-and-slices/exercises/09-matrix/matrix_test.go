package matrix

import (
	"reflect"
	"slices"
	"testing"
)

func clone2D(m [][]int) [][]int {
	if m == nil {
		return nil
	}
	out := make([][]int, len(m))
	for i, row := range m {
		out[i] = slices.Clone(row)
	}
	return out
}

func TestTranspose(t *testing.T) {
	tests := []struct {
		in, want [][]int
	}{
		{[][]int{{1, 2, 3}, {4, 5, 6}}, [][]int{{1, 4}, {2, 5}, {3, 6}}},
		{[][]int{{1, 2, 3}}, [][]int{{1}, {2}, {3}}},
		{[][]int{{1}, {2}, {3}}, [][]int{{1, 2, 3}}},
		{[][]int{{7}}, [][]int{{7}}},
		{[][]int{{1, 2}, {3, 4}}, [][]int{{1, 3}, {2, 4}}},
		{
			[][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}},
			[][]int{{1, 4, 7}, {2, 5, 8}, {3, 6, 9}},
		},
	}
	for _, tt := range tests {
		in := clone2D(tt.in)
		got := Transpose(in)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("Transpose(%v) = %v, want %v", tt.in, got, tt.want)
		}
		if !reflect.DeepEqual(in, tt.in) {
			t.Errorf("Transpose modified its input: %v became %v", tt.in, in)
		}
	}
}

func TestTransposeEmpty(t *testing.T) {
	for _, in := range [][][]int{nil, {}, {{}}, {{}, {}}} {
		if got := Transpose(in); len(got) != 0 {
			t.Errorf("Transpose(%v) = %v, want an empty matrix", in, got)
		}
	}
}

func TestTransposeIsNew(t *testing.T) {
	m := [][]int{{1, 2}, {3, 4}}
	got := Transpose(m)
	if len(got) != 2 || len(got[0]) != 2 {
		t.Fatalf("Transpose(%v) = %v, want [[1 3] [2 4]]", m, got)
	}
	got[0][0] = 100
	if m[0][0] != 1 {
		t.Errorf("Transpose result shares memory with its input: m = %v", m)
	}
}

func TestTransposeTwiceIsIdentity(t *testing.T) {
	m := [][]int{{1, 2, 3, 4}, {5, 6, 7, 8}, {9, 10, 11, 12}}
	if got := Transpose(Transpose(m)); !reflect.DeepEqual(got, m) {
		t.Errorf("Transpose(Transpose(%v)) = %v, want the original", m, got)
	}
}

func TestSpiral(t *testing.T) {
	tests := []struct {
		in   [][]int
		want []int
	}{
		{[][]int{{1}}, []int{1}},
		{[][]int{{1, 2}, {3, 4}}, []int{1, 2, 4, 3}},
		{[][]int{{1, 2, 3}}, []int{1, 2, 3}},
		{[][]int{{1}, {2}, {3}}, []int{1, 2, 3}},
		{[][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}, []int{1, 2, 3, 6, 9, 8, 7, 4, 5}},
		{
			[][]int{{1, 2, 3, 4}, {5, 6, 7, 8}, {9, 10, 11, 12}},
			[]int{1, 2, 3, 4, 8, 12, 11, 10, 9, 5, 6, 7},
		},
		{
			[][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}, {10, 11, 12}},
			[]int{1, 2, 3, 6, 9, 12, 11, 10, 7, 4, 5, 8},
		},
		{
			[][]int{{1, 2}, {3, 4}, {5, 6}},
			[]int{1, 2, 4, 6, 5, 3},
		},
		{
			[][]int{{1, 2, 3, 4}, {5, 6, 7, 8}, {9, 10, 11, 12}, {13, 14, 15, 16}},
			[]int{1, 2, 3, 4, 8, 12, 16, 15, 14, 13, 9, 5, 6, 7, 11, 10},
		},
	}
	for _, tt := range tests {
		in := clone2D(tt.in)
		got := Spiral(in)
		if !slices.Equal(got, tt.want) {
			t.Errorf("Spiral(%v) = %v, want %v", tt.in, got, tt.want)
		}
		if !reflect.DeepEqual(in, tt.in) {
			t.Errorf("Spiral modified its input: %v became %v", tt.in, in)
		}
	}
}

func TestSpiralEmpty(t *testing.T) {
	for _, in := range [][][]int{nil, {}, {{}}, {{}, {}}} {
		if got := Spiral(in); len(got) != 0 {
			t.Errorf("Spiral(%v) = %v, want an empty slice", in, got)
		}
	}
}

func TestSpiralVisitsEveryElementOnce(t *testing.T) {
	for rows := 1; rows <= 6; rows++ {
		for cols := 1; cols <= 6; cols++ {
			m := make([][]int, rows)
			for i := range m {
				m[i] = make([]int, cols)
				for j := range m[i] {
					m[i][j] = i*cols + j
				}
			}
			got := Spiral(m)
			if len(got) != rows*cols {
				t.Errorf("Spiral of a %dx%d matrix returned %d elements, want %d", rows, cols, len(got), rows*cols)
				continue
			}
			seen := make([]bool, rows*cols)
			for _, v := range got {
				if v < 0 || v >= len(seen) || seen[v] {
					t.Errorf("Spiral of a %dx%d matrix = %v: element %d is invalid or repeated", rows, cols, got, v)
					break
				}
				seen[v] = true
			}
		}
	}
}
