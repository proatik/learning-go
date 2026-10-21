package sudoku

import (
	"errors"
	"fmt"
)

// Validate returns an error if g holds a value outside 0..9 or a digit
// repeated within a row, column or 3x3 box.
func Validate(g [9][9]int) error {
	// seen[unit][digit] marks a digit as used in a row, column or box.
	// Arrays of bools are a compact fixed-size set.
	var rows, cols, boxes [9][10]bool
	for r := range 9 {
		for c := range 9 {
			v := g[r][c]
			if v < 0 || v > 9 {
				return fmt.Errorf("sudoku: invalid value %d at %d,%d", v, r, c)
			}
			if v == 0 {
				continue
			}
			b := r/3*3 + c/3
			switch {
			case rows[r][v]:
				return fmt.Errorf("sudoku: %d repeated in row %d", v, r)
			case cols[c][v]:
				return fmt.Errorf("sudoku: %d repeated in column %d", v, c)
			case boxes[b][v]:
				return fmt.Errorf("sudoku: %d repeated in box %d", v, b)
			}
			rows[r][v], cols[c][v], boxes[b][v] = true, true, true
		}
	}
	return nil
}

// IsSolved reports whether g is valid and completely filled.
func IsSolved(g [9][9]int) bool {
	for _, row := range g {
		for _, v := range row {
			if v == 0 {
				return false
			}
		}
	}
	return Validate(g) == nil
}

// used returns the set of digits 1..9 already present in the row, column
// and box of (row, col). Index 0 (empty) and invalid values are ignored.
func used(g [9][9]int, row, col int) [10]bool {
	var u [10]bool
	mark := func(v int) {
		if 1 <= v && v <= 9 {
			u[v] = true
		}
	}
	br, bc := row/3*3, col/3*3
	for i := range 9 {
		mark(g[row][i])         // row
		mark(g[i][col])         // column
		mark(g[br+i/3][bc+i%3]) // box, walked as i = 0..8
	}
	return u
}

func inRange(i int) bool { return 0 <= i && i <= 8 }

// Candidates returns the digits that may legally be placed in the empty
// cell (row, col), in ascending order.
func Candidates(g [9][9]int, row, col int) []int {
	if !inRange(row) || !inRange(col) || g[row][col] != 0 {
		return nil
	}
	u := used(g, row, col)
	var out []int
	for d := 1; d <= 9; d++ {
		if !u[d] {
			out = append(out, d)
		}
	}
	return out
}

// Place writes v into g[row][col] if the move is legal.
func Place(g *[9][9]int, row, col, v int) error {
	switch {
	case g == nil:
		return errors.New("sudoku: nil grid")
	case !inRange(row) || !inRange(col):
		return fmt.Errorf("sudoku: cell %d,%d out of range", row, col)
	case v < 1 || v > 9:
		return fmt.Errorf("sudoku: digit %d out of range", v)
	case g[row][col] != 0: // g[row] auto-dereferences the array pointer
		return fmt.Errorf("sudoku: cell %d,%d is not empty", row, col)
	}
	for _, d := range Candidates(*g, row, col) {
		if d == v {
			g[row][col] = v
			return nil
		}
	}
	return fmt.Errorf("sudoku: %d conflicts at %d,%d", v, row, col)
}
