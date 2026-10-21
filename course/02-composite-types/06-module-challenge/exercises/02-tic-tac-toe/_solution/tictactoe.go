package tictactoe

import (
	"errors"
	"fmt"
)

// lines lists the eight winning lines as [3][2]int{{row, col}, ...}.
var lines = [8][3][2]int{
	{{0, 0}, {0, 1}, {0, 2}}, {{1, 0}, {1, 1}, {1, 2}}, {{2, 0}, {2, 1}, {2, 2}}, // rows
	{{0, 0}, {1, 0}, {2, 0}}, {{0, 1}, {1, 1}, {2, 1}}, {{0, 2}, {1, 2}, {2, 2}}, // columns
	{{0, 0}, {1, 1}, {2, 2}}, {{0, 2}, {1, 1}, {2, 0}}, // diagonals
}

// Parse converts three rows of "X", "O" and "." into a board.
func Parse(rows []string) ([3][3]rune, error) {
	var b [3][3]rune
	if len(rows) != 3 {
		return [3][3]rune{}, fmt.Errorf("tictactoe: got %d rows, want 3", len(rows))
	}
	for r, row := range rows {
		cells := []rune(row)
		if len(cells) != 3 {
			return [3][3]rune{}, fmt.Errorf("tictactoe: row %d has %d cells, want 3", r, len(cells))
		}
		for c, ch := range cells {
			if ch != 'X' && ch != 'O' && ch != '.' {
				return [3][3]rune{}, fmt.Errorf("tictactoe: invalid cell %q at %d,%d", ch, r, c)
			}
			b[r][c] = ch
		}
	}
	return b, nil
}

// hasLine reports whether player owns any complete line.
func hasLine(b [3][3]rune, player rune) bool {
	for _, l := range lines {
		if b[l[0][0]][l[0][1]] == player &&
			b[l[1][0]][l[1][1]] == player &&
			b[l[2][0]][l[2][1]] == player {
			return true
		}
	}
	return false
}

// counts returns the number of X and O cells, or an error for a bad cell.
func counts(b [3][3]rune) (x, o int, err error) {
	for _, row := range b {
		for _, ch := range row {
			switch ch {
			case 'X':
				x++
			case 'O':
				o++
			case '.':
			default:
				return 0, 0, fmt.Errorf("tictactoe: invalid cell %q", ch)
			}
		}
	}
	return x, o, nil
}

// Winner returns 'X' or 'O' if that player has three in a row, '.' if
// nobody does, and an error if the board cannot occur in a real game.
func Winner(b [3][3]rune) (rune, error) {
	x, o, err := counts(b)
	if err != nil {
		return 0, err
	}
	if d := x - o; d != 0 && d != 1 {
		return 0, fmt.Errorf("tictactoe: impossible counts X=%d O=%d", x, o)
	}
	xWins, oWins := hasLine(b, 'X'), hasLine(b, 'O')
	switch {
	case xWins && oWins:
		return 0, errors.New("tictactoe: both players have a line")
	case xWins && x != o+1:
		return 0, errors.New("tictactoe: O moved after X won")
	case oWins && x != o:
		return 0, errors.New("tictactoe: X moved after O won")
	case xWins:
		return 'X', nil
	case oWins:
		return 'O', nil
	}
	return '.', nil
}

// Status describes a legal board: "X wins", "O wins", "draw",
// "X to move" or "O to move".
func Status(b [3][3]rune) (string, error) {
	w, err := Winner(b)
	if err != nil {
		return "", err
	}
	if w != '.' {
		return string(w) + " wins", nil
	}
	x, o, _ := counts(b) // already validated by Winner
	switch {
	case x+o == 9:
		return "draw", nil
	case x == o:
		return "X to move", nil
	default:
		return "O to move", nil
	}
}
