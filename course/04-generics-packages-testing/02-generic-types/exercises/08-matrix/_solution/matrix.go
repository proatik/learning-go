package matrix

import (
	"errors"
	"fmt"
)

// Number is the set of element types a Matrix supports. It permits + and *.
type Number interface {
	~int | ~int32 | ~int64 | ~float32 | ~float64
}

// Matrix is a rows × cols matrix of numbers.
type Matrix[T Number] struct {
	rows, cols int
	data       []T // row-major: element (i, j) is data[i*cols+j]
}

// ErrShape reports matrices whose dimensions don't fit the operation.
var ErrShape = errors.New("matrix: shape mismatch")

// newMatrix returns a rows × cols matrix of zeros.
func newMatrix[T Number](rows, cols int) *Matrix[T] {
	return &Matrix[T]{rows: rows, cols: cols, data: make([]T, rows*cols)}
}

// FromRows builds a matrix from a copy of rows; every row must have the same
// non-zero length.
func FromRows[T Number](rows [][]T) (*Matrix[T], error) {
	if len(rows) == 0 || len(rows[0]) == 0 {
		return nil, ErrShape
	}
	m := newMatrix[T](len(rows), len(rows[0]))
	for i, row := range rows {
		if len(row) != m.cols {
			return nil, ErrShape
		}
		copy(m.data[i*m.cols:], row)
	}
	return m, nil
}

// Identity returns the n × n identity matrix.
func Identity[T Number](n int) *Matrix[T] {
	n = max(n, 0)
	m := newMatrix[T](n, n)
	for i := range n {
		m.data[i*n+i] = 1
	}
	return m
}

// Dims returns the number of rows and columns.
func (m *Matrix[T]) Dims() (rows, cols int) {
	return m.rows, m.cols
}

// At returns the element in row i, column j.
func (m *Matrix[T]) At(i, j int) T {
	return m.data[i*m.cols+j]
}

// Rows returns a copy of the matrix as a slice of rows.
func (m *Matrix[T]) Rows() [][]T {
	out := make([][]T, m.rows)
	for i := range out {
		out[i] = append([]T(nil), m.data[i*m.cols:(i+1)*m.cols]...)
	}
	return out
}

// Add returns m + o; both must have the same shape.
func (m *Matrix[T]) Add(o *Matrix[T]) (*Matrix[T], error) {
	if m.rows != o.rows || m.cols != o.cols {
		return nil, m.shapeError("add", o)
	}
	out := newMatrix[T](m.rows, m.cols)
	for i := range out.data {
		out.data[i] = m.data[i] + o.data[i]
	}
	return out, nil
}

// Mul returns the matrix product m × o; m.cols must equal o.rows.
func (m *Matrix[T]) Mul(o *Matrix[T]) (*Matrix[T], error) {
	if m.cols != o.rows {
		return nil, m.shapeError("mul", o)
	}
	out := newMatrix[T](m.rows, o.cols)
	for i := range m.rows {
		for j := range o.cols {
			var sum T
			for k := range m.cols {
				sum += m.At(i, k) * o.At(k, j)
			}
			out.data[i*out.cols+j] = sum
		}
	}
	return out, nil
}

// Scale returns m with every element multiplied by k.
func (m *Matrix[T]) Scale(k T) *Matrix[T] {
	out := newMatrix[T](m.rows, m.cols)
	for i, v := range m.data {
		out.data[i] = v * k
	}
	return out
}

// Transpose returns the cols × rows matrix with rows and columns swapped.
func (m *Matrix[T]) Transpose() *Matrix[T] {
	out := newMatrix[T](m.cols, m.rows)
	for i := range m.rows {
		for j := range m.cols {
			out.data[j*out.cols+i] = m.At(i, j)
		}
	}
	return out
}

// shapeError describes why op can't combine m and o.
func (m *Matrix[T]) shapeError(op string, o *Matrix[T]) error {
	return fmt.Errorf("%s %dx%d and %dx%d: %w", op, m.rows, m.cols, o.rows, o.cols, ErrShape)
}
