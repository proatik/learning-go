// Package graph provides a generic directed graph with BFS and topological
// sorting.
package graph

import (
	"cmp"
	"errors"
	"iter"
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
	// TODO
	return ""
}

// Is makes errors.Is(err, ErrCycle) true for any *CycleError.
func (e *CycleError[N]) Is(target error) bool {
	// TODO
	return false
}

// Graph is a directed graph. Create one with New.
type Graph[N cmp.Ordered] struct {
	adj map[N][]N // node → successors in insertion order; every node is a key
}

// New returns an empty graph.
func New[N cmp.Ordered]() *Graph[N] {
	// TODO
	return nil
}

// AddNode adds n if it is not already in the graph.
func (g *Graph[N]) AddNode(n N) {
	// TODO
}

// AddEdge adds the edge from → to, adding both nodes. A repeated edge is
// ignored.
func (g *Graph[N]) AddEdge(from, to N) {
	// TODO
}

// Nodes yields every node in sorted order.
func (g *Graph[N]) Nodes() iter.Seq[N] {
	return func(yield func(N) bool) {
		// TODO
	}
}

// Neighbors yields the successors of n in the order their edges were added.
func (g *Graph[N]) Neighbors(n N) iter.Seq[N] {
	return func(yield func(N) bool) {
		// TODO
	}
}

// BFS yields every node reachable from start with its distance in edges,
// breadth first. Successors are visited in the order their edges were
// added. It yields nothing if start is not in the graph.
func (g *Graph[N]) BFS(start N) iter.Seq2[N, int] {
	return func(yield func(N, int) bool) {
		// TODO
	}
}

// TopoSort orders the nodes so every edge goes from an earlier node to a
// later one, choosing the smallest available node at each step. If the
// graph has a cycle it returns nil and a *CycleError.
func (g *Graph[N]) TopoSort() ([]N, error) {
	// TODO
	return nil, nil
}
