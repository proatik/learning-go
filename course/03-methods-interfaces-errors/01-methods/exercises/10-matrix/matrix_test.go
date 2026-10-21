package matrix_test

import (
	"fmt"
	"slices"
	"testing"

	matrix "learngo/course/03-methods-interfaces-errors/01-methods/exercises/10-matrix"
)

// grid reads m back into a [][]float64 using only Rows, Cols and At.
func grid(m *matrix.Matrix) [][]float64 {
	out := make([][]float64, m.Rows())
	for i := range out {
		out[i] = make([]float64, m.Cols())
		for j := range out[i] {
			out[i][j], _ = m.At(i, j)
		}
	}
	return out
}

func sameGrid(a, b [][]float64) bool {
	return slices.EqualFunc(a, b, func(x, y []float64) bool { return slices.Equal(x, y) })
}

func mustRows(t *testing.T, rows [][]float64) *matrix.Matrix {
	t.Helper()
	m := matrix.FromRows(rows)
	if m == nil {
		t.Fatalf("FromRows(%v) = nil, want a matrix", rows)
	}
	return m
}

func TestNew(t *testing.T) {
	tests := []struct {
		rows, cols int
		ok         bool
	}{
		{2, 3, true},
		{1, 1, true},
		{0, 3, false},
		{3, 0, false},
		{-1, 2, false},
	}
	for _, tt := range tests {
		m := matrix.New(tt.rows, tt.cols)
		if !tt.ok {
			if m != nil {
				t.Errorf("New(%d, %d) = %v, want nil", tt.rows, tt.cols, grid(m))
			}
			continue
		}
		if m == nil {
			t.Errorf("New(%d, %d) = nil, want a %dx%d matrix", tt.rows, tt.cols, tt.rows, tt.cols)
			continue
		}
		if m.Rows() != tt.rows || m.Cols() != tt.cols {
			t.Errorf("New(%d, %d) shape = %dx%d, want %dx%d", tt.rows, tt.cols, m.Rows(), m.Cols(), tt.rows, tt.cols)
		}
		for i := range tt.rows {
			for j := range tt.cols {
				if v, ok := m.At(i, j); v != 0 || !ok {
					t.Errorf("New(%d, %d).At(%d, %d) = (%v, %v), want (0, true)", tt.rows, tt.cols, i, j, v, ok)
				}
			}
		}
	}
}

func TestFromRows(t *testing.T) {
	valid := [][][]float64{
		{{1, 2}, {3, 4}},
		{{1, 2, 3}},
		{{1}, {2}, {3}},
		{{-1.5, 0}, {2.25, 1e9}},
	}
	for _, rows := range valid {
		m := mustRows(t, rows)
		if got := grid(m); !sameGrid(got, rows) {
			t.Errorf("FromRows(%v) holds %v, want %v", rows, got, rows)
		}
	}
	invalid := [][][]float64{nil, {}, {{}}, {{1, 2}, {3}}, {{1}, {}}, {{1}, {2, 3}}}
	for _, rows := range invalid {
		if m := matrix.FromRows(rows); m != nil {
			t.Errorf("FromRows(%v) = %v, want nil", rows, grid(m))
		}
	}
}

func TestFromRowsCopies(t *testing.T) {
	rows := [][]float64{{1, 2}, {3, 4}}
	m := mustRows(t, rows)
	rows[0][0] = 99
	if v, _ := m.At(0, 0); v != 1 {
		t.Errorf("changing the input after FromRows changed the matrix: At(0, 0) = %v, want 1", v)
	}
}

func TestIdentity(t *testing.T) {
	tests := []struct {
		n    int
		want [][]float64
	}{
		{1, [][]float64{{1}}},
		{2, [][]float64{{1, 0}, {0, 1}}},
		{3, [][]float64{{1, 0, 0}, {0, 1, 0}, {0, 0, 1}}},
	}
	for _, tt := range tests {
		m := matrix.Identity(tt.n)
		if m == nil {
			t.Errorf("Identity(%d) = nil, want %v", tt.n, tt.want)
			continue
		}
		if got := grid(m); !sameGrid(got, tt.want) {
			t.Errorf("Identity(%d) = %v, want %v", tt.n, got, tt.want)
		}
	}
	if m := matrix.Identity(0); m != nil {
		t.Errorf("Identity(0) = %v, want nil", grid(m))
	}
}

func TestNilShape(t *testing.T) {
	var m *matrix.Matrix
	if m.Rows() != 0 || m.Cols() != 0 {
		t.Errorf("nil matrix Rows(), Cols() = %d, %d, want 0, 0", m.Rows(), m.Cols())
	}
}

func TestAtSet(t *testing.T) {
	m := mustRows(t, [][]float64{{1, 2, 3}, {4, 5, 6}})
	if v, ok := m.At(1, 2); v != 6 || !ok {
		t.Errorf("At(1, 2) = (%v, %v), want (6, true)", v, ok)
	}
	for _, ij := range [][2]int{{2, 0}, {0, 3}, {-1, 0}, {0, -1}} {
		if v, ok := m.At(ij[0], ij[1]); v != 0 || ok {
			t.Errorf("At(%d, %d) = (%v, %v), want (0, false)", ij[0], ij[1], v, ok)
		}
		if m.Set(ij[0], ij[1], 7) {
			t.Errorf("Set(%d, %d, 7) = true, want false", ij[0], ij[1])
		}
	}
	if !m.Set(0, 1, -8) {
		t.Errorf("Set(0, 1, -8) = false, want true")
	}
	want := [][]float64{{1, -8, 3}, {4, 5, 6}}
	if got := grid(m); !sameGrid(got, want) {
		t.Errorf("after Set(0, 1, -8) matrix = %v, want %v", got, want)
	}
}

func TestEqual(t *testing.T) {
	a := mustRows(t, [][]float64{{1, 2}, {3, 4}})
	b := mustRows(t, [][]float64{{1, 2}, {3, 4}})
	c := mustRows(t, [][]float64{{1, 2}, {3, 5}})
	flat := mustRows(t, [][]float64{{1, 2, 3, 4}})
	var n1, n2 *matrix.Matrix
	tests := []struct {
		desc string
		x, y *matrix.Matrix
		want bool
	}{
		{"same values", a, b, true},
		{"itself", a, a, true},
		{"different value", a, c, false},
		{"different shape", a, flat, false},
		{"nil vs matrix", n1, a, false},
		{"matrix vs nil", a, n1, false},
		{"nil vs nil", n1, n2, true},
	}
	for _, tt := range tests {
		if got := tt.x.Equal(tt.y); got != tt.want {
			t.Errorf("%s: Equal() = %v, want %v", tt.desc, got, tt.want)
		}
	}
}

func TestScale(t *testing.T) {
	m := mustRows(t, [][]float64{{1, -2}, {0, 4}})
	m.Scale(2.5)
	want := [][]float64{{2.5, -5}, {0, 10}}
	if got := grid(m); !sameGrid(got, want) {
		t.Errorf("Scale(2.5) = %v, want %v", got, want)
	}
}

func TestTranspose(t *testing.T) {
	tests := []struct {
		in, want [][]float64
	}{
		{[][]float64{{1, 2}, {3, 4}}, [][]float64{{1, 3}, {2, 4}}},
		{[][]float64{{1, 2, 3}}, [][]float64{{1}, {2}, {3}}},
		{[][]float64{{1, 2, 3}, {4, 5, 6}}, [][]float64{{1, 4}, {2, 5}, {3, 6}}},
	}
	for _, tt := range tests {
		m := mustRows(t, tt.in)
		tr := m.Transpose()
		if tr == nil {
			t.Errorf("FromRows(%v).Transpose() = nil, want %v", tt.in, tt.want)
			continue
		}
		if got := grid(tr); !sameGrid(got, tt.want) {
			t.Errorf("FromRows(%v).Transpose() = %v, want %v", tt.in, got, tt.want)
		}
		if got := grid(m); !sameGrid(got, tt.in) {
			t.Errorf("Transpose changed its receiver to %v, want %v", got, tt.in)
		}
	}
}

func TestAdd(t *testing.T) {
	a := mustRows(t, [][]float64{{1, 2}, {3, 4}})
	b := mustRows(t, [][]float64{{10, 20}, {30, 40}})
	sum := a.Add(b)
	if sum == nil {
		t.Fatalf("Add of two 2x2 matrices = nil, want a matrix")
	}
	want := [][]float64{{11, 22}, {33, 44}}
	if got := grid(sum); !sameGrid(got, want) {
		t.Errorf("%v.Add(%v) = %v, want %v", grid(a), grid(b), got, want)
	}
	if got := grid(a); !sameGrid(got, [][]float64{{1, 2}, {3, 4}}) {
		t.Errorf("Add changed its receiver to %v", got)
	}
	if got := a.Add(mustRows(t, [][]float64{{1, 2}})); got != nil {
		t.Errorf("Add with a different shape = %v, want nil", grid(got))
	}
	if got := a.Add(nil); got != nil {
		t.Errorf("Add(nil) = %v, want nil", grid(got))
	}
}

func TestMul(t *testing.T) {
	tests := []struct {
		a, b, want [][]float64
	}{
		{[][]float64{{1, 2}, {3, 4}}, [][]float64{{5, 6}, {7, 8}}, [][]float64{{19, 22}, {43, 50}}},
		{[][]float64{{1, 2, 3}}, [][]float64{{4}, {5}, {6}}, [][]float64{{32}}},
		{[][]float64{{4}, {5}, {6}}, [][]float64{{1, 2, 3}}, [][]float64{{4, 8, 12}, {5, 10, 15}, {6, 12, 18}}},
		{[][]float64{{1, 2, 3}, {4, 5, 6}}, [][]float64{{1, 0}, {0, 1}, {1, 1}}, [][]float64{{4, 5}, {10, 11}}},
	}
	for _, tt := range tests {
		p := mustRows(t, tt.a).Mul(mustRows(t, tt.b))
		if p == nil {
			t.Errorf("%v.Mul(%v) = nil, want %v", tt.a, tt.b, tt.want)
			continue
		}
		if got := grid(p); !sameGrid(got, tt.want) {
			t.Errorf("%v.Mul(%v) = %v, want %v", tt.a, tt.b, got, tt.want)
		}
	}

	a := mustRows(t, [][]float64{{1, 2}, {3, 4}})
	id := matrix.Identity(2)
	if id == nil {
		t.Fatalf("Identity(2) = nil, want a matrix")
	}
	if p := a.Mul(id); p == nil || !sameGrid(grid(p), grid(a)) {
		t.Errorf("A.Mul(Identity(2)) should equal A = %v", grid(a))
	}
	if p := a.Mul(matrix.New(3, 1)); p != nil {
		t.Errorf("2x2 Mul 3x1 = %v, want nil", grid(p))
	}
	if p := a.Mul(nil); p != nil {
		t.Errorf("Mul(nil) = %v, want nil", grid(p))
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		rows [][]float64
		want string
	}{
		{[][]float64{{1, 0}, {0, 1}}, "[1 0]\n[0 1]"},
		{[][]float64{{1.5, -2}}, "[1.5 -2]"},
		{[][]float64{{7}}, "[7]"},
		{[][]float64{{0.1}, {1e21}, {-0.25}}, "[0.1]\n[1e+21]\n[-0.25]"},
	}
	for _, tt := range tests {
		if got := mustRows(t, tt.rows).String(); got != tt.want {
			t.Errorf("FromRows(%v).String() = %q, want %q", tt.rows, got, tt.want)
		}
	}
}

func ExampleMatrix_String() {
	a := matrix.FromRows([][]float64{{1, 2}, {3, 4}})
	fmt.Println(a)
	fmt.Println("--")
	fmt.Println(a.Transpose())
	fmt.Println("--")
	fmt.Println(a.Mul(matrix.Identity(2)).Equal(a))
	// Output:
	// [1 2]
	// [3 4]
	// --
	// [1 3]
	// [2 4]
	// --
	// true
}
