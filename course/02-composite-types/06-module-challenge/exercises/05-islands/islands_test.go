package islands

import (
	"slices"
	"strings"
	"testing"
)

var sample = []string{
	"##..#",
	"#...#",
	"..#..",
	".....",
	"###.#",
}

func TestIslands(t *testing.T) {
	tests := []struct {
		name  string
		grid  []string
		sizes []int
	}{
		{"sample", sample, []int{3, 3, 2, 1, 1}},
		{"nil", nil, []int{}},
		{"empty rows", []string{"", ""}, []int{}},
		{"all water", []string{"...", "...", "..."}, []int{}},
		{"all land", []string{"###", "###"}, []int{6}},
		{"single", []string{"#"}, []int{1}},
		{"diagonals do not connect", []string{"#.#", ".#.", "#.#"}, []int{1, 1, 1, 1, 1}},
		{"other bytes are water", []string{"#~#", "~~~", "#X#"}, []int{1, 1, 1, 1}},
		{"ring", []string{"#####", "#...#", "#.#.#", "#...#", "#####"}, []int{16, 1}},
		{"snake", []string{"#####", "....#", "#####", "#....", "#####"}, []int{17}},
		{"ragged rows", []string{"##", "#####", "", "#", "..#"}, []int{7, 1, 1}},
		{"ragged connects", []string{"#", "#.#", "###"}, []int{6}},
	}
	for _, tt := range tests {
		if got := CountIslands(tt.grid); got != len(tt.sizes) {
			t.Errorf("CountIslands(%s) = %d, want %d", tt.name, got, len(tt.sizes))
		}
		got := IslandSizes(tt.grid)
		if !slices.Equal(got, tt.sizes) {
			t.Errorf("IslandSizes(%s) = %v, want %v", tt.name, got, tt.sizes)
		}
	}
}

func TestIslandsLarge(t *testing.T) {
	// A 300x300 checkerboard of 2x2 land blocks separated by water.
	var grid []string
	for r := range 300 {
		var sb strings.Builder
		for c := range 300 {
			if r%3 != 2 && c%3 != 2 {
				sb.WriteByte('#')
			} else {
				sb.WriteByte('.')
			}
		}
		grid = append(grid, sb.String())
	}
	if got := CountIslands(grid); got != 100*100 {
		t.Fatalf("CountIslands(300x300 blocks) = %d, want %d", got, 100*100)
	}
	full := slices.Repeat([]string{strings.Repeat("#", 300)}, 300)
	if got := IslandSizes(full); !slices.Equal(got, []int{90000}) {
		t.Fatalf("IslandSizes(300x300 land) = %v, want [90000]", got)
	}
}

func toBytes(rows ...string) [][]byte {
	g := make([][]byte, len(rows))
	for i, r := range rows {
		g[i] = []byte(r)
	}
	return g
}

func toStrings(g [][]byte) []string {
	out := make([]string, len(g))
	for i, r := range g {
		out[i] = string(r)
	}
	return out
}

func TestFloodFill(t *testing.T) {
	tests := []struct {
		name     string
		grid     []string
		row, col int
		fill     byte
		want     []string
		changed  int
	}{
		{"example", []string{"aab", "abb", "bba"}, 0, 2, 'x', []string{"aax", "axx", "xxa"}, 5},
		{"other region", []string{"aab", "abb", "bba"}, 0, 0, 'x', []string{"xxb", "xbb", "bba"}, 3},
		{"isolated corner", []string{"aab", "abb", "bba"}, 2, 2, 'z', []string{"aab", "abb", "bbz"}, 1},
		{"whole grid", []string{"...", "..."}, 1, 2, '#', []string{"###", "###"}, 6},
		{"no diagonals", []string{"a.a", ".a.", "a.a"}, 1, 1, 'b', []string{"a.a", ".b.", "a.a"}, 1},
		{"ring interior untouched", []string{"#####", "#...#", "#####"}, 0, 0, 'o', []string{"ooooo", "o...o", "ooooo"}, 12},
		{"ragged", []string{"aa", "aaaa", "", "a"}, 0, 0, 'b', []string{"bb", "bbbb", "", "a"}, 6},
	}
	for _, tt := range tests {
		g := toBytes(tt.grid...)
		changed := FloodFill(g, tt.row, tt.col, tt.fill)
		if changed != tt.changed {
			t.Errorf("FloodFill(%s) changed %d cells, want %d", tt.name, changed, tt.changed)
		}
		if got := toStrings(g); !slices.Equal(got, tt.want) {
			t.Errorf("FloodFill(%s) grid = %q, want %q", tt.name, got, tt.want)
		}
	}
}

func TestFloodFillNoOp(t *testing.T) {
	grid := []string{"aab", "abb", "bba"}
	tests := []struct {
		name     string
		row, col int
		fill     byte
	}{
		{"same colour", 0, 0, 'a'},
		{"row -1", -1, 0, 'x'},
		{"row 3", 3, 0, 'x'},
		{"col -1", 0, -1, 'x'},
		{"col 3", 1, 3, 'x'},
	}
	for _, tt := range tests {
		g := toBytes(grid...)
		if n := FloodFill(g, tt.row, tt.col, tt.fill); n != 0 {
			t.Errorf("FloodFill(%s) = %d, want 0", tt.name, n)
		}
		if got := toStrings(g); !slices.Equal(got, grid) {
			t.Errorf("FloodFill(%s) changed the grid to %q", tt.name, got)
		}
	}
	if n := FloodFill(nil, 0, 0, 'x'); n != 0 {
		t.Errorf("FloodFill(nil, 0, 0) = %d, want 0", n)
	}
	ragged := toBytes("a", "aaa")
	if n := FloodFill(ragged, 0, 2, 'x'); n != 0 {
		t.Errorf("FloodFill(cell past the end of a short row) = %d, want 0", n)
	}
}

func TestFloodFillLarge(t *testing.T) {
	g := make([][]byte, 300)
	for i := range g {
		g[i] = []byte(strings.Repeat(".", 300))
	}
	if n := FloodFill(g, 150, 150, '#'); n != 90000 {
		t.Fatalf("FloodFill(300x300) = %d, want 90000", n)
	}
	if g[0][0] != '#' || g[299][299] != '#' {
		t.Fatalf("FloodFill(300x300) did not reach the corners")
	}
}
