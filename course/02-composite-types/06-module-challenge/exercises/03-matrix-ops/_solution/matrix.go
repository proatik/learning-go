package matrix

import (
	"errors"
	"fmt"
)

// dims returns the number of rows and columns of m, or an error if m is
// ragged.
func dims(m [][]int) (rows, cols int, err error) {
	if len(m) == 0 {
		return 0, 0, nil
	}
	cols = len(m[0])
	for i, row := range m {
		if len(row) != cols {
			return 0, 0, fmt.Errorf("matrix: row %d has %d columns, want %d", i, len(row), cols)
		}
	}
	return len(m), cols, nil
}

// newMatrix allocates a rows×cols matrix of zeros.
func newMatrix(rows, cols int) [][]int {
	m := make([][]int, rows)
	for i := range m {
		m[i] = make([]int, cols)
	}
	return m
}

// Multiply returns the matrix product a×b.
func Multiply(a, b [][]int) ([][]int, error) {
	n, k, err := dims(a)
	if err != nil {
		return nil, err
	}
	k2, p, err := dims(b)
	if err != nil {
		return nil, err
	}
	if n == 0 || k == 0 || k2 == 0 || p == 0 {
		return nil, errors.New("matrix: empty matrix")
	}
	if k != k2 {
		return nil, fmt.Errorf("matrix: cannot multiply %dx%d by %dx%d", n, k, k2, p)
	}

	c := newMatrix(n, p)
	for i := range n {
		for j := range p {
			for x := range k {
				c[i][j] += a[i][x] * b[x][j]
			}
		}
	}
	return c, nil
}

// Rotate returns a new matrix: m rotated 90 degrees clockwise.
func Rotate(m [][]int) ([][]int, error) {
	rows, cols, err := dims(m)
	if err != nil {
		return nil, err
	}
	if rows == 0 || cols == 0 {
		return [][]int{}, nil
	}
	// Row r of the input becomes column rows-1-r of the output.
	out := newMatrix(cols, rows)
	for r := range rows {
		for c := range cols {
			out[c][rows-1-r] = m[r][c]
		}
	}
	return out, nil
}

// Spiral returns the elements of m in clockwise spiral order.
func Spiral(m [][]int) ([]int, error) {
	rows, cols, err := dims(m)
	if err != nil {
		return nil, err
	}
	out := make([]int, 0, rows*cols)
	// Shrink the boundaries after walking each side.
	top, bottom, left, right := 0, rows-1, 0, cols-1
	for top <= bottom && left <= right {
		for c := left; c <= right; c++ {
			out = append(out, m[top][c])
		}
		top++
		for r := top; r <= bottom; r++ {
			out = append(out, m[r][right])
		}
		right--
		if top <= bottom { // a remaining row to walk back along
			for c := right; c >= left; c-- {
				out = append(out, m[bottom][c])
			}
			bottom--
		}
		if left <= right { // a remaining column to walk up
			for r := bottom; r >= top; r-- {
				out = append(out, m[r][left])
			}
			left++
		}
	}
	return out, nil
}
