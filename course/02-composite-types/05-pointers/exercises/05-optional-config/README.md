# Optional Config Fields

A `Config` uses **pointer fields** so "not set" (`nil`) is different from "set to the zero value" (`new(0)`, `new(false)`, `new("")`).

```go
type Config struct {
	Retries *int    // default 3
	Cache   *bool   // default true
	Name    *string // default "app"
}

type Settings struct {
	Retries int
	Cache   bool
	Name    string
}

func Effective(c Config) Settings
func Merge(base, override Config) Config
func FromSettings(s Settings) Config
```

| Function | Behavior |
|---|---|
| `Effective(c)` | Each nil field → its default; each set field → its value |
| `Merge(base, override)` | Each field: `override`'s if non-nil, otherwise `base`'s (which may be nil) |
| `FromSettings(s)` | A `Config` with **every** field set to the value from `s` |

| Call | Result |
|---|---|
| `Effective(Config{})` | `{3 true app}` |
| `Effective(Config{Retries: new(0), Cache: new(false), Name: new("")})` | `{0 false }` |
| `Merge(Config{Retries: new(5)}, Config{Name: new("api")})` → `Effective` | `{5 true api}` |
| `FromSettings(Settings{0, false, "x"})` | `Retries`, `Cache`, `Name` all non-nil: `0`, `false`, `"x"` |

## Rules
- Build pointers with Go 1.26+ `new(expr)`: `Retries: new(s.Retries)`. No helper functions, no temporary variables.
- `Merge` must not modify `base` or `override`.
- Every `FromSettings` call returns fresh pointers: two Configs never share a variable.
