package graph

import (
	"errors"
	"maps"
	"slices"
	"testing"
)

func sample() map[string][]string {
	return map[string][]string{
		"home":    {"park", "shop"},
		"park":    {"school"},
		"shop":    {"school", "mall"},
		"school":  {"library"},
		"library": {},
		"mall":    {"library", "shop"}, // cycle shop <-> mall
		"island":  {},
		"dock":    {"ferry"}, // "ferry" is only a neighbour, never a key
	}
}

func deepCopy(g map[string][]string) map[string][]string {
	out := make(map[string][]string, len(g))
	for k, v := range g {
		out[k] = slices.Clone(v)
	}
	return out
}

func deepEqual(a, b map[string][]string) bool {
	return maps.EqualFunc(a, b, func(x, y []string) bool { return slices.Equal(x, y) })
}

func TestAddEdge(t *testing.T) {
	g := map[string][]string{}
	AddEdge(g, "a", "b")
	AddEdge(g, "a", "c")
	AddEdge(g, "a", "b") // duplicate
	AddEdge(g, "b", "c")
	AddEdge(g, "c", "c") // self-loop
	AddEdge(g, "d", "a")

	want := map[string][]string{
		"a": {"b", "c"},
		"b": {"c"},
		"c": {"c"},
		"d": {"a"},
	}
	if !deepEqual(g, want) {
		t.Fatalf("after AddEdge calls g = %v, want %v", g, want)
	}

	g2 := map[string][]string{}
	AddEdge(g2, "x", "y")
	ys, ok := g2["y"]
	if !ok {
		t.Fatalf("AddEdge(g, \"x\", \"y\") did not add \"y\" as a key: %v", g2)
	}
	if len(ys) != 0 {
		t.Fatalf("g[\"y\"] = %v, want empty", ys)
	}
	AddEdge(g2, "y", "x")
	if !slices.Equal(g2["x"], []string{"y"}) || !slices.Equal(g2["y"], []string{"x"}) {
		t.Fatalf("AddEdge must not overwrite existing edges: %v", g2)
	}
}

func TestShortestPath(t *testing.T) {
	tests := []struct {
		from, to string
		want     []string
	}{
		{"home", "library", []string{"home", "park", "school", "library"}},
		{"home", "school", []string{"home", "park", "school"}},
		{"home", "mall", []string{"home", "shop", "mall"}},
		{"shop", "library", []string{"shop", "school", "library"}},
		{"mall", "school", []string{"mall", "shop", "school"}},
		{"home", "home", []string{"home"}},
		{"island", "island", []string{"island"}},
		{"ferry", "ferry", []string{"ferry"}},
		{"dock", "ferry", []string{"dock", "ferry"}},
	}
	for _, tt := range tests {
		g := sample()
		before := deepCopy(g)
		got, err := ShortestPath(g, tt.from, tt.to)
		if err != nil || !slices.Equal(got, tt.want) {
			t.Errorf("ShortestPath(%q, %q) = %v, %v, want %v, nil", tt.from, tt.to, got, err, tt.want)
		}
		if !deepEqual(g, before) {
			t.Errorf("ShortestPath(%q, %q) modified the graph", tt.from, tt.to)
		}
	}
}

func TestShortestPathTies(t *testing.T) {
	// Two equally short routes; the first-listed neighbour wins.
	g := map[string][]string{
		"s": {"b", "a"},
		"a": {"t"},
		"b": {"t"},
	}
	got, err := ShortestPath(g, "s", "t")
	if want := []string{"s", "b", "t"}; err != nil || !slices.Equal(got, want) {
		t.Fatalf("ShortestPath(s, t) = %v, %v, want %v, nil", got, err, want)
	}
	// A longer route discovered first must not win over a shorter one.
	g = map[string][]string{
		"s":  {"a", "t2"},
		"a":  {"b"},
		"b":  {"t"},
		"t2": {"t"},
	}
	got, err = ShortestPath(g, "s", "t")
	if want := []string{"s", "t2", "t"}; err != nil || !slices.Equal(got, want) {
		t.Fatalf("ShortestPath(s, t) = %v, %v, want %v, nil", got, err, want)
	}
}

func TestShortestPathErrors(t *testing.T) {
	tests := []struct {
		from, to string
		want     error
	}{
		{"library", "home", ErrNoPath},
		{"home", "island", ErrNoPath},
		{"island", "home", ErrNoPath},
		{"ferry", "dock", ErrNoPath},
		{"school", "shop", ErrNoPath},
		{"nowhere", "home", ErrUnknownNode},
		{"home", "nowhere", ErrUnknownNode},
		{"", "home", ErrUnknownNode},
		{"x", "y", ErrUnknownNode},
	}
	for _, tt := range tests {
		got, err := ShortestPath(sample(), tt.from, tt.to)
		if !errors.Is(err, tt.want) {
			t.Errorf("ShortestPath(%q, %q) error = %v, want %v", tt.from, tt.to, err, tt.want)
		}
		if got != nil {
			t.Errorf("ShortestPath(%q, %q) = %v, want nil on error", tt.from, tt.to, got)
		}
	}
	if _, err := ShortestPath(nil, "a", "a"); !errors.Is(err, ErrUnknownNode) {
		t.Errorf("ShortestPath(nil graph) error = %v, want ErrUnknownNode", err)
	}
}

func TestDistances(t *testing.T) {
	tests := []struct {
		from string
		want map[string]int
	}{
		{"home", map[string]int{"home": 0, "park": 1, "shop": 1, "school": 2, "mall": 2, "library": 3}},
		{"mall", map[string]int{"mall": 0, "library": 1, "shop": 1, "school": 2}},
		{"library", map[string]int{"library": 0}},
		{"dock", map[string]int{"dock": 0, "ferry": 1}},
		{"ferry", map[string]int{"ferry": 0}},
	}
	for _, tt := range tests {
		g := sample()
		before := deepCopy(g)
		got := Distances(g, tt.from)
		if !maps.Equal(got, tt.want) {
			t.Errorf("Distances(%q) = %v, want %v", tt.from, got, tt.want)
		}
		if !deepEqual(g, before) {
			t.Errorf("Distances(%q) modified the graph", tt.from)
		}
	}
	if got := Distances(sample(), "nowhere"); got != nil {
		t.Errorf("Distances(unknown) = %v, want nil", got)
	}
}

func TestLongChain(t *testing.T) {
	g := map[string][]string{}
	names := make([]string, 1000)
	for i := range names {
		names[i] = string(rune('a'+i%26)) + string(rune('0'+i/26%10)) + string(rune('0'+i/260))
	}
	for i := 1; i < len(names); i++ {
		AddEdge(g, names[i-1], names[i])
	}
	AddEdge(g, names[0], names[500]) // shortcut
	got, err := ShortestPath(g, names[0], names[999])
	if err != nil {
		t.Fatalf("ShortestPath(chain) error = %v", err)
	}
	if len(got) != 1+1+499 {
		t.Fatalf("ShortestPath(chain) has %d nodes, want %d", len(got), 501)
	}
	if got[0] != names[0] || got[1] != names[500] || got[len(got)-1] != names[999] {
		t.Fatalf("ShortestPath(chain) = %v..., want it to use the shortcut", got[:3])
	}
}
