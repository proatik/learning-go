// Package render renders a titled list of items in named formats.
package render

import "errors"

// Func renders a title and its items.
type Func func(title string, items []string) string

// Errors returned by Registry methods.
var (
	ErrEmptyName = errors.New("render: empty name")
	ErrNilFunc   = errors.New("render: nil func")
	ErrDuplicate = errors.New("render: already registered")
	ErrUnknown   = errors.New("render: unknown format")
)

// ---- Old API: delete everything in this block once Registry works. ----

// registry is filled by other packages' init functions. Every importer shares it.
var registry = map[string]Func{}

// Register adds f under name to the shared registry, silently replacing
// any earlier entry.
func Register(name string, f Func) { registry[name] = f }

// Render renders with the shared registry.
func Render(name, title string, items []string) (string, error) {
	f, ok := registry[name]
	if !ok {
		return "", ErrUnknown
	}
	return f(title, items), nil
}

// ---- New API ----

// Registry maps format names to render functions. The zero value is an
// empty registry ready to use.
type Registry struct {
	funcs map[string]Func
}

// Register adds f under name. It fails with ErrEmptyName, ErrNilFunc, or an
// error wrapping ErrDuplicate when name is already registered.
func (r *Registry) Register(name string, f Func) error {
	// TODO
	return nil
}

// Render renders title and items with the format registered as name, or
// returns an error wrapping ErrUnknown.
func (r *Registry) Render(name, title string, items []string) (string, error) {
	// TODO
	return "", nil
}

// Names returns the registered format names in sorted order.
func (r *Registry) Names() []string {
	// TODO
	return nil
}
