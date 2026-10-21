package graph

import (
	"errors"
	"fmt"
	"slices"
)

// Errors returned by ShortestPath.
var (
	ErrUnknownNode = errors.New("graph: unknown node")
	ErrNoPath      = errors.New("graph: no path")
)

// AddEdge adds the directed edge from -> to, ignoring duplicates, and
// makes sure both nodes are keys of g.
func AddEdge(g map[string][]string, from, to string) {
	if !slices.Contains(g[from], to) {
		g[from] = append(g[from], to)
	}
	if _, ok := g[to]; !ok {
		g[to] = []string{}
	}
}

// exists reports whether node is a key of g or appears in any edge list.
func exists(g map[string][]string, node string) bool {
	if _, ok := g[node]; ok {
		return true
	}
	for _, next := range g {
		if slices.Contains(next, node) {
			return true
		}
	}
	return false
}

// bfs explores g from start and returns each reached node's distance and
// predecessor. Neighbours are visited in slice order.
func bfs(g map[string][]string, start string) (dist map[string]int, prev map[string]string) {
	dist = map[string]int{start: 0}
	prev = map[string]string{}
	queue := []string{start}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		for _, n := range g[cur] { // a missing key yields a nil slice: no edges
			if _, seen := dist[n]; seen {
				continue
			}
			dist[n] = dist[cur] + 1
			prev[n] = cur // fixed on first discovery
			queue = append(queue, n)
		}
	}
	return dist, prev
}

// Distances returns the BFS distance from `from` to every reachable node.
func Distances(g map[string][]string, from string) map[string]int {
	if !exists(g, from) {
		return nil
	}
	dist, _ := bfs(g, from)
	return dist
}

// ShortestPath returns the nodes on a shortest path from `from` to `to`.
func ShortestPath(g map[string][]string, from, to string) ([]string, error) {
	for _, n := range []string{from, to} {
		if !exists(g, n) {
			return nil, fmt.Errorf("%w: %q", ErrUnknownNode, n)
		}
	}
	dist, prev := bfs(g, from)
	if _, ok := dist[to]; !ok {
		return nil, fmt.Errorf("%w from %q to %q", ErrNoPath, from, to)
	}

	// Walk the predecessors back from the target, then reverse.
	path := []string{to}
	for cur := to; cur != from; {
		cur = prev[cur]
		path = append(path, cur)
	}
	slices.Reverse(path)
	return path, nil
}
