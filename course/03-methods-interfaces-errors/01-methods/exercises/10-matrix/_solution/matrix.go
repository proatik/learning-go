package matrix

import (
	"slices"
	"strconv"
	"strings"
)

// Matrix is a dense rows×cols matrix of float64 stored in row-major order.
// Create one with New, FromRows or Identity.
type Matrix struct {
	rows, cols int
	data       []float64 // element (i, j) is data[i*cols+j]
}

// New returns a rows×cols zero matrix, or nil if either size is not positive.
func New(rows, cols int) *Matrix {
	if rows <= 0 || cols <= 0 {
		return nil
	}
	return &Matrix{rows: rows, cols: cols, data: make([]float64, rows*cols)}
}

// FromRows returns a matrix holding a copy of rows, or nil if rows is empty,
// contains an empty row, or is ragged.
func FromRows(rows [][]float64) *Matrix {
	if len(rows) == 0 {
		return nil
	}
	m := New(len(rows), len(rows[0]))
	if m == nil {
		return nil
	}
	for i, row := range rows {
		if len(row) != m.cols {
			return nil
		}
		copy(m.data[i*m.cols:], row)
	}
	return m
}

// Identity returns the n×n identity matrix, or nil if n is not positive.
func Identity(n int) *Matrix {
	m := New(n, n)
	if m == nil {
		return nil
	}
	for i := range n {
		m.data[i*n+i] = 1
	}
	return m
}

// Rows returns the number of rows (0 for a nil matrix).
func (m *Matrix) Rows() int {
	if m == nil {
		return 0
	}
	return m.rows
}

// Cols returns the number of columns (0 for a nil matrix).
func (m *Matrix) Cols() int {
	if m == nil {
		return 0
	}
	return m.cols
}

// inRange reports whether (i, j) is a valid index of m.
func (m *Matrix) inRange(i, j int) bool {
	return i >= 0 && i < m.Rows() && j >= 0 && j < m.Cols()
}

// At returns element (i, j); ok is false if the indices are out of range.
func (m *Matrix) At(i, j int) (float64, bool) {
	if !m.inRange(i, j) {
		return 0, false
	}
	return m.data[i*m.cols+j], true
}

// Set stores v at (i, j) and reports whether the indices were in range.
func (m *Matrix) Set(i, j int, v float64) bool {
	if !m.inRange(i, j) {
		return false
	}
	m.data[i*m.cols+j] = v
	return true
}

// Equal reports whether m and o have the same shape and values.
func (m *Matrix) Equal(o *Matrix) bool {
	if m == nil || o == nil {
		return m == o
	}
	return m.rows == o.rows && m.cols == o.cols && slices.Equal(m.data, o.data)
}

// Scale multiplies every element by k in place.
func (m *Matrix) Scale(k float64) {
	if m == nil {
		return
	}
	for i := range m.data {
		m.data[i] *= k
	}
}

// Transpose returns a new matrix with rows and columns swapped.
func (m *Matrix) Transpose() *Matrix {
	t := New(m.Cols(), m.Rows())
	if t == nil {
		return nil
	}
	for i := range m.rows {
		for j := range m.cols {
			t.data[j*t.cols+i] = m.data[i*m.cols+j]
		}
	}
	return t
}

// Add returns m + o, or nil if the shapes differ.
func (m *Matrix) Add(o *Matrix) *Matrix {
	if m == nil || o == nil || m.rows != o.rows || m.cols != o.cols {
		return nil
	}
	sum := New(m.rows, m.cols)
	for i := range sum.data {
		sum.data[i] = m.data[i] + o.data[i]
	}
	return sum
}

// Mul returns the matrix product m × o, or nil if m.Cols() != o.Rows().
func (m *Matrix) Mul(o *Matrix) *Matrix {
	if m == nil || o == nil || m.cols != o.rows {
		return nil
	}
	p := New(m.rows, o.cols)
	for i := range m.rows {
		for j := range o.cols {
			var s float64
			for k := range m.cols {
				s += m.data[i*m.cols+k] * o.data[k*o.cols+j]
			}
			p.data[i*p.cols+j] = s
		}
	}
	return p
}

// String formats each row as "[v1 v2 ...]", rows separated by newlines.
func (m *Matrix) String() string {
	lines := make([]string, m.Rows())
	for i := range lines {
		vals := make([]string, m.cols)
		for j := range vals {
			vals[j] = strconv.FormatFloat(m.data[i*m.cols+j], 'g', -1, 64)
		}
		lines[i] = "[" + strings.Join(vals, " ") + "]"
	}
	return strings.Join(lines, "\n")
}
