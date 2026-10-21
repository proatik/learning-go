package tictactoe

import "testing"

// board builds a board directly, independent of Parse.
func board(r0, r1, r2 string) [3][3]rune {
	var b [3][3]rune
	for r, row := range []string{r0, r1, r2} {
		for c, ch := range []rune(row) {
			b[r][c] = ch
		}
	}
	return b
}

func TestParse(t *testing.T) {
	rows := []string{"XO.", ".X.", "O.X"}
	got, err := Parse(rows)
	if err != nil {
		t.Fatalf("Parse(%q) error = %v, want nil", rows, err)
	}
	if want := board("XO.", ".X.", "O.X"); got != want {
		t.Fatalf("Parse(%q) = %q, want %q", rows, got, want)
	}
	empty, err := Parse([]string{"...", "...", "..."})
	if err != nil || empty != board("...", "...", "...") {
		t.Fatalf("Parse(empty board) = %q, %v, want all '.', nil", empty, err)
	}
}

func TestParseInvalid(t *testing.T) {
	tests := [][]string{
		nil,
		{"...", "..."},
		{"...", "...", "...", "..."},
		{"..", "...", "..."},
		{"...", "....", "..."},
		{"...", "...", ""},
		{"x..", "...", "..."},
		{"...", ".0.", "..."},
		{"...", "...", ". ."},
		{"...", "...", "..é"},
	}
	for _, rows := range tests {
		got, err := Parse(rows)
		if err == nil {
			t.Errorf("Parse(%q) error = nil, want error", rows)
		}
		if got != [3][3]rune{} {
			t.Errorf("Parse(%q) = %q, want zero array on error", rows, got)
		}
	}
}

func TestWinner(t *testing.T) {
	tests := []struct {
		b    [3][3]rune
		want rune
	}{
		{board("...", "...", "..."), '.'},
		{board("X..", "...", "..."), '.'},
		{board("XO.", "...", "..."), '.'},
		{board("XXX", "OO.", "..."), 'X'},
		{board("OO.", "XXX", "..."), 'X'},
		{board("OO.", "...", "XXX"), 'X'},
		{board("XO.", "XO.", "X.."), 'X'},
		{board(".XO", ".XO", "X.."), '.'},
		{board("OX.", "OX.", ".X."), 'X'},
		{board("X.O", ".XO", "..X"), 'X'},
		{board("XXO", ".O.", "O.X"), 'O'},
		{board("OOO", "XX.", "X.."), 'O'},
		{board("XXO", "XO.", "O.."), 'O'},
		{board("XOX", "XOO", "OXX"), '.'}, // draw
		{board("XOX", "OXO", "XOX"), 'X'}, // two diagonals with the last move
		{board("XXX", "OXO", "OXO"), 'X'}, // row + column with the last move
	}
	for _, tt := range tests {
		got, err := Winner(tt.b)
		if err != nil || got != tt.want {
			t.Errorf("Winner(%q) = %q, %v, want %q, nil", tt.b, got, err, tt.want)
		}
	}
}

func TestWinnerImpossible(t *testing.T) {
	tests := []struct {
		name string
		b    [3][3]rune
	}{
		{"O moved first", board("O..", "...", "...")},
		{"X moved twice", board("XX.", "...", "...")},
		{"too many O", board("OOX", "...", "...")},
		{"both win", board("XXX", "OOO", "...")},
		{"O moved after X won", board("XXX", "OO.", "O..")},
		{"X moved after O won", board("OOO", "XX.", "X.X")},
		{"bad cell", board("X?.", "...", "...")},
		{"zero cell", [3][3]rune{}},
		{"lowercase", board("x..", "...", "...")},
	}
	for _, tt := range tests {
		got, err := Winner(tt.b)
		if err == nil {
			t.Errorf("%s: Winner(%q) = %q, nil, want error", tt.name, tt.b, got)
		} else if got != 0 {
			t.Errorf("%s: Winner(%q) = %q, want 0 with error", tt.name, tt.b, got)
		}
	}
}

func TestStatus(t *testing.T) {
	tests := []struct {
		b    [3][3]rune
		want string
	}{
		{board("...", "...", "..."), "X to move"},
		{board("X..", "...", "..."), "O to move"},
		{board("XO.", "...", "..."), "X to move"},
		{board("XXX", "OO.", "..."), "X wins"},
		{board("XXO", ".O.", "O.X"), "O wins"},
		{board("XOX", "XOO", "OXX"), "draw"},
		{board("XOX", "OXO", "XOX"), "X wins"},
		{board("XOX", "XOO", "OX."), "X to move"},
	}
	for _, tt := range tests {
		got, err := Status(tt.b)
		if err != nil || got != tt.want {
			t.Errorf("Status(%q) = %q, %v, want %q, nil", tt.b, got, err, tt.want)
		}
	}
}

func TestStatusImpossible(t *testing.T) {
	for _, b := range [][3][3]rune{
		board("XXX", "OOO", "..."),
		board("O..", "...", "..."),
		board("X#.", "...", "..."),
	} {
		if got, err := Status(b); err == nil || got != "" {
			t.Errorf("Status(%q) = %q, %v, want \"\", error", b, got, err)
		}
	}
}

func TestBoardIsAValue(t *testing.T) {
	b := board("X..", "...", "...")
	c := b // arrays copy
	c[1][1] = 'O'
	if b[1][1] != '.' {
		t.Fatalf("modifying a copy changed the original board")
	}
	if got, err := Status(c); err != nil || got != "X to move" {
		t.Fatalf("Status(%q) = %q, %v, want \"X to move\", nil", c, got, err)
	}
}
