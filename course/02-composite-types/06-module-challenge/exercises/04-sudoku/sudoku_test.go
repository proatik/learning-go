package sudoku

import (
	"slices"
	"testing"
)

var solved = [9][9]int{
	{5, 3, 4, 6, 7, 8, 9, 1, 2},
	{6, 7, 2, 1, 9, 5, 3, 4, 8},
	{1, 9, 8, 3, 4, 2, 5, 6, 7},
	{8, 5, 9, 7, 6, 1, 4, 2, 3},
	{4, 2, 6, 8, 5, 3, 7, 9, 1},
	{7, 1, 3, 9, 2, 4, 8, 5, 6},
	{9, 6, 1, 5, 3, 7, 2, 8, 4},
	{2, 8, 7, 4, 1, 9, 6, 3, 5},
	{3, 4, 5, 2, 8, 6, 1, 7, 9},
}

var puzzle = [9][9]int{
	{5, 3, 0, 0, 7, 0, 0, 0, 0},
	{6, 0, 0, 1, 9, 5, 0, 0, 0},
	{0, 9, 8, 0, 0, 0, 0, 6, 0},
	{8, 0, 0, 0, 6, 0, 0, 0, 3},
	{4, 0, 0, 8, 0, 3, 0, 0, 1},
	{7, 0, 0, 0, 2, 0, 0, 0, 6},
	{0, 6, 0, 0, 0, 0, 2, 8, 0},
	{0, 0, 0, 4, 1, 9, 0, 0, 5},
	{0, 0, 0, 0, 8, 0, 0, 7, 9},
}

func with(g [9][9]int, r, c, v int) [9][9]int {
	g[r][c] = v // g is a copy
	return g
}

func TestValidate(t *testing.T) {
	valid := map[string][9][9]int{
		"empty":    {},
		"solved":   solved,
		"puzzle":   puzzle,
		"one cell": with([9][9]int{}, 4, 4, 9),
	}
	for name, g := range valid {
		if err := Validate(g); err != nil {
			t.Errorf("Validate(%s) = %v, want nil", name, err)
		}
	}

	invalid := map[string][9][9]int{
		"row dup":      with(with([9][9]int{}, 0, 0, 3), 0, 8, 3),
		"column dup":   with(with([9][9]int{}, 0, 5, 7), 8, 5, 7),
		"box dup":      with(with([9][9]int{}, 3, 3, 2), 5, 5, 2),
		"box dup 2":    with(with([9][9]int{}, 6, 8, 1), 8, 6, 1),
		"negative":     with([9][9]int{}, 2, 2, -1),
		"ten":          with([9][9]int{}, 8, 8, 10),
		"solved swap":  with(solved, 0, 0, 3),
		"puzzle clash": with(puzzle, 0, 2, 5),
	}
	for name, g := range invalid {
		if err := Validate(g); err == nil {
			t.Errorf("Validate(%s) = nil, want error", name)
		}
	}
}

func TestIsSolved(t *testing.T) {
	tests := []struct {
		name string
		g    [9][9]int
		want bool
	}{
		{"solved", solved, true},
		{"puzzle", puzzle, false},
		{"empty", [9][9]int{}, false},
		{"one hole", with(solved, 4, 4, 0), false},
		{"full but wrong", with(solved, 0, 0, 3), false},
		{"full but out of range", with(solved, 8, 8, 10), false},
	}
	for _, tt := range tests {
		if got := IsSolved(tt.g); got != tt.want {
			t.Errorf("IsSolved(%s) = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestCandidates(t *testing.T) {
	tests := []struct {
		name     string
		g        [9][9]int
		row, col int
		want     []int
	}{
		{"empty grid", [9][9]int{}, 0, 0, []int{1, 2, 3, 4, 5, 6, 7, 8, 9}},
		{"one placed", with([9][9]int{}, 0, 0, 5), 0, 1, []int{1, 2, 3, 4, 6, 7, 8, 9}},
		{"puzzle 0,2", puzzle, 0, 2, []int{1, 2, 4}},
		{"puzzle 4,4", puzzle, 4, 4, []int{5}},
		{"puzzle 8,0", puzzle, 8, 0, []int{1, 2, 3}},
		{"solved hole", with(solved, 4, 4, 0), 4, 4, []int{5}},
	}
	for _, tt := range tests {
		got := Candidates(tt.g, tt.row, tt.col)
		if !slices.Equal(got, tt.want) || len(got) != len(tt.want) {
			t.Errorf("Candidates(%s, %d, %d) = %v, want %v", tt.name, tt.row, tt.col, got, tt.want)
		}
	}
}

func TestCandidatesNone(t *testing.T) {
	tests := []struct {
		name     string
		row, col int
	}{
		{"filled cell", 0, 0},
		{"row -1", -1, 0},
		{"row 9", 9, 0},
		{"col -1", 0, -1},
		{"col 9", 0, 9},
	}
	for _, tt := range tests {
		if got := Candidates(puzzle, tt.row, tt.col); len(got) != 0 {
			t.Errorf("Candidates(puzzle, %d, %d) [%s] = %v, want empty", tt.row, tt.col, tt.name, got)
		}
	}
}

func TestPlace(t *testing.T) {
	var g [9][9]int
	if err := Place(&g, 0, 0, 5); err != nil {
		t.Fatalf("Place(&g, 0, 0, 5) = %v, want nil", err)
	}
	if g[0][0] != 5 {
		t.Fatalf("after Place, g[0][0] = %d, want 5 (write through the pointer)", g[0][0])
	}
	if err := Place(&g, 8, 8, 9); err != nil || g[8][8] != 9 {
		t.Fatalf("Place(&g, 8, 8, 9) = %v, g[8][8] = %d, want nil, 9", err, g[8][8])
	}

	before := g
	bad := []struct {
		name        string
		row, col, v int
	}{
		{"row conflict", 0, 8, 5},
		{"column conflict", 4, 0, 5},
		{"box conflict", 2, 2, 5},
		{"occupied", 0, 0, 1},
		{"digit 0", 1, 1, 0},
		{"digit 10", 1, 1, 10},
		{"digit -3", 1, 1, -3},
		{"row 9", 9, 0, 1},
		{"col -1", 0, -1, 1},
	}
	for _, tt := range bad {
		if err := Place(&g, tt.row, tt.col, tt.v); err == nil {
			t.Errorf("Place(&g, %d, %d, %d) [%s] = nil, want error", tt.row, tt.col, tt.v, tt.name)
		}
		if g != before {
			t.Fatalf("Place(&g, %d, %d, %d) [%s] changed the grid on error", tt.row, tt.col, tt.v, tt.name)
		}
	}

	if err := Place(nil, 0, 0, 1); err == nil {
		t.Errorf("Place(nil, 0, 0, 1) = nil, want error")
	}
}

func TestPlaceSolvesPuzzle(t *testing.T) {
	g := puzzle
	for r := range 9 {
		for c := range 9 {
			if g[r][c] == 0 {
				if err := Place(&g, r, c, solved[r][c]); err != nil {
					t.Fatalf("Place(&g, %d, %d, %d) = %v, want nil", r, c, solved[r][c], err)
				}
			}
		}
	}
	if g != solved {
		t.Fatalf("filled puzzle does not match the solution")
	}
	if !IsSolved(g) {
		t.Fatalf("IsSolved(filled puzzle) = false, want true")
	}
	if puzzle[0][2] != 0 {
		t.Fatalf("Place on a copy modified the original puzzle")
	}
}
