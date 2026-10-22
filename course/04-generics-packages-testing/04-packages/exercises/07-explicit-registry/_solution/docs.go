// Package docs wires the render registry used by the documentation tool.
package docs

import (
	"learngo/course/04-generics-packages-testing/04-packages/exercises/07-explicit-registry/builtin"
	"learngo/course/04-generics-packages-testing/04-packages/exercises/07-explicit-registry/render"
)

// Default returns a new registry holding the built-in formats.
func Default() (*render.Registry, error) {
	r := &render.Registry{}
	if err := builtin.RegisterAll(r); err != nil {
		return nil, err
	}
	return r, nil
}
