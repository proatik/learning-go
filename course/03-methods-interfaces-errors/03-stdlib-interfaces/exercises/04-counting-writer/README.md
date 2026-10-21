# Counting Writer

Build a **decorator**: an `io.Writer` that wraps another writer and counts the bytes that pass through it.

```go
type CountingWriter struct {
	W io.Writer // the wrapped writer; nil means "count only, discard the bytes"
	N int64     // total bytes successfully written so far
}

func (c *CountingWriter) Write(p []byte) (int, error)

// Measure wraps w in a CountingWriter, passes it to write, and returns how
// many bytes write produced together with write's error.
func Measure(w io.Writer, write func(io.Writer) error) (int64, error)
```

## Examples

| Code | Result |
|---|---|
| `c := &CountingWriter{W: &buf}; fmt.Fprintf(c, "hi %d", 42)` | `buf` holds `"hi 42"`, `c.N == 5` |
| two writes `"abc"` then `"de"` | `c.N == 5` |
| `c := &CountingWriter{}; io.WriteString(c, "héllo")` | `c.N == 6` (bytes, not runes) |
| `Measure(&buf, func(w io.Writer) error { _, err := io.WriteString(w, "x"); return err })` | `1, nil` |

## Rules
- `Write` forwards `p` to `W` and adds the `n` **returned by `W`** to `N` (even when `W` also returns an error), then returns that same `n, err`.
- If `W` is `nil`, `Write` accepts everything: `N += len(p)`, returns `len(p), nil`.
- `Measure` returns the count **and** whatever error `write` returned (bytes written before the error still count).
- `Write` has a **pointer receiver** (it updates `N`), so callers pass `&CountingWriter{...}`.
