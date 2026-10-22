// Package graph provides a generic directed graph with BFS and topological
// sorting.
package graph

import (
	"cmp"
	"errors"
	"fmt"
	"iter"
	"maps"
	"slices"
)

// ErrCycle matches every *CycleError.
var ErrCycle = errors.New("graph: cycle")

// CycleError reports the nodes that could not be ordered because they are
// on, or depend on, a cycle.
type CycleError[N cmp.Ordered] struct {
	Nodes []N // sorted
}

// Error returns "graph: cycle among [a b]".
func (e *CycleError[N]) Error() string {
	return fmt.Sprintf("graph: cycle among %v", e.Nodes)
}

// Is makes errors.Is(err, ErrCycle) true for any *CycleError.
func (e *CycleError[N]) Is(target error) bool {
	return target == ErrCycle
}

// Graph is a directed graph. Create one with New.
type Graph[N cmp.Ordered] struct {
	adj map[N][]N // node → successors in insertion order; every node is a key
}

// New returns an empty graph.
func New[N cmp.Ordered]() *Graph[N] {
	return &Graph[N]{adj: make(map[N][]N)}
}

// AddNode adds n if it is not already in the graph.
func (g *Graph[N]) AddNode(n N) {
	if _, ok := g.adj[n]; !ok {
		g.adj[n] = nil
	}
}

// AddEdge adds the edge from → to, adding both nodes. A repeated edge is
// ignored.
func (g *Graph[N]) AddEdge(from, to N) {
	g.AddNode(from)
	g.AddNode(to)
	if !slices.Contains(g.adj[from], to) {
		g.adj[from] = append(g.adj[from], to)
	}
}

// Nodes yields every node in sorted order.
func (g *Graph[N]) Nodes() iter.Seq[N] {
	return slices.Values(slices.Sorted(maps.Keys(g.adj)))
}

// Neighbors yields the successors of n in the order their edges were added.
func (g *Graph[N]) Neighbors(n N) iter.Seq[N] {
	return slices.Values(g.adj[n])
}

// BFS yields every node reachable from start with its distance in edges,
// breadth first. Successors are visited in the order their edges were
// added. It yields nothing if start is not in the graph.
func (g *Graph[N]) BFS(start N) iter.Seq2[N, int] {
	return func(yield func(N, int) bool) {
		if _, ok := g.adj[start]; !ok {
			return
		}
		dist := map[N]int{start: 0}
		queue := []N{start}
		for len(queue) > 0 {
			n := queue[0]
			queue = queue[1:]
			if !yield(n, dist[n]) {
				return
			}
			for _, next := range g.adj[n] {
				if _, seen := dist[next]; !seen {
					dist[next] = dist[n] + 1
					queue = append(queue, next)
				}
			}
		}
	}
}

// TopoSort orders the nodes so every edge goes from an earlier node to a
// later one, choosing the smallest available node at each step. If the
// graph has a cycle it returns nil and a *CycleError.
func (g *Graph[N]) TopoSort() ([]N, error) {
	indegree := make(map[N]int, len(g.adj))
	for _, succs := range g.adj {
		for _, s := range succs {
			indegree[s]++
		}
	}
	var ready []N
	for n := range g.adj {
		if indegree[n] == 0 {
			ready = append(ready, n)
		}
	}
	order := make([]N, 0, len(g.adj))
	for len(ready) > 0 {
		i := slices.Index(ready, slices.Min(ready))
		n := ready[i]
		ready = slices.Delete(ready, i, i+1)
		order = append(order, n)
		for _, s := range g.adj[n] {
			indegree[s]--
			if indegree[s] == 0 {
				ready = append(ready, s)
			}
		}
	}
	if len(order) < len(g.adj) {
		var stuck []N
		for n := range g.adj {
			if indegree[n] > 0 {
				stuck = append(stuck, n)
			}
		}
		slices.Sort(stuck)
		return nil, &CycleError[N]{Nodes: stuck}
	}
	return order, nil
}
