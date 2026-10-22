# Lazy Config

Load an app's configuration **once**, on first use, even when many goroutines ask for it at the same moment.

```go
type Config struct {
	Host string
	Port int
}

type App struct {
	config func() (Config, error) // set by New; runs the loader at most once
}

func New(load func() (Config, error)) *App

func (a *App) Config() (Config, error)
func (a *App) Addr() (string, error)
```

## Examples

| Code | Result |
|---|---|
| `a := New(load)` | `load` is **not** called yet |
| `a.Config()` | calls `load`, returns `{Host:localhost Port:8080}, nil` |
| `a.Config()` again | same result, `load` not called again |
| 100 goroutines call `a.Config()` at once | `load` runs exactly once; all 100 get its result |
| `a.Addr()` | `"localhost:8080", nil` |
| loader returns `errMissing` | every `Config()` returns `errMissing`; `load` still ran once |
| …then `a.Addr()` | `"", "appconf: config file missing"` |

## Rules
- `New` must not call `load`.
- `load` runs at most once per `App`, no matter how many goroutines call `Config` or `Addr`. Callers that arrive while it is running **wait** for it and get its result.
- An error is cached like a value: it is **not** retried.
- If `load` panics, every `Config` call panics with the same value (and `load` still ran once).
- Two `App`s each have their own loader and run it once each.
- `Addr` formats `"<host>:<port>"` (`"%s:%d"`). On a load error it returns `""` and `fmt.Errorf("appconf: %w", err)`.
