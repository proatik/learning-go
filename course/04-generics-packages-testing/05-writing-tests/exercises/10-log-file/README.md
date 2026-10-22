# Files in a temp dir

Test a small append-only log file with `t.TempDir` and `t.Cleanup`, and catch every broken version of it.

```go
var ErrClosed = errors.New("log is closed")

type Log struct { /* an open *os.File */ }

// Open opens the log at path for appending, creating the file if needed.
func Open(path string) (*Log, error)

func (l *Log) Write(msg string) error // appends msg + "\n"; ErrClosed after Close
func (l *Log) Close() error           // closing twice is allowed

// ReadLines returns the lines of the log at path, without their newlines.
func ReadLines(path string) ([]string, error)
```

`Log` is already written (in the read-only checker). You write the tests in `logfile_test.go`.

## Examples

| Steps | `ReadLines(path)` |
|---|---|
| `Open`, no writes | `[]` (no lines), `nil` |
| `Write("a")`, `Write("b")` | `["a" "b"]`, `nil` |
| then `Close`, `Open` again, `Write("c")` | `["a" "b" "c"]`, `nil` |
| file does not exist | `nil`, an error where `errors.Is(err, fs.ErrNotExist)` |

## Rules

- **Every file lives in `t.TempDir()`.** The checker fails if your tests leave anything in the exercise folder.
- **Close every `Log` with `t.Cleanup`**, in a helper like `openLog(t, path)`. On Windows, an open file stops `TempDir` from being removed.
- After `Close`, `Write` returns `ErrClosed` (check with `errors.Is`).
- Reopening must **append**: earlier lines stay.
- Compare slices with `slices.Equal` and report `"ReadLines() = %q, want %q"`.
- **The checker** requires calls to `TempDir` and `Cleanup`, runs your tests on the correct `Log` (they must pass), then once per planted bug (they must fail every time).
