# Command Registry

A plugin-style registry: handlers of different types register under a name, and a text line is dispatched to the right one.

```go
type Handler interface {
	Handle(args []string) string
}

// HandlerFunc adapts an ordinary function to the Handler interface.
type HandlerFunc func(args []string) string

func (f HandlerFunc) Handle(args []string) string

type Echo struct{}               // joins its args with a space
type Counter struct{ /* n int */ } // *Counter counts its calls

func (Echo) Handle(args []string) string
func (c *Counter) Handle(args []string) string

func Upper(h Handler) Handler   // decorator: upper-cases h's output

type Registry struct{ /* unexported fields */ } // zero value ready to use

func (r *Registry) Register(name string, h Handler) bool
func (r *Registry) Names() []string
func (r *Registry) Dispatch(line string) string
```

## Examples

| Code | Result |
|---|---|
| `Echo{}.Handle([]string{"hi", "there"})` | `"hi there"` |
| `c := &Counter{}; c.Handle(nil); c.Handle(nil)` | `"count=1"`, then `"count=2"` |
| `HandlerFunc(f).Handle(args)` | `f(args)` |
| `Upper(Echo{}).Handle([]string{"go"})` | `"GO"` |
| `r.Register("echo", Echo{})` / again | `true` / `false` (duplicate) |
| `r.Dispatch("  echo   a  b ")` | `"a b"` |
| `r.Dispatch("nope x")` | `"unknown command: nope"` |
| `r.Dispatch("   ")` | `""` |
| `r.Dispatch("help")` (no `help` registered) | `"commands: count echo"` |

## Rules
- `Register` returns `false` and stores nothing when:
  - `name` is empty or contains any whitespace (`strings.ContainsFunc(name, unicode.IsSpace)`),
  - `name` is already registered,
  - `h` is a nil interface, a nil `HandlerFunc`, or a nil `*Counter`.
- `Names` returns registered names sorted ascending (none → empty, `nil` is fine).
- `Dispatch` splits `line` with `strings.Fields`:
  - no fields → `""`
  - first field is a registered name → that handler's `Handle(rest)` where `rest` is the remaining fields (an empty slice when there are none)
  - first field is `help` and `help` is **not** registered → `"commands: "` + names joined by a space, or `"no commands"` if none are registered
  - otherwise → `"unknown command: <name>"`
- Names are case-sensitive.
- `Counter`'s result is `"count=<number of calls so far>"`. `Echo` with no args returns `""`.
