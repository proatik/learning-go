# Text Pipeline

Build streaming text filters as **`io.Reader` decorators** — each wraps another reader and transforms bytes as they flow through — plus an `io.Writer` that counts what passes, and a `Run` function that chains everything with `io.Copy`.

```go
type Stage func(io.Reader) io.Reader

func Upper(r io.Reader) io.Reader       // a–z → A–Z
func Rot13(r io.Reader) io.Reader       // rotate letters by 13
func StripDigits(r io.Reader) io.Reader // drop 0–9

type CountingWriter struct {
	W     io.Writer
	Bytes int64
	Lines int
}

func (c *CountingWriter) Write(p []byte) (int, error)

func Run(dst io.Writer, src io.Reader, stages ...Stage) (int64, error)
```

## Examples

| Code | Result |
|---|---|
| `io.ReadAll(Upper(strings.NewReader("héllo")))` | `"HéLLO"` |
| `Rot13` of `"Hello, World!"` | `"Uryyb, Jbeyq!"` |
| `StripDigits` of `"a1b22c333"` | `"abc"` |
| `Run(&out, strings.NewReader("hello 2026"), StripDigits, Upper)` | out `"HELLO "`, `6, nil` |
| `cw := &CountingWriter{W: &buf}; cw.Write([]byte("a\nb\n"))` | `4, nil`; `cw.Bytes == 4`, `cw.Lines == 2` |

## Rules
- Each stage returns a value of **your own reader type** (a struct holding the source `io.Reader`) with a `Read` method. Don't read the whole input up front — transform it chunk by chunk.
- Work on **bytes**, ASCII only: `Upper` changes only `a`–`z`; `Rot13` rotates only `a`–`z` and `A`–`Z` (case kept); `StripDigits` removes only `'0'`–`'9'`. Every other byte (including UTF-8 multi-byte characters) passes through unchanged.
- Follow the `io.Reader` contract (the tests run `iotest.TestReader`):
  - transform the `n` bytes read **before** returning, and return the source's error unchanged — including data that arrives together with an error
  - `StripDigits` must **never** return `0, nil` for a non-empty `p`: if a whole chunk was digits, read again
- **CountingWriter.Write** forwards `p` to `W` and returns `W`'s `n, err`. It adds only the bytes **actually written** (`p[:n]`) to `Bytes`, and the `'\n'` among them to `Lines`. If `W` is `nil`, it accepts everything: `len(p), nil`.
- **Run** applies `stages` left to right (`stages[0]` wraps `src` first), skips `nil` stages, and copies the result to `dst` with `io.Copy`, returning its result.
