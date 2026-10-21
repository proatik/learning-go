package graph

import "errors"

// Errors returned by ShortestPath.
var (
	ErrUnknownNode = errors.New("graph: unknown node")
	ErrNoPath      = errors.New("graph: no path")
)

// AddEdge adds the directed edge from -> to, ignoring duplicates, and
// makes sure both nodes are keys of g.
func AddEdge(g map[string][]string, from, to string) {
	// TODO
}

// Distances returns the BFS distance from `from` to every reachable node.
func Distances(g map[string][]string, from string) map[string]int {
	// TODO
	return nil
}

// ShortestPath returns the nodes on a shortest path from `from` to `to`.
func ShortestPath(g map[string][]string, from, to string) ([]string, error) {
	// TODO
	return nil, nil
}
