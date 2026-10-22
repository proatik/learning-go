# Functional Options

Give `server.New` sensible defaults and optional, named settings, then use them from an `api` package.

```go
// package server  (server/server.go)
const (
	DefaultName     = "server"
	DefaultTimeout  = 30 * time.Second
	DefaultMaxConns = 100
)

var (
	ErrNoAddr   = errors.New("server: address is empty")
	ErrTimeout  = errors.New("server: timeout must be positive")
	ErrMaxConns = errors.New("server: maxConns must be at least 1")
)

type Server struct{ /* unexported fields */ }
type Option func(*Server)

func WithName(name string) Option
func WithTimeout(d time.Duration) Option
func WithMaxConns(n int) Option
func WithMiddleware(names ...string) Option

func New(addr string, opts ...Option) (*Server, error)

func (s *Server) Addr() string
func (s *Server) Name() string
func (s *Server) Timeout() time.Duration
func (s *Server) MaxConns() int
func (s *Server) Middleware() []string
func (s *Server) String() string

// package api  (api.go)
func NewProduction(addr string, extra ...server.Option) (*server.Server, error)
```

## Examples

| Call | `String()` of the result |
|---|---|
| `server.New(":8080")` | `"server@:8080 timeout=30s maxConns=100 middleware=[]"` |
| `server.New(":80", WithName("web"), WithTimeout(5*time.Second))` | `"web@:80 timeout=5s maxConns=100 middleware=[]"` |
| `server.New(":80", WithMiddleware("log"), WithMiddleware("auth", "gzip"))` | `"server@:80 timeout=30s maxConns=100 middleware=[log auth gzip]"` |
| `api.NewProduction(":443")` | `"api@:443 timeout=10s maxConns=1000 middleware=[recover log]"` |
| `api.NewProduction(":443", WithMaxConns(5))` | `"api@:443 timeout=10s maxConns=5 middleware=[recover log]"` |

| Call | Error |
|---|---|
| `server.New("")` | `ErrNoAddr` |
| `server.New("", WithTimeout(0), WithMaxConns(0))` | wraps `ErrNoAddr`, `ErrTimeout` **and** `ErrMaxConns` |

## Rules
- `New` starts from the defaults, applies the options **in order** (a later option overrides an earlier one), then validates.
- `WithMiddleware` **appends**; calling it several times accumulates the names.
- Validation reports **every** problem with `errors.Join`, in the order address → timeout → max conns:
  - `addr == ""` → `ErrNoAddr`;
  - `timeout <= 0` → `ErrTimeout`;
  - `maxConns < 1` → `ErrMaxConns`.
  - On error, `New` returns `nil`.
- `Middleware` returns a **copy** (empty or nil when there is none).
- `String` is `"<name>@<addr> timeout=<timeout> maxConns=<n> middleware=<list>"`. Use `%v` for the duration (`5s`, `1m30s`) and for the `[]string` (`[a b]`).
- `NewProduction` applies `WithName("api")`, `WithTimeout(10*time.Second)`, `WithMaxConns(1000)`, `WithMiddleware("recover", "log")`, **then** `extra...`, so callers can override.
