package tictactoe

// Parse converts three rows of "X", "O" and "." into a board.
func Parse(rows []string) ([3][3]rune, error) {
	// TODO
	return [3][3]rune{}, nil
}

// Winner returns 'X' or 'O' if that player has three in a row, '.' if
// nobody does, and an error if the board cannot occur in a real game.
func Winner(b [3][3]rune) (rune, error) {
	// TODO
	return 0, nil
}

// Status describes a legal board: "X wins", "O wins", "draw",
// "X to move" or "O to move".
func Status(b [3][3]rune) (string, error) {
	// TODO
	return "", nil
}
