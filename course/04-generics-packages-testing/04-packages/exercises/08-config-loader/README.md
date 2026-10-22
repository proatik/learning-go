# Config Loader

Write a `config` package that starts from defaults, applies settings from a `map[string]string` (as if read from a file or environment), and reports **every** problem at once.

```go
// package config  (config/config.go)
type Config struct {
	Host    string
	Port    int
	Debug   bool
	Timeout time.Duration
}

var (
	ErrUnknownKey = errors.New("config: unknown key")
	ErrBadValue   = errors.New("config: bad value")
	ErrInvalid    = errors.New("config: invalid")
)

func Default() Config
func (c Config) Validate() error
func Load(settings map[string]string) (Config, error)

// package daemon  (daemon.go)
func Describe(settings map[string]string) (string, error)
```

## Examples

| Call | Result |
|---|---|
| `Default()` | `Config{Host: "localhost", Port: 8080, Debug: false, Timeout: 5 * time.Second}` |
| `Load(nil)` | `Default(), nil` |
| `Load(map[string]string{"port": "9090", "debug": "true"})` | defaults with `Port: 9090, Debug: true` |
| `Load(map[string]string{"colour": "red"})` | error wrapping `ErrUnknownKey` |
| `Load(map[string]string{"port": "abc"})` | error wrapping `ErrBadValue` |
| `Load(map[string]string{"port": "70000"})` | error wrapping `ErrInvalid` |
| `Config{}.Validate()` | wraps `ErrInvalid`, with three problems |
| `daemon.Describe(map[string]string{"host": "example.com"})` | `"example.com:8080 debug=false timeout=5s"` |

`Load(map[string]string{"zzz": "1", "port": "x", "host": ""})` returns this error (keys in sorted order, then validation):

```text
config: bad value: port="x"
config: unknown key: "zzz"
config: invalid: host is empty
```

## Rules
- Keys:
  - `"host"`: used as is;
  - `"port"`: `strconv.Atoi`;
  - `"debug"`: `strconv.ParseBool`;
  - `"timeout"`: `time.ParseDuration`, which parses strings like `"5s"`, `"1m30s"` or `"250ms"`.
- `Load` starts from `Default()` and visits the keys in **sorted order** (`slices.Sorted(maps.Keys(settings))`), so errors come out in a stable order.
  - An unknown key → `fmt.Errorf("%w: %q", ErrUnknownKey, key)`.
  - A value that fails to parse → `fmt.Errorf("%w: %s=%q", ErrBadValue, key, value)`. The field keeps its previous value.
- After the keys, `Load` always runs `Validate` on the result and adds its error too.
- Combine everything with `errors.Join`. On any error, return `Config{}` and the joined error.
- `Validate` checks, in this order, and joins the problems (`nil` when valid):
  - empty `Host` → `fmt.Errorf("%w: host is empty", ErrInvalid)`;
  - `Port` outside `1..65535` → `fmt.Errorf("%w: port %d out of range", ErrInvalid, c.Port)`;
  - `Timeout <= 0` → `fmt.Errorf("%w: timeout must be positive", ErrInvalid)`.
- `Describe` loads the settings and returns `"<host>:<port> debug=<bool> timeout=<duration>"`, or `""` and the error.
