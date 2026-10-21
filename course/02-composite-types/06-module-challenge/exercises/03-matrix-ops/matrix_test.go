package matrix

import (
	"slices"
	"testing"
)

func equal(a, b [][]int) bool {
	return slices.EqualFunc(a, b, func(x, y []int) bool { return slices.Equal(x, y) })
}

func clone(m [][]int) [][]int {
	out := make([][]int, len(m))
	for i, row := range m {
		out[i] = slices.Clone(row)
	}
	return out
}

func TestMultiply(t *testing.T) {
	tests := []struct {
		a, b, want [][]int
	}{
		{[][]int{{1, 2}, {3, 4}}, [][]int{{5, 6}, {7, 8}}, [][]int{{19, 22}, {43, 50}}},
		{[][]int{{2}}, [][]int{{-3}}, [][]int{{-6}}},
		{[][]int{{1, 2, 3}}, [][]int{{4}, {5}, {6}}, [][]int{{32}}},
		{[][]int{{4}, {5}, {6}}, [][]int{{1, 2, 3}}, [][]int{{4, 8, 12}, {5, 10, 15}, {6, 12, 18}}},
		{[][]int{{1, 0}, {0, 1}}, [][]int{{9, 8, 7}, {6, 5, 4}}, [][]int{{9, 8, 7}, {6, 5, 4}}},
		{[][]int{{1, 2, 3}, {4, 5, 6}}, [][]int{{7, 8}, {9, 10}, {11, 12}}, [][]int{{58, 64}, {139, 154}}},
		{[][]int{{0, 0}, {0, 0}}, [][]int{{1, 2}, {3, 4}}, [][]int{{0, 0}, {0, 0}}},
	}
	for _, tt := range tests {
		a0, b0 := clone(tt.a), clone(tt.b)
		got, err := Multiply(tt.a, tt.b)
		if err != nil || !equal(got, tt.want) {
			t.Errorf("Multiply(%v, %v) = %v, %v, want %v, nil", tt.a, tt.b, got, err, tt.want)
		}
		if !equal(tt.a, a0) || !equal(tt.b, b0) {
			t.Errorf("Multiply modified its input")
		}
	}
}

func TestMultiplyInvalid(t *testing.T) {
	tests := []struct {
		name string
		a, b [][]int
	}{
		{"nil a", nil, [][]int{{1}}},
		{"nil b", [][]int{{1}}, nil},
		{"zero-width a", [][]int{{}}, [][]int{{1}}},
		{"zero-width b", [][]int{{1}}, [][]int{{}}},
		{"mismatch", [][]int{{1, 2}}, [][]int{{1, 2}}},
		{"mismatch 2", [][]int{{1, 2, 3}, {4, 5, 6}}, [][]int{{1, 2}, {3, 4}}},
		{"ragged a", [][]int{{1, 2}, {3}}, [][]int{{1}, {2}}},
		{"ragged b", [][]int{{1, 2}}, [][]int{{1}, {2, 3}}},
	}
	for _, tt := range tests {
		if got, err := Multiply(tt.a, tt.b); err == nil || got != nil {
			t.Errorf("%s: Multiply(%v, %v) = %v, %v, want nil, error", tt.name, tt.a, tt.b, got, err)
		}
	}
}

func TestRotate(t *testing.T) {
	tests := []struct {
		m, want [][]int
	}{
		{[][]int{{1}}, [][]int{{1}}},
		{[][]int{{1, 2}, {3, 4}}, [][]int{{3, 1}, {4, 2}}},
		{[][]int{{1, 2, 3}, {4, 5, 6}}, [][]int{{4, 1}, {5, 2}, {6, 3}}},
		{[][]int{{1, 2, 3}}, [][]int{{1}, {2}, {3}}},
		{[][]int{{1}, {2}, {3}}, [][]int{{3, 2, 1}}},
		{[][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}, [][]int{{7, 4, 1}, {8, 5, 2}, {9, 6, 3}}},
	}
	for _, tt := range tests {
		m0 := clone(tt.m)
		got, err := Rotate(tt.m)
		if err != nil || !equal(got, tt.want) {
			t.Errorf("Rotate(%v) = %v, %v, want %v, nil", tt.m, got, err, tt.want)
			continue
		}
		if !equal(tt.m, m0) {
			t.Errorf("Rotate modified its input: %v, was %v", tt.m, m0)
		}
		got[0][0] = 999 // must not write through to the input
		if !equal(tt.m, m0) {
			t.Errorf("Rotate result shares memory with its input")
		}
	}
}

func TestRotateFourTimes(t *testing.T) {
	m := [][]int{{1, 2, 3, 4}, {5, 6, 7, 8}, {9, 10, 11, 12}}
	cur := m
	for range 4 {
		next, err := Rotate(cur)
		if err != nil {
			t.Fatalf("Rotate(%v) error = %v", cur, err)
		}
		cur = next
	}
	if !equal(cur, m) {
		t.Fatalf("four rotations = %v, want %v", cur, m)
	}
}

func TestRotateEmptyAndRagged(t *testing.T) {
	for _, m := range [][][]int{nil, {}, {{}}, {{}, {}}} {
		got, err := Rotate(m)
		if err != nil || got == nil || len(got) != 0 {
			t.Errorf("Rotate(%v) = %v, %v, want [][]int{}, nil", m, got, err)
		}
	}
	if got, err := Rotate([][]int{{1, 2}, {3}}); err == nil || got != nil {
		t.Errorf("Rotate(ragged) = %v, %v, want nil, error", got, err)
	}
}

func TestSpiral(t *testing.T) {
	tests := []struct {
		m    [][]int
		want []int
	}{
		{[][]int{{1}}, []int{1}},
		{[][]int{{1, 2}, {3, 4}}, []int{1, 2, 4, 3}},
		{[][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}, []int{1, 2, 3, 6, 9, 8, 7, 4, 5}},
		{[][]int{{1, 2, 3, 4}, {5, 6, 7, 8}, {9, 10, 11, 12}}, []int{1, 2, 3, 4, 8, 12, 11, 10, 9, 5, 6, 7}},
		{[][]int{{1, 2, 3, 4}}, []int{1, 2, 3, 4}},
		{[][]int{{1}, {2}, {3}}, []int{1, 2, 3}},
		{[][]int{{1, 2}, {3, 4}, {5, 6}}, []int{1, 2, 4, 6, 5, 3}},
		{[][]int{{1, 2, 3}, {4, 5, 6}}, []int{1, 2, 3, 6, 5, 4}},
		{
			[][]int{{1, 2, 3, 4}, {5, 6, 7, 8}, {9, 10, 11, 12}, {13, 14, 15, 16}},
			[]int{1, 2, 3, 4, 8, 12, 16, 15, 14, 13, 9, 5, 6, 7, 11, 10},
		},
	}
	for _, tt := range tests {
		m0 := clone(tt.m)
		got, err := Spiral(tt.m)
		if err != nil || !slices.Equal(got, tt.want) || len(got) != len(tt.want) {
			t.Errorf("Spiral(%v) = %v, %v, want %v, nil", tt.m, got, err, tt.want)
		}
		if !equal(tt.m, m0) {
			t.Errorf("Spiral modified its input")
		}
	}
}

func TestSpiralEmptyAndRagged(t *testing.T) {
	for _, m := range [][][]int{nil, {}, {{}}} {
		got, err := Spiral(m)
		if err != nil || got == nil || len(got) != 0 {
			t.Errorf("Spiral(%v) = %v, %v, want []int{}, nil", m, got, err)
		}
	}
	if got, err := Spiral([][]int{{1}, {2, 3}}); err == nil || got != nil {
		t.Errorf("Spiral(ragged) = %v, %v, want nil, error", got, err)
	}
}
