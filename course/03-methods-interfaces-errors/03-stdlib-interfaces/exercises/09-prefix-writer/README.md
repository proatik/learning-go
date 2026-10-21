# Prefix Writer and Tee Logger

Write a writer decorator that puts a prefix at the start of **every line**, even when lines arrive split across many `Write` calls. Then combine it with `io.MultiWriter` into a logger.

```go
type PrefixWriter struct {
	// unexported fields of your choice
}

func NewPrefixWriter(w io.Writer, prefix string) *PrefixWriter
func (pw *PrefixWriter) Write(p []byte) (int, error)

// NewLogger returns a writer that prefixes every line and sends the result to all of ws.
func NewLogger(prefix string, ws ...io.Writer) io.Writer
```

## Examples

With `pw := NewPrefixWriter(&buf, "> ")`:

| Writes (in order) | `buf` afterwards |
|---|---|
| `"hello\n"` | `"> hello\n"` |
| `"hel"`, `"lo\nwor"`, `"ld\n"` | `"> hello\n> world\n"` |
| `"a\nb"` | `"> a\n> b"` |
| `"a\n"`, then nothing more | `"> a\n"` (no dangling `"> "`) |
| `"\n\n"` | `"> \n> \n"` |
| `""` | `""` |

`l := NewLogger("[app] ", &console, &file); fmt.Fprintln(l, "started")` puts `"[app] started\n"` in **both** buffers.

## Rules
- A line starts at the very first byte ever written and right after each `\n`. Write the prefix **lazily**: only when the first byte of a new line arrives (so `"a\n"` never produces a trailing prefix).
- The result must not depend on how the bytes are split into `Write` calls: one byte at a time gives the same output as one big write.
- On success return `len(p), nil` (the caller's byte count, **not** the number of bytes including prefixes).
- If the wrapped writer fails, return `0` and its error.
- An empty prefix passes bytes through unchanged. Don't modify `p`.
