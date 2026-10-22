// Package build plans the order in which targets are built.
package build

import (
	"errors"

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

// Plan returns the build order: every target after all of its dependencies,
// alphabetical where there is a choice.
func Plan(targets []Target) ([]string, error) {
	// TODO: build a graph.Graph[string] with an edge dep → target, then TopoSort.
	_ = graph.New[string]
	return nil, nil
}

// Affected returns, sorted, the targets to rebuild when changed changes:
// changed itself and everything that depends on it, directly or not.
func Affected(targets []Target, changed string) ([]string, error) {
	// TODO: the same graph, then BFS from changed.
	return nil, nil
}
