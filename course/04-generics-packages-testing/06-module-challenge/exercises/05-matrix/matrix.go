// Package matrix provides a generic dense matrix of numbers.
package matrix

import (
	"errors"
	"iter"
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
	// TODO
	return nil, nil
}

// FromRows builds a matrix from a copy of rows. Every row must have the same
// non-zero length.
func FromRows[T Number](rows [][]T) (*Matrix[T], error) {
	// TODO
	return nil, nil
}

// Identity returns the n × n identity matrix, or nil if n < 1.
func Identity[T Number](n int) *Matrix[T] {
	// TODO
	return nil
}

// Dims returns the number of rows and columns.
func (m *Matrix[T]) Dims() (rows, cols int) {
	// TODO
	return 0, 0
}

// At returns element (r, c).
func (m *Matrix[T]) At(r, c int) T {
	// TODO
	var zero T
	return zero
}

// Set stores v at (r, c).
func (m *Matrix[T]) Set(r, c int, v T) {
	// TODO
}

// Add returns m + o. Both must have the same dimensions.
func (m *Matrix[T]) Add(o *Matrix[T]) (*Matrix[T], error) {
	// TODO
	return nil, nil
}

// Mul returns the matrix product m × o. m must have as many columns as o
// has rows.
func (m *Matrix[T]) Mul(o *Matrix[T]) (*Matrix[T], error) {
	// TODO
	return nil, nil
}

// Transpose returns a new cols × rows matrix with rows and columns swapped.
func (m *Matrix[T]) Transpose() *Matrix[T] {
	// TODO
	return nil
}

// Sum returns the sum of all elements.
func (m *Matrix[T]) Sum() T {
	// TODO
	var zero T
	return zero
}

// Apply returns a new matrix of the same size holding f of every element.
// The element type may change.
func (m *Matrix[T]) Apply[U Number](f func(T) U) *Matrix[U] {
	// TODO
	return nil
}

// Rows yields each row index with a copy of that row.
func (m *Matrix[T]) Rows() iter.Seq2[int, []T] {
	return func(yield func(int, []T) bool) {
		// TODO
	}
}

// String formats the matrix one row per line, like "[1 2]\n[3 4]".
func (m *Matrix[T]) String() string {
	// TODO
	return ""
}
