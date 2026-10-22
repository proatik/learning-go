// Package render renders a titled list of items in named formats.
package render

import (
	"errors"
	"fmt"
	"maps"
	"slices"
)

// Func renders a title and its items.
type Func func(title string, items []string) string

// Errors returned by Registry methods.
var (
	ErrEmptyName = errors.New("render: empty name")
	ErrNilFunc   = errors.New("render: nil func")
	ErrDuplicate = errors.New("render: already registered")
	ErrUnknown   = errors.New("render: unknown format")
)

// Registry maps format names to render functions. The zero value is an
// empty registry ready to use.
type Registry struct {
	funcs map[string]Func
}

// Register adds f under name. It fails with ErrEmptyName, ErrNilFunc, or an
// error wrapping ErrDuplicate when name is already registered.
func (r *Registry) Register(name string, f Func) error {
	if name == "" {
		return ErrEmptyName
	}
	if f == nil {
		return ErrNilFunc
	}
	if _, ok := r.funcs[name]; ok {
		return fmt.Errorf("%w: %q", ErrDuplicate, name)
	}
	if r.funcs == nil {
		r.funcs = map[string]Func{}
	}
	r.funcs[name] = f
	return nil
}

// Render renders title and items with the format registered as name, or
// returns an error wrapping ErrUnknown.
func (r *Registry) Render(name, title string, items []string) (string, error) {
	f, ok := r.funcs[name]
	if !ok {
		return "", fmt.Errorf("%w: %q", ErrUnknown, name)
	}
	return f(title, items), nil
}

// Names returns the registered format names in sorted order.
func (r *Registry) Names() []string {
	return slices.Sorted(maps.Keys(r.funcs))
}
