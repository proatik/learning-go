// Package matrix provides a generic dense matrix of numbers.
package matrix

import (
	"errors"
	"fmt"
	"iter"
	"strings"
)

// Number is the set of element types a Matrix can hold.
type Number interface {
	~int | ~int32 | ~int64 | ~float32 | ~float64
}

// ErrShape is wrapped by every error about mismatched or invalid sizes.
var ErrShape = errors.New("matrix: bad shape")

// Matrix is a rows × cols matrix stored row by row in data.
type Matrix[T Number] struct {
	rows, cols int
	data       []T // element (r, c) is data[r*cols+c]
}

// New returns a zero rows × cols matrix. Both sizes must be at least 1.
func New[T Number](rows, cols int) (*Matrix[T], error) {
	if rows < 1 || cols < 1 {
		return nil, fmt.Errorf("%w: %dx%d", ErrShape, rows, cols)
	}
	return &Matrix[T]{rows: rows, cols: cols, data: make([]T, rows*cols)}, nil
}

// FromRows builds a matrix from a copy of rows. Every row must have the same
// non-zero length.
func FromRows[T Number](rows [][]T) (*Matrix[T], error) {
	if len(rows) == 0 {
		return nil, fmt.Errorf("%w: no rows", ErrShape)
	}
	m, err := New[T](len(rows), len(rows[0]))
	if err != nil {
		return nil, err
	}
	for r, row := range rows {
		if len(row) != m.cols {
			return nil, fmt.Errorf("%w: row %d has %d values, want %d", ErrShape, r, len(row), m.cols)
		}
		copy(m.data[r*m.cols:], row)
	}
	return m, nil
}

// Identity returns the n × n identity matrix, or nil if n < 1.
func Identity[T Number](n int) *Matrix[T] {
	m, err := New[T](n, n)
	if err != nil {
		return nil
	}
	for i := range n {
		m.Set(i, i, 1)
	}
	return m
}

// Dims returns the number of rows and columns.
func (m *Matrix[T]) Dims() (rows, cols int) {
	return m.rows, m.cols
}

// At returns element (r, c).
func (m *Matrix[T]) At(r, c int) T {
	return m.data[r*m.cols+c]
}

// Set stores v at (r, c).
func (m *Matrix[T]) Set(r, c int, v T) {
	m.data[r*m.cols+c] = v
}

// Add returns m + o. Both must have the same dimensions.
func (m *Matrix[T]) Add(o *Matrix[T]) (*Matrix[T], error) {
	if m.rows != o.rows || m.cols != o.cols {
		return nil, fmt.Errorf("%w: %dx%d + %dx%d", ErrShape, m.rows, m.cols, o.rows, o.cols)
	}
	out, _ := New[T](m.rows, m.cols)
	for i := range m.data {
		out.data[i] = m.data[i] + o.data[i]
	}
	return out, nil
}

// Mul returns the matrix product m × o. m must have as many columns as o
// has rows.
func (m *Matrix[T]) Mul(o *Matrix[T]) (*Matrix[T], error) {
	if m.cols != o.rows {
		return nil, fmt.Errorf("%w: %dx%d * %dx%d", ErrShape, m.rows, m.cols, o.rows, o.cols)
	}
	out, _ := New[T](m.rows, o.cols)
	for r := range m.rows {
		for c := range o.cols {
			var sum T
			for k := range m.cols {
				sum += m.At(r, k) * o.At(k, c)
			}
			out.Set(r, c, sum)
		}
	}
	return out, nil
}

// Transpose returns a new cols × rows matrix with rows and columns swapped.
func (m *Matrix[T]) Transpose() *Matrix[T] {
	out, _ := New[T](m.cols, m.rows)
	for r := range m.rows {
		for c := range m.cols {
			out.Set(c, r, m.At(r, c))
		}
	}
	return out
}

// Sum returns the sum of all elements.
func (m *Matrix[T]) Sum() T {
	var total T
	for _, v := range m.data {
		total += v
	}
	return total
}

// Apply returns a new matrix of the same size holding f of every element.
// The element type may change.
func (m *Matrix[T]) Apply[U Number](f func(T) U) *Matrix[U] {
	out, _ := New[U](m.rows, m.cols)
	for i, v := range m.data {
		out.data[i] = f(v)
	}
	return out
}

// Rows yields each row index with a copy of that row.
func (m *Matrix[T]) Rows() iter.Seq2[int, []T] {
	return func(yield func(int, []T) bool) {
		for r := range m.rows {
			row := make([]T, m.cols)
			copy(row, m.data[r*m.cols:(r+1)*m.cols])
			if !yield(r, row) {
				return
			}
		}
	}
}

// String formats the matrix one row per line, like "[1 2]\n[3 4]".
func (m *Matrix[T]) String() string {
	lines := make([]string, 0, m.rows)
	for _, row := range m.Rows() {
		lines = append(lines, fmt.Sprint(row))
	}
	return strings.Join(lines, "\n")
}
