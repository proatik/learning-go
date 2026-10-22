# ID Generator

`ids.Next()` uses one package-level counter shared by the **whole program**, so two directories can never both start at `u-1`, and tests leak into each other. Replace it with a `Generator` type. Keep `ids.Next` working for old callers, but deprecate it.

```go
// package ids  (ids/ids.go)
type Generator struct{ /* unexported fields */ } // the zero value is ready: IDs "1", "2", …

func New(prefix string) *Generator
func (g *Generator) Next() string
func (g *Generator) Issued() int

// Deprecated: … (old API, kept for compatibility)
func Next() string

// package users  (users.go)
var ErrEmptyName = errors.New("users: empty name")

type Directory struct{ /* unexported fields */ }

func NewDirectory(prefix string) *Directory
func (d *Directory) Add(name string) (string, error)
func (d *Directory) Name(id string) (string, bool)
func (d *Directory) Len() int
```

## Examples

| Code | Result |
|---|---|
| `g := ids.New("u-"); g.Next(); g.Next()` | `"u-1"`, `"u-2"` |
| `g.Issued()` (after the two calls) | `2` |
| `var g ids.Generator; g.Next()` | `"1"` |
| `a, b := ids.New("x"), ids.New("x"); a.Next(); b.Next()` | `"x1"`, `"x1"` (independent) |
| `d := users.NewDirectory("u-"); d.Add("Ann")` | `"u-1", nil` |
| `d.Name("u-1")` | `"Ann", true` |
| a second `users.NewDirectory("u-")` | its first `Add` also returns `"u-1"` |
| `d.Add("")` | `"", ErrEmptyName` (no ID is used up) |

## Rules
- A `Generator` counts from 1: the n-th `Next()` returns `prefix + n`. `Issued` returns how many IDs it has produced.
- **Remove the package-level `counter` variable.** The old `ids.Next()` must call `Next` on a package-level **default `Generator`** with an empty prefix. The standard library keeps old package-level APIs alive the same way (e.g. `http.DefaultClient`).
- Give `ids.Next` a doc comment with a paragraph that starts with `Deprecated:` and mentions `Generator`. A test reads the comment.
- `users.Directory` owns its own `*ids.Generator` (created in `NewDirectory`) and a `map[string]string` from ID to name. It must **not** call the deprecated `ids.Next`.
- `Add` checks the name **before** taking an ID, so a rejected name doesn't use up a number.
