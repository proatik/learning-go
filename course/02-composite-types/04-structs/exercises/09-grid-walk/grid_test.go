package grid

import (
	"fmt"
	"maps"
	"strings"
	"testing"
)

func TestMove(t *testing.T) {
	tests := []struct {
		p      Pos
		dir    rune
		want   Pos
		wantOK bool
	}{
		{Pos{2, 3}, 'U', Pos{2, 4}, true},
		{Pos{2, 3}, 'D', Pos{2, 2}, true},
		{Pos{2, 3}, 'R', Pos{3, 3}, true},
		{Pos{2, 3}, 'L', Pos{1, 3}, true},
		{Pos{}, 'D', Pos{0, -1}, true},
		{Pos{2, 3}, 'x', Pos{2, 3}, false},
		{Pos{2, 3}, 'u', Pos{2, 3}, false},
		{Pos{-1, -1}, ' ', Pos{-1, -1}, false},
	}
	for _, tt := range tests {
		got, ok := Move(tt.p, tt.dir)
		if got != tt.want || ok != tt.wantOK {
			t.Errorf("Move(%+v, %q) = %+v, %v; want %+v, %v", tt.p, tt.dir, got, ok, tt.want, tt.wantOK)
		}
	}
}

func TestRun(t *testing.T) {
	tests := []struct {
		path     string
		visits   map[Pos]int
		end      Pos
		min, max Pos
	}{
		{"", map[Pos]int{{0, 0}: 1}, Pos{}, Pos{}, Pos{}},
		{
			"UURDL",
			map[Pos]int{{0, 0}: 1, {0, 1}: 2, {0, 2}: 1, {1, 2}: 1, {1, 1}: 1},
			Pos{0, 1}, Pos{0, 0}, Pos{1, 2},
		},
		{
			"LDRRU",
			map[Pos]int{{0, 0}: 1, {-1, 0}: 1, {-1, -1}: 1, {0, -1}: 1, {1, -1}: 1, {1, 0}: 1},
			Pos{1, 0}, Pos{-1, -1}, Pos{1, 0},
		},
		{"RL", map[Pos]int{{0, 0}: 2, {1, 0}: 1}, Pos{0, 0}, Pos{0, 0}, Pos{1, 0}},
		{"UDUD", map[Pos]int{{0, 0}: 3, {0, 1}: 2}, Pos{0, 0}, Pos{0, 0}, Pos{0, 1}},
		{"DDD", map[Pos]int{{0, 0}: 1, {0, -1}: 1, {0, -2}: 1, {0, -3}: 1}, Pos{0, -3}, Pos{0, -3}, Pos{0, 0}},
	}
	for _, tt := range tests {
		w, err := Run(tt.path)
		if err != nil {
			t.Errorf("Run(%q) returned error %v", tt.path, err)
			continue
		}
		if !maps.Equal(w.Visits, tt.visits) {
			t.Errorf("Run(%q).Visits = %v, want %v", tt.path, w.Visits, tt.visits)
		}
		if w.End != tt.end || w.Min != tt.min || w.Max != tt.max {
			t.Errorf("Run(%q): End %+v, Min %+v, Max %+v; want End %+v, Min %+v, Max %+v",
				tt.path, w.End, w.Min, w.Max, tt.end, tt.min, tt.max)
		}
	}
}

func TestRunInvalid(t *testing.T) {
	tests := []struct {
		path  string
		char  rune
		index int
	}{
		{"UUX", 'X', 2},
		{"u", 'u', 0},
		{"R R", ' ', 1},
		{"RRRRL?", '?', 5},
	}
	for _, tt := range tests {
		w, err := Run(tt.path)
		if err == nil {
			t.Errorf("Run(%q): err = nil, want an error", tt.path)
			continue
		}
		msg := err.Error()
		if !strings.Contains(msg, fmt.Sprintf("%q", tt.char)) || !strings.Contains(msg, fmt.Sprint(tt.index)) {
			t.Errorf("Run(%q) error = %q, want it to mention %q and index %d", tt.path, msg, tt.char, tt.index)
		}
		if w.Visits != nil || w.End != (Pos{}) || w.Min != (Pos{}) || w.Max != (Pos{}) {
			t.Errorf("Run(%q) = %+v alongside the error, want Walk{}", tt.path, w)
		}
	}
}

func TestFirstRevisit(t *testing.T) {
	tests := []struct {
		path   string
		want   Pos
		wantOK bool
	}{
		{"UURDL", Pos{0, 1}, true},
		{"RL", Pos{0, 0}, true},
		{"UDUD", Pos{0, 0}, true},
		{"RRUULLDD", Pos{0, 0}, true},
		{"RRUULLDDRR", Pos{0, 0}, true},
		{"RRDLLUUL", Pos{0, 0}, true},
		{"RUULDDD", Pos{0, 0}, true},
		{"RRR", Pos{}, false},
		{"LDRRU", Pos{}, false},
		{"", Pos{}, false},
		{"RLX", Pos{}, false},
		{"UUX", Pos{}, false},
	}
	for _, tt := range tests {
		got, ok := FirstRevisit(tt.path)
		if got != tt.want || ok != tt.wantOK {
			t.Errorf("FirstRevisit(%q) = %+v, %v; want %+v, %v", tt.path, got, ok, tt.want, tt.wantOK)
		}
	}
}

func TestRender(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{"UURDL", "##\nE#\nS."},
		{"LDRRU", "#SE\n###"},
		{"RL", "S#"},
		{"UDUD", "#\nS"},
		{"RRUULLDD", "###\n#.#\nS##"},
		{"DDD", "S\n#\n#\nE"},
		{"", "S"},
		{"RRDLL", "S##\nE##"},
	}
	for _, tt := range tests {
		w, err := Run(tt.path)
		if err != nil {
			t.Fatalf("Run(%q) returned error %v", tt.path, err)
		}
		if got := Render(w); got != tt.want {
			t.Errorf("Render(Run(%q)) =\n%s\nwant\n%s", tt.path, got, tt.want)
		}
	}
	if got := Render(Walk{}); got != "" {
		t.Errorf("Render(Walk{}) = %q, want \"\"", got)
	}
}

func TestRenderHandmade(t *testing.T) {
	// Render must use only the Walk it is given.
	w := Walk{
		Visits: map[Pos]int{{0, 0}: 1, {2, 1}: 1, {1, 1}: 4},
		End:    Pos{2, 1},
		Min:    Pos{0, 0},
		Max:    Pos{2, 1},
	}
	want := ".#E\nS.."
	if got := Render(w); got != want {
		t.Errorf("Render(%+v) =\n%s\nwant\n%s", w, got, want)
	}
}
