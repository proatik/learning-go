# Explicit Registry

Today `builtin` fills a **package-level** map in `render` from an `init()` function. Replace that with a `Registry` value that callers create and fill **explicitly**.

```go
// package render  (render/render.go)
type Func func(title string, items []string) string

var (
	ErrEmptyName = errors.New("render: empty name")
	ErrNilFunc   = errors.New("render: nil func")
	ErrDuplicate = errors.New("render: already registered")
	ErrUnknown   = errors.New("render: unknown format")
)

type Registry struct{ /* unexported fields */ } // the zero value is an empty registry

func (r *Registry) Register(name string, f Func) error
func (r *Registry) Render(name, title string, items []string) (string, error)
func (r *Registry) Names() []string

// package builtin  (builtin/builtin.go)
func Plain(title string, items []string) string
func Markdown(title string, items []string) string
func Numbered(title string, items []string) string
func RegisterAll(r *render.Registry) error

// package docs  (docs.go)
func Default() (*render.Registry, error)
```

## Examples

| Call | Result |
|---|---|
| `Plain("Todo", []string{"a", "b"})` | `"Todo: a, b"` |
| `Plain("Todo", nil)` | `"Todo: (none)"` |
| `Markdown("Todo", []string{"a", "b"})` | `"# Todo\n\n- a\n- b"` |
| `Markdown("Todo", nil)` | `"# Todo\n\n_No items._"` |
| `Numbered("Todo", []string{"a", "b"})` | `"Todo\n1. a\n2. b"` |
| `Numbered("Todo", nil)` | `"Todo"` |
| `var r render.Registry; r.Register("p", builtin.Plain)` | `nil` |
| `r.Register("p", builtin.Plain)` again | error wrapping `ErrDuplicate`: `render: already registered: "p"` |
| `r.Render("nope", "T", nil)` | error wrapping `ErrUnknown`: `render: unknown format: "nope"` |
| `docs.Default()` then `.Names()` | `[markdown numbered plain]` |

## Rules
- `Register`, checked in order:
  - empty `name` → `ErrEmptyName`;
  - `f == nil` → `ErrNilFunc`;
  - name already used → `fmt.Errorf("%w: %q", ErrDuplicate, name)`.
  - Nothing is registered on error.
- `Render` calls the function registered under `name`. An unknown name → `""` and `fmt.Errorf("%w: %q", ErrUnknown, name)`.
- `Names` returns the registered names **sorted**: empty (or nil) for an empty registry.
- The zero value `render.Registry{}` works: create the map lazily in `Register`.
- `RegisterAll` registers `"plain"`, `"markdown"` and `"numbered"`, in that order, and returns the first error.
- `Default` returns a **new** registry with the built-ins on every call; registries never share state.
- **Delete the old API:**
  - the package-level `registry` map and the package-level `Register` / `Render` functions in `render.go`;
  - the `init` function in `builtin.go`.
  - A test checks the source files.
