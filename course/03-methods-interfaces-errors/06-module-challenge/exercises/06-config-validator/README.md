# Config Validator

Validate a configuration and report **every** problem at once: each problem is a typed `*FieldError` wrapping a sentinel, all of them combined with `errors.Join`. Then write `Fields`, which walks any error tree and collects the field names.

```go
var (
	ErrRequired   = errors.New("required")
	ErrOutOfRange = errors.New("out of range")
	ErrInvalid    = errors.New("invalid")
)

type FieldError struct {
	Field string
	Err   error // one of the sentinels above
}

func (e *FieldError) Error() string
func (e *FieldError) Unwrap() error

type Config struct {
	Name     string
	Port     int
	Env      string
	Replicas int
	Tags     []string
}

func (c Config) Validate() error
func Fields(err error) []string
```

## Examples

| Code | Result |
|---|---|
| `(&FieldError{"port", ErrOutOfRange}).Error()` | `"port: out of range"` |
| `Config{Name: "api", Port: 8080, Env: "dev", Replicas: 1}.Validate()` | `nil` |
| `Config{Port: 0, Env: "dev", Replicas: 1}.Validate().Error()` | `"name: required\nport: out of range"` |
| `errors.Is(err, ErrRequired)` | `true` — `Join` + `Unwrap` make every sentinel reachable |
| `errors.AsType[*FieldError](err)` | the **first** one: `{name required}`, `true` |
| `Fields(fmt.Errorf("load: %w", err))` | `[name port]` |

## Rules
- **FieldError**: `Error()` is `Field + ": " + Err.Error()`; `Unwrap()` returns `Err`.
- **Validate** checks the fields in this order and adds one `&FieldError{...}` per problem:

  | Field name | Problem |
  |---|---|
  | `name` | blank after `strings.TrimSpace` → `ErrRequired` |
  | `port` | not in `1..65535` → `ErrOutOfRange` |
  | `env` | `""` → `ErrRequired`; not exactly `dev`, `staging` or `prod` → `ErrInvalid` |
  | `replicas` | not in `1..10` → `ErrOutOfRange`; **otherwise**, if env is `prod` and replicas `< 2` → `ErrInvalid` |
  | `tags[i]` | blank tag → `ErrRequired`; otherwise the same text as an earlier tag → `ErrInvalid` |

- At most one problem per field (per tag). Return `errors.Join(problems...)` — which is `nil` when there are none. A valid config must return a true `nil`.
- **Fields** returns the `Field` of every `*FieldError` found in the tree, **depth-first in order**:
  - check the error itself with a type assertion, then descend
  - an error with `Unwrap() []error` (from `errors.Join` or `fmt.Errorf` with several `%w`) → visit each child in order
  - an error with `Unwrap() error` → visit the one child
  - `nil` or no field errors → `nil`
