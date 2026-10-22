// Package build plans the order in which targets are built.
package build

import (
	"errors"
	"fmt"
	"slices"

	"learngo/course/04-generics-packages-testing/06-module-challenge/exercises/07-build-graph/graph"
)

// ErrUnknownDep is wrapped when a target depends on a name that is not a
// target.
var ErrUnknownDep = errors.New("build: unknown dependency")

// Target is something to build and the targets it needs first.
type Target struct {
	Name string
	Deps []string
}

// depGraph returns a graph with an edge dep → target for every dependency.
func depGraph(targets []Target) (*graph.Graph[string], error) {
	known := make(map[string]bool, len(targets))
	for _, t := range targets {
		known[t.Name] = true
	}
	g := graph.New[string]()
	for _, t := range targets {
		g.AddNode(t.Name)
		for _, d := range t.Deps {
			if !known[d] {
				return nil, fmt.Errorf("%w: %q needs %q", ErrUnknownDep, t.Name, d)
			}
			g.AddEdge(d, t.Name)
		}
	}
	return g, nil
}

// Plan returns the build order: every target after all of its dependencies,
// alphabetical where there is a choice.
func Plan(targets []Target) ([]string, error) {
	g, err := depGraph(targets)
	if err != nil {
		return nil, err
	}
	order, err := g.TopoSort()
	if err != nil {
		return nil, fmt.Errorf("build: %w", err)
	}
	return order, nil
}

// Affected returns, sorted, the targets to rebuild when changed changes:
// changed itself and everything that depends on it, directly or not.
func Affected(targets []Target, changed string) ([]string, error) {
	g, err := depGraph(targets)
	if err != nil {
		return nil, err
	}
	out := []string{}
	for n := range g.BFS(changed) {
		out = append(out, n)
	}
	slices.Sort(out)
	return out, nil
}
