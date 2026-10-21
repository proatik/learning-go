package matrix

// Matrix is a dense rows×cols matrix of float64 stored in row-major order.
// Create one with New, FromRows or Identity.
type Matrix struct {
	rows, cols int
	data       []float64 // element (i, j) is data[i*cols+j]
}

// New returns a rows×cols zero matrix, or nil if either size is not positive.
func New(rows, cols int) *Matrix {
	// TODO
	return nil
}

// FromRows returns a matrix holding a copy of rows, or nil if rows is empty,
// contains an empty row, or is ragged.
func FromRows(rows [][]float64) *Matrix {
	// TODO
	return nil
}

// Identity returns the n×n identity matrix, or nil if n is not positive.
func Identity(n int) *Matrix {
	// TODO
	return nil
}

// Rows returns the number of rows (0 for a nil matrix).
func (m *Matrix) Rows() int {
	// TODO
	return 0
}

// Cols returns the number of columns (0 for a nil matrix).
func (m *Matrix) Cols() int {
	// TODO
	return 0
}

// At returns element (i, j); ok is false if the indices are out of range.
func (m *Matrix) At(i, j int) (float64, bool) {
	// TODO
	return 0, false
}

// Set stores v at (i, j) and reports whether the indices were in range.
func (m *Matrix) Set(i, j int, v float64) bool {
	// TODO
	return false
}

// Equal reports whether m and o have the same shape and values.
func (m *Matrix) Equal(o *Matrix) bool {
	// TODO
	return false
}

// Scale multiplies every element by k in place.
func (m *Matrix) Scale(k float64) {
	// TODO
}

// Transpose returns a new matrix with rows and columns swapped.
func (m *Matrix) Transpose() *Matrix {
	// TODO
	return nil
}

// Add returns m + o, or nil if the shapes differ.
func (m *Matrix) Add(o *Matrix) *Matrix {
	// TODO
	return nil
}

// Mul returns the matrix product m × o, or nil if m.Cols() != o.Rows().
func (m *Matrix) Mul(o *Matrix) *Matrix {
	// TODO
	return nil
}

// String formats each row as "[v1 v2 ...]", rows separated by newlines.
func (m *Matrix) String() string {
	// TODO
	return ""
}
