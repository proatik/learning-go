# Mini Logging Library

Design a small logging package around one tiny interface. Concrete loggers write to any `io.Writer` or to memory; **decorators** filter by level or fan out to several loggers — and because they all are `Logger`s, they nest freely.

```go
var ErrUnknownLevel = errors.New("logx: unknown level")

type Level int

const (
	Debug Level = iota
	Info
	Warn
	Error
)

func (l Level) String() string
func ParseLevel(s string) (Level, error)

type Logger interface {
	Log(level Level, msg string) error
}

type TextLogger struct { w io.Writer; prefix string }
func NewTextLogger(w io.Writer, prefix string) *TextLogger
func (l *TextLogger) Log(level Level, msg string) error
func (l *TextLogger) With(prefix string) *TextLogger

type MultiLogger []Logger
func (m MultiLogger) Log(level Level, msg string) error

type LevelFilter struct {
	Min  Level
	Next Logger
}
func (f LevelFilter) Log(level Level, msg string) error

type Entry struct {
	Level Level
	Msg   string
}
func (e Entry) String() string

type MemoryLogger struct{ Entries []Entry }
func (m *MemoryLogger) Log(level Level, msg string) error

func Logf(l Logger, level Level, format string, args ...any) error
```

## Examples

| Code | Result |
|---|---|
| `Warn.String()`, `Level(9).String()` | `"WARN"`, `"Level(9)"` |
| `ParseLevel(" info ")` | `Info, nil` |
| `ParseLevel("fatal")` | `logx: unknown level: "fatal"` |
| `NewTextLogger(w, "app: ").Log(Info, "started")` | writes `"app: INFO  started\n"` |
| `l.With("db: ").Log(Error, "down")` | writes `"app: db: ERROR down\n"` |
| `LevelFilter{Min: Warn, Next: l}.Log(Info, "x")` | `nil`, nothing written |
| `Entry{Warn, "disk"}.String()` | `"WARN disk"` |

## Rules
- **Level.String**: `DEBUG`, `INFO`, `WARN`, `ERROR`; any other value → `"Level(N)"`.
- **ParseLevel** ignores case and surrounding spaces; anything else → `fmt.Errorf("%w: %q", ErrUnknownLevel, s)` with the **original** `s`.
- **TextLogger.Log** writes exactly `fmt.Fprintf(w, "%s%-5s %s\n", prefix, level, msg)` (the level is left-aligned in 5 columns). A write error is returned as `fmt.Errorf("logx: %w", err)`.
- **With** returns a **new** logger on the same writer with prefix `parent prefix + prefix`; the parent is unchanged.
- **MultiLogger.Log** calls every non-nil logger in order — even after one fails — and returns `errors.Join` of their errors (`nil` when all succeed, or when empty).
- **LevelFilter.Log** forwards to `Next` only if `level >= Min` and `Next != nil`, returning `Next`'s error unchanged; otherwise returns `nil`.
- **MemoryLogger.Log** appends an `Entry` and returns `nil` (pointer receiver — store `&MemoryLogger{}` in a `Logger`).
- **Logf** logs `fmt.Sprintf(format, args...)`; a `nil` logger returns `nil`.
