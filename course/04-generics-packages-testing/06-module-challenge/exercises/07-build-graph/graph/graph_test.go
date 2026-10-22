package graph

import (
	"cmp"
	"errors"
	"fmt"
	"slices"
	"testing"
)

func newGraph[N cmp.Ordered](t *testing.T, edges ...[2]N) *Graph[N] {
	t.Helper()
	g := New[N]()
	if g == nil {
		t.Fatalf("New() = nil, want an empty graph")
	}
	for _, e := range edges {
		g.AddEdge(e[0], e[1])
	}
	return g
}

// bfs renders g.BFS(start) as "node:dist" strings.
func bfs[N cmp.Ordered](g *Graph[N], start N) []string {
	var out []string
	for n, d := range g.BFS(start) {
		out = append(out, fmt.Sprintf("%v:%d", n, d))
	}
	return out
}

func TestNodesAndNeighbors(t *testing.T) {
	g := newGraph(t, [2]string{"b", "c"}, [2]string{"b", "a"}, [2]string{"b", "c"}, [2]string{"d", "b"})
	g.AddNode("z")
	g.AddNode("a")
	if got, want := slices.Collect(g.Nodes()), []string{"a", "b", "c", "d", "z"}; !slices.Equal(got, want) {
		t.Errorf("Nodes() = %q, want %q", got, want)
	}
	if got, want := slices.Collect(g.Neighbors("b")), []string{"c", "a"}; !slices.Equal(got, want) {
		t.Errorf("Neighbors(%q) = %q, want %q (insertion order, no duplicates)", "b", got, want)
	}
	if got := slices.Collect(g.Neighbors("z")); len(got) != 0 {
		t.Errorf("Neighbors(%q) = %q, want none", "z", got)
	}
	if got := slices.Collect(g.Neighbors("missing")); len(got) != 0 {
		t.Errorf("Neighbors(%q) = %q, want none", "missing", got)
	}
}

func TestBFS(t *testing.T) {
	g := newGraph(t,
		[2]string{"a", "b"}, [2]string{"a", "c"}, [2]string{"b", "d"},
		[2]string{"c", "d"}, [2]string{"d", "e"}, [2]string{"e", "a"},
	)
	g.AddNode("x")
	tests := []struct {
		start string
		want  []string
	}{
		{"a", []string{"a:0", "b:1", "c:1", "d:2", "e:3"}},
		{"d", []string{"d:0", "e:1", "a:2", "b:3", "c:3"}},
		{"x", []string{"x:0"}},
		{"nope", nil},
	}
	for _, tt := range tests {
		if got := bfs(g, tt.start); !slices.Equal(got, tt.want) {
			t.Errorf("BFS(%q) = %q, want %q", tt.start, got, tt.want)
		}
	}

	var got []string
	for n := range g.BFS("a") {
		got = append(got, n)
		if n == "c" {
			break
		}
	}
	if want := []string{"a", "b", "c"}; !slices.Equal(got, want) {
		t.Errorf("BFS(%q) with break at %q = %q, want %q", "a", "c", got, want)
	}
}

func TestTopoSort(t *testing.T) {
	g := newGraph(t,
		[2]string{"shirt", "tie"}, [2]string{"tie", "jacket"}, [2]string{"trousers", "shoes"},
		[2]string{"trousers", "belt"}, [2]string{"belt", "jacket"}, [2]string{"shirt", "belt"},
		[2]string{"socks", "shoes"},
	)
	g.AddNode("watch")
	got, err := g.TopoSort()
	want := []string{"shirt", "socks", "tie", "trousers", "belt", "jacket", "shoes", "watch"}
	if err != nil || !slices.Equal(got, want) {
		t.Errorf("TopoSort() = %q, %v, want %q, <nil>", got, err, want)
	}

	ints := newGraph(t, [2]int{3, 1}, [2]int{2, 1}, [2]int{1, 0}, [2]int{10, 0})
	gotInts, err := ints.TopoSort()
	if wantInts := []int{2, 3, 1, 10, 0}; err != nil || !slices.Equal(gotInts, wantInts) {
		t.Errorf("TopoSort() = %v, %v, want %v, <nil>", gotInts, err, wantInts)
	}

	empty := newGraph[string](t)
	if got, err := empty.TopoSort(); err != nil || len(got) != 0 {
		t.Errorf("empty TopoSort() = %q, %v, want [], <nil>", got, err)
	}
}

func TestTopoSortCycle(t *testing.T) {
	tests := []struct {
		name    string
		edges   [][2]string
		want    []string
		wantMsg string
	}{
		{
			"triangle with a tail",
			[][2]string{{"a", "b"}, {"b", "c"}, {"c", "a"}, {"c", "d"}, {"e", "a"}},
			[]string{"a", "b", "c", "d"},
			"graph: cycle among [a b c d]",
		},
		{"self loop", [][2]string{{"x", "x"}, {"w", "y"}}, []string{"x"}, "graph: cycle among [x]"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := newGraph(t, tt.edges...)
			order, err := g.TopoSort()
			if order != nil {
				t.Errorf("TopoSort() order = %q, want nil", order)
			}
			if !errors.Is(err, ErrCycle) {
				t.Fatalf("TopoSort() error = %v, want %v", err, ErrCycle)
			}
			var ce *CycleError[string]
			if !errors.As(err, &ce) {
				t.Fatalf("errors.As(%v, *CycleError[string]) = false, want true", err)
			}
			if !slices.Equal(ce.Nodes, tt.want) {
				t.Errorf("CycleError.Nodes = %q, want %q", ce.Nodes, tt.want)
			}
			if err.Error() != tt.wantMsg {
				t.Errorf("TopoSort() error = %q, want %q", err.Error(), tt.wantMsg)
			}
		})
	}
}
