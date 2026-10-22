package matrix

import (
	"errors"
	"fmt"
	"slices"
	"testing"
)

// mk builds a matrix without relying on FromRows.
func mk[T Number](rows, cols int, data ...T) *Matrix[T] {
	return &Matrix[T]{rows: rows, cols: cols, data: data}
}

// show formats a matrix from its fields, without relying on Rows.
func show[T Number](m *Matrix[T]) string {
	if m == nil {
		return "<nil>"
	}
	out := make([][]T, m.rows)
	for i := range m.rows {
		if (i+1)*m.cols <= len(m.data) {
			out[i] = m.data[i*m.cols : (i+1)*m.cols]
		}
	}
	return fmt.Sprintf("%dx%d %v", m.rows, m.cols, out)
}

// same reports whether a and b have the same shape and elements.
func same[T Number](a, b *Matrix[T]) bool {
	return a != nil && b != nil && a.rows == b.rows && a.cols == b.cols && slices.Equal(a.data, b.data)
}

func TestFromRows(t *testing.T) {
	in := [][]int{{1, 2, 3}, {4, 5, 6}}
	m, err := FromRows(in)
	if err != nil {
		t.Fatalf("FromRows(%v) error = %v, want nil", in, err)
	}
	if want := mk(2, 3, 1, 2, 3, 4, 5, 6); !same(m, want) {
		t.Fatalf("FromRows(%v) = %s, want %s", in, show(m), show(want))
	}
	in[0][0] = 99
	if m.data[0] != 1 {
		t.Errorf("FromRows must copy its input: changing the input changed the matrix to %s", show(m))
	}
	if r, c := m.Dims(); r != 2 || c != 3 {
		t.Errorf("Dims() = (%d, %d), want (2, 3)", r, c)
	}
	if got := m.At(1, 0); got != 4 {
		t.Errorf("At(1, 0) = %d, want 4", got)
	}
	if got := m.At(0, 2); got != 3 {
		t.Errorf("At(0, 2) = %d, want 3", got)
	}
}

func TestFromRowsErrors(t *testing.T) {
	tests := [][][]float64{
		nil,
		{},
		{{}},
		{{1, 2}, {3}},
		{{1}, {2, 3}},
		{{1, 2}, {}},
	}
	for _, in := range tests {
		m, err := FromRows(in)
		if !errors.Is(err, ErrShape) {
			t.Errorf("FromRows(%v) error = %v, want %v", in, err, ErrShape)
		}
		if m != nil {
			t.Errorf("FromRows(%v) = %s, want nil", in, show(m))
		}
	}
}

func TestIdentity(t *testing.T) {
	tests := []struct {
		n    int
		want *Matrix[float64]
	}{
		{1, mk[float64](1, 1, 1)},
		{2, mk[float64](2, 2, 1, 0, 0, 1)},
		{3, mk[float64](3, 3, 1, 0, 0, 0, 1, 0, 0, 0, 1)},
		{0, mk[float64](0, 0)},
		{-2, mk[float64](0, 0)},
	}
	for _, tt := range tests {
		got := Identity[float64](tt.n)
		if !same(got, tt.want) {
			t.Errorf("Identity(%d) = %s, want %s", tt.n, show(got), show(tt.want))
		}
	}
}

func TestRowsIsCopy(t *testing.T) {
	m := mk(2, 2, 1, 2, 3, 4)
	rows := m.Rows()
	if len(rows) != 2 || !slices.Equal(rows[0], []int{1, 2}) || !slices.Equal(rows[1], []int{3, 4}) {
		t.Fatalf("Rows() of [[1 2] [3 4]] = %v, want [[1 2] [3 4]]", rows)
	}
	rows[0][0] = 99
	rows[1] = append(rows[1], 5)
	if want := mk(2, 2, 1, 2, 3, 4); !same(m, want) {
		t.Errorf("modifying the Rows() result changed the matrix: %s, want %s", show(m), show(want))
	}
}

func TestAdd(t *testing.T) {
	a := mk(2, 2, 1, 2, 3, 4)
	b := mk(2, 2, 10, -20, 0, 1)
	got, err := a.Add(b)
	if want := mk(2, 2, 11, -18, 3, 5); err != nil || !same(got, want) {
		t.Errorf("%s.Add(%s) = (%s, %v), want (%s, nil)", show(a), show(b), show(got), err, show(want))
	}
	if !same(a, mk(2, 2, 1, 2, 3, 4)) || !same(b, mk(2, 2, 10, -20, 0, 1)) {
		t.Errorf("Add modified its operands: a = %s, b = %s", show(a), show(b))
	}
}

func TestMul(t *testing.T) {
	tests := []struct {
		a, b, want *Matrix[int]
	}{
		{mk(2, 2, 1, 2, 3, 4), mk(2, 2, 1, 0, 0, 1), mk(2, 2, 1, 2, 3, 4)},
		{mk(2, 2, 1, 2, 3, 4), mk(2, 2, 5, 6, 7, 8), mk(2, 2, 19, 22, 43, 50)},
		{mk(2, 3, 1, 2, 3, 4, 5, 6), mk(3, 2, 7, 8, 9, 10, 11, 12), mk(2, 2, 58, 64, 139, 154)},
		{mk(1, 3, 1, 2, 3), mk(3, 1, 4, 5, 6), mk(1, 1, 32)},
		{mk(3, 1, 1, 2, 3), mk(1, 2, -1, 2), mk(3, 2, -1, 2, -2, 4, -3, 6)},
	}
	for _, tt := range tests {
		got, err := tt.a.Mul(tt.b)
		if err != nil || !same(got, tt.want) {
			t.Errorf("%s.Mul(%s) = (%s, %v), want (%s, nil)", show(tt.a), show(tt.b), show(got), err, show(tt.want))
		}
	}
}

func TestShapeErrors(t *testing.T) {
	a23 := mk(2, 3, 1, 2, 3, 4, 5, 6)
	a32 := mk(3, 2, 1, 2, 3, 4, 5, 6)
	a22 := mk(2, 2, 1, 2, 3, 4)
	tests := []struct {
		name    string
		do      func() (*Matrix[int], error)
		wantMsg string
	}{
		{"2x3.Add(3x2)", func() (*Matrix[int], error) { return a23.Add(a32) }, "add 2x3 and 3x2: matrix: shape mismatch"},
		{"2x2.Add(2x3)", func() (*Matrix[int], error) { return a22.Add(a23) }, "add 2x2 and 2x3: matrix: shape mismatch"},
		{"2x3.Mul(2x3)", func() (*Matrix[int], error) { return a23.Mul(a23) }, "mul 2x3 and 2x3: matrix: shape mismatch"},
		{"3x2.Mul(3x2)", func() (*Matrix[int], error) { return a32.Mul(a32) }, "mul 3x2 and 3x2: matrix: shape mismatch"},
	}
	for _, tt := range tests {
		got, err := tt.do()
		if !errors.Is(err, ErrShape) {
			t.Errorf("%s error = %v, want %v", tt.name, err, ErrShape)
			continue
		}
		if err.Error() != tt.wantMsg {
			t.Errorf("%s error = %q, want %q", tt.name, err.Error(), tt.wantMsg)
		}
		if got != nil {
			t.Errorf("%s = %s, want nil on error", tt.name, show(got))
		}
	}
}

func TestScale(t *testing.T) {
	m := mk(2, 2, 1, 2, 3, 4)
	got := m.Scale(10)
	if want := mk(2, 2, 10, 20, 30, 40); !same(got, want) {
		t.Errorf("%s.Scale(10) = %s, want %s", show(m), show(got), show(want))
	}
	if got := m.Scale(0); !same(got, mk(2, 2, 0, 0, 0, 0)) {
		t.Errorf("%s.Scale(0) = %s, want 2x2 [[0 0] [0 0]]", show(m), show(got))
	}
	if !same(m, mk(2, 2, 1, 2, 3, 4)) {
		t.Errorf("Scale modified its receiver: %s", show(m))
	}
}

func TestTranspose(t *testing.T) {
	tests := []struct {
		m, want *Matrix[float64]
	}{
		{mk(1, 3, 1.0, 2, 3), mk(3, 1, 1.0, 2, 3)},
		{mk(2, 3, 1.0, 2, 3, 4, 5, 6), mk(3, 2, 1.0, 4, 2, 5, 3, 6)},
		{mk(2, 2, 1.5, -2, 0, 4), mk(2, 2, 1.5, 0, -2, 4)},
	}
	for _, tt := range tests {
		if got := tt.m.Transpose(); !same(got, tt.want) {
			t.Errorf("%s.Transpose() = %s, want %s", show(tt.m), show(got), show(tt.want))
		}
	}
}

type cents int64

func TestNamedNumberType(t *testing.T) {
	prices := mk[cents](1, 2, 199, 250)
	qty := mk[cents](2, 1, 3, 2)
	got, err := prices.Mul(qty)
	if want := mk[cents](1, 1, 1097); err != nil || !same(got, want) {
		t.Errorf("%s.Mul(%s) = (%s, %v), want (%s, nil)", show(prices), show(qty), show(got), err, show(want))
	}
}
