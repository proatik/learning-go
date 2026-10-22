package matrix

import (
	"errors"
	"fmt"
	"slices"
	"testing"
)

type Cents int64

func mustRows[T Number](t *testing.T, rows [][]T) *Matrix[T] {
	t.Helper()
	m, err := FromRows(rows)
	if err != nil || m == nil {
		t.Fatalf("FromRows(%v) = %v, %v, want a matrix, <nil>", rows, m, err)
	}
	return m
}

// rowsOf collects m.Rows() into a [][]T.
func rowsOf[T Number](m *Matrix[T]) [][]T {
	var out [][]T
	for _, row := range m.Rows() {
		out = append(out, row)
	}
	return out
}

func equal[T Number](a, b [][]T) bool {
	return slices.EqualFunc(a, b, func(x, y []T) bool { return slices.Equal(x, y) })
}

func checkRows[T Number](t *testing.T, op string, m *Matrix[T], want [][]T) {
	t.Helper()
	if m == nil {
		t.Fatalf("%s = nil, want %v", op, want)
	}
	if r, c := m.Dims(); r != len(want) || c != len(want[0]) {
		t.Fatalf("%s.Dims() = %d, %d, want %d, %d", op, r, c, len(want), len(want[0]))
	}
	if got := rowsOf(m); !equal(got, want) {
		t.Fatalf("%s rows = %v, want %v", op, got, want)
	}
}

func TestNew(t *testing.T) {
	m, err := New[float64](2, 3)
	if err != nil || m == nil {
		t.Fatalf("New[float64](2, 3) = %v, %v, want a matrix, <nil>", m, err)
	}
	checkRows(t, "New(2, 3)", m, [][]float64{{0, 0, 0}, {0, 0, 0}})

	for _, size := range [][2]int{{0, 3}, {2, 0}, {-1, 2}, {0, 0}} {
		m, err := New[int](size[0], size[1])
		if m != nil || !errors.Is(err, ErrShape) {
			t.Errorf("New(%d, %d) = %v, %v, want <nil>, %v", size[0], size[1], m, err, ErrShape)
		}
	}
}

func TestFromRows(t *testing.T) {
	in := [][]int{{1, 2, 3}, {4, 5, 6}}
	m := mustRows(t, in)
	checkRows(t, "FromRows", m, [][]int{{1, 2, 3}, {4, 5, 6}})
	in[0][0] = 99
	if got := m.At(0, 0); got != 1 {
		t.Errorf("At(0, 0) after changing the input = %d, want 1 (FromRows must copy)", got)
	}

	bad := [][][]int{nil, {}, {{}}, {{1, 2}, {3}}, {{1}, {2, 3}}}
	for _, rows := range bad {
		m, err := FromRows(rows)
		if m != nil || !errors.Is(err, ErrShape) {
			t.Errorf("FromRows(%v) = %v, %v, want <nil>, %v", rows, m, err, ErrShape)
		}
	}
}

func TestAtSetDims(t *testing.T) {
	m, _ := New[int](3, 2)
	if m == nil {
		t.Fatalf("New[int](3, 2) = nil, want a matrix")
	}
	m.Set(2, 1, 7)
	m.Set(0, 1, -4)
	if got := m.At(2, 1); got != 7 {
		t.Errorf("At(2, 1) = %d, want 7", got)
	}
	checkRows(t, "after Set", m, [][]int{{0, -4}, {0, 0}, {0, 7}})
}

func TestIdentity(t *testing.T) {
	checkRows(t, "Identity[int](3)", Identity[int](3), [][]int{{1, 0, 0}, {0, 1, 0}, {0, 0, 1}})
	checkRows(t, "Identity[float32](1)", Identity[float32](1), [][]float32{{1}})
	if got := Identity[int](0); got != nil {
		t.Errorf("Identity(0) = %v, want <nil>", got)
	}
}

func TestAdd(t *testing.T) {
	a := mustRows(t, [][]int{{1, 2}, {3, 4}})
	b := mustRows(t, [][]int{{10, 20}, {30, -4}})
	sum, err := a.Add(b)
	if err != nil {
		t.Fatalf("Add error = %v, want <nil>", err)
	}
	checkRows(t, "a.Add(b)", sum, [][]int{{11, 22}, {33, 0}})
	checkRows(t, "a after Add", a, [][]int{{1, 2}, {3, 4}})

	c := mustRows(t, [][]int{{1, 2, 3}, {4, 5, 6}})
	got, err := a.Add(c)
	if want := "matrix: bad shape: 2x2 + 2x3"; got != nil || !errors.Is(err, ErrShape) || err.Error() != want {
		t.Errorf("2x2.Add(2x3) = %v, %v, want <nil>, %q", got, err, want)
	}
}

func TestMul(t *testing.T) {
	a := mustRows(t, [][]int{{1, 2, 3}, {4, 5, 6}})
	b := mustRows(t, [][]int{{7, 8}, {9, 10}, {11, 12}})
	p, err := a.Mul(b)
	if err != nil {
		t.Fatalf("Mul error = %v, want <nil>", err)
	}
	checkRows(t, "2x3.Mul(3x2)", p, [][]int{{58, 64}, {139, 154}})

	q, err := b.Mul(a)
	if err != nil {
		t.Fatalf("Mul error = %v, want <nil>", err)
	}
	checkRows(t, "3x2.Mul(2x3)", q, [][]int{{39, 54, 69}, {49, 68, 87}, {59, 82, 105}})

	id, err := a.Mul(Identity[int](3))
	if err != nil {
		t.Fatalf("Mul(Identity) error = %v, want <nil>", err)
	}
	checkRows(t, "a.Mul(I)", id, [][]int{{1, 2, 3}, {4, 5, 6}})

	row := mustRows(t, [][]float64{{0.5, 2}})
	col := mustRows(t, [][]float64{{4}, {0.25}})
	dot, err := row.Mul(col)
	if err != nil {
		t.Fatalf("Mul error = %v, want <nil>", err)
	}
	checkRows(t, "1x2.Mul(2x1)", dot, [][]float64{{2.5}})

	got, err := a.Mul(a)
	if want := "matrix: bad shape: 2x3 * 2x3"; got != nil || !errors.Is(err, ErrShape) || err.Error() != want {
		t.Errorf("2x3.Mul(2x3) = %v, %v, want <nil>, %q", got, err, want)
	}
}

func TestTranspose(t *testing.T) {
	a := mustRows(t, [][]int{{1, 2, 3}, {4, 5, 6}})
	checkRows(t, "Transpose", a.Transpose(), [][]int{{1, 4}, {2, 5}, {3, 6}})
	checkRows(t, "a after Transpose", a, [][]int{{1, 2, 3}, {4, 5, 6}})
	checkRows(t, "Transpose of 1x1", mustRows(t, [][]int{{9}}).Transpose(), [][]int{{9}})
}

func TestSumAndNamedTypes(t *testing.T) {
	prices := mustRows(t, [][]Cents{{199, 250}, {1, -50}})
	if got := prices.Sum(); got != 400 {
		t.Errorf("Sum() = %d, want 400", got)
	}
	if got := mustRows(t, [][]float64{{0.25, 0.5}}).Sum(); got != 0.75 {
		t.Errorf("Sum() = %v, want 0.75", got)
	}
}

func TestApply(t *testing.T) {
	prices := mustRows(t, [][]Cents{{150, 275}, {5, 0}})
	dollars := prices.Apply(func(c Cents) float64 { return float64(c) / 100 })
	checkRows(t, "Apply(cents → dollars)", dollars, [][]float64{{1.5, 2.75}, {0.05, 0}})

	sq := mustRows(t, [][]int{{1, -2, 3}}).Apply(func(n int) int { return n * n })
	checkRows(t, "Apply(square)", sq, [][]int{{1, 4, 9}})
}

func TestRows(t *testing.T) {
	m := mustRows(t, [][]int{{1, 2}, {3, 4}, {5, 6}})
	var idx []int
	for i, row := range m.Rows() {
		idx = append(idx, i)
		row[0] = 100
	}
	if want := []int{0, 1, 2}; !slices.Equal(idx, want) {
		t.Errorf("Rows() indexes = %v, want %v", idx, want)
	}
	checkRows(t, "after changing yielded rows", m, [][]int{{1, 2}, {3, 4}, {5, 6}})

	n := 0
	for i := range m.Rows() {
		n++
		if i == 1 {
			break
		}
	}
	if n != 2 {
		t.Errorf("Rows() with break at row 1 ran %d times, want 2", n)
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		m    fmt.Stringer
		want string
	}{
		{mustRows(t, [][]int{{1, 2}, {3, 4}}), "[1 2]\n[3 4]"},
		{mustRows(t, [][]float64{{0.5, -1}}), "[0.5 -1]"},
		{mustRows(t, [][]int{{7}}), "[7]"},
	}
	for _, tt := range tests {
		if got := tt.m.String(); got != tt.want {
			t.Errorf("String() = %q, want %q", got, tt.want)
		}
	}
}
