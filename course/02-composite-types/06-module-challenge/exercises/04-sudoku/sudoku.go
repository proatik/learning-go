package sudoku

// Validate returns an error if g holds a value outside 0..9 or a digit
// repeated within a row, column or 3x3 box.
func Validate(g [9][9]int) error {
	// TODO
	return nil
}

// IsSolved reports whether g is valid and completely filled.
func IsSolved(g [9][9]int) bool {
	// TODO
	return false
}

// Candidates returns the digits that may legally be placed in the empty
// cell (row, col), in ascending order.
func Candidates(g [9][9]int, row, col int) []int {
	// TODO
	return nil
}

// Place writes v into g[row][col] if the move is legal.
func Place(g *[9][9]int, row, col, v int) error {
	// TODO
	return nil
}
