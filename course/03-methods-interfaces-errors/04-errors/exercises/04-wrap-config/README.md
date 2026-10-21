# Wrap with Context

Load a `host:port` config file. Every error must say **which file** failed and still let callers test the cause with `errors.Is` / `errors.As`.

```go
type Config struct {
	Host string
	Port int
}

var (
	ErrNotFound    = errors.New("file not found")
	ErrMissingHost = errors.New("missing host")
	ErrMissingPort = errors.New("missing port")
	ErrPortRange   = errors.New("port out of range")
)

func parse(content string) (Config, error)                         // helper: no file name in its errors
func LoadConfig(files map[string]string, name string) (Config, error) // adds `config "<name>": `
```

`files` stands in for a disk: it maps a file name to its content.

## Examples

With `name = "app.conf"`:

| `files["app.conf"]` | Result / `err.Error()` | Matches |
|---|---|---|
| `"localhost:8080"` | `Config{"localhost", 8080}, nil` | |
| `"  db.local:5432\n"` | `Config{"db.local", 5432}, nil` | |
| *(no such key)* | `config "app.conf": file not found` | `ErrNotFound` |
| `"localhost"` | `config "app.conf": missing port` | `ErrMissingPort` |
| `":80"` | `config "app.conf": missing host` | `ErrMissingHost` |
| `"h:abc"` | `config "app.conf": port "abc": strconv.Atoi: parsing "abc": invalid syntax` | `*strconv.NumError`, `strconv.ErrSyntax` |
| `"h:70000"` | `config "app.conf": port 70000: port out of range` | `ErrPortRange` |

## Rules
- `parse`:
  - Trim `content` with `strings.TrimSpace`.
  - Split it at the **first** `:`. No `:` → `ErrMissingPort`. Empty host → `ErrMissingHost`. Return both sentinels unwrapped.
  - Convert the port with `strconv.Atoi`. On failure, return `fmt.Errorf("port %q: %w", p, err)`, which keeps the `*strconv.NumError` in the chain.
  - A valid port is `1` to `65535`. Otherwise return `fmt.Errorf("port %d: %w", n, ErrPortRange)`.
- `LoadConfig`:
  - A missing file → `fmt.Errorf("config %q: %w", name, ErrNotFound)`.
  - Any `parse` error → `fmt.Errorf("config %q: %w", name, err)`. Wrap it **once**, here.
- Return a zero `Config` with every error.
