package matrix

import "errors"

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

// FromRows builds a matrix from a copy of rows; every row must have the same
// non-zero length.
func FromRows[T Number](rows [][]T) (*Matrix[T], error) {
	// TODO
	return nil, nil
}

// Identity returns the n × n identity matrix.
func Identity[T Number](n int) *Matrix[T] {
	// TODO
	return &Matrix[T]{}
}

// Dims returns the number of rows and columns.
func (m *Matrix[T]) Dims() (rows, cols int) {
	// TODO
	return 0, 0
}

// At returns the element in row i, column j.
func (m *Matrix[T]) At(i, j int) T {
	// TODO
	return 0
}

// Rows returns a copy of the matrix as a slice of rows.
func (m *Matrix[T]) Rows() [][]T {
	// TODO
	return nil
}

// Add returns m + o; both must have the same shape.
func (m *Matrix[T]) Add(o *Matrix[T]) (*Matrix[T], error) {
	// TODO
	return nil, nil
}

// Mul returns the matrix product m × o; m.cols must equal o.rows.
func (m *Matrix[T]) Mul(o *Matrix[T]) (*Matrix[T], error) {
	// TODO
	return nil, nil
}

// Scale returns m with every element multiplied by k.
func (m *Matrix[T]) Scale(k T) *Matrix[T] {
	// TODO
	return &Matrix[T]{}
}

// Transpose returns the cols × rows matrix with rows and columns swapped.
func (m *Matrix[T]) Transpose() *Matrix[T] {
	// TODO
	return &Matrix[T]{}
}
