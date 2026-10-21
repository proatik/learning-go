# Stream Toolkit

Solve five small tasks by **composing** the `io` and `bufio` helpers instead of writing read loops.

```go
func Head(r io.Reader, n int64) (string, error)            // first n bytes      → io.LimitReader
func Join(sep string, parts ...io.Reader) io.Reader        // parts with sep between → io.MultiReader
func Broadcast(msg string, ws ...io.Writer) error          // msg to every writer → io.MultiWriter
func ReadAndLog(r io.Reader, log io.Writer) (string, error) // read all, copy to log → io.TeeReader
func WriteLines(w io.Writer, lines []string) error         // buffered output     → bufio.Writer
```

## Examples

| Call | Result |
|---|---|
| `Head(strings.NewReader("0123456789"), 4)` | `"0123", nil` |
| `Head(strings.NewReader("abc"), 10)` | `"abc", nil` |
| `Head(strings.NewReader("abc"), 0)` | `"", nil` |
| `io.ReadAll(Join(", ", strings.NewReader("a"), strings.NewReader("b"), strings.NewReader("c")))` | `"a, b, c"` |
| `io.ReadAll(Join("-"))` | `""` |
| `Broadcast("hi\n", &b1, &b2)` | `b1` and `b2` both hold `"hi\n"` |
| `ReadAndLog(strings.NewReader("data"), &log)` | `"data", nil`; `log` holds `"data"` |
| `WriteLines(&buf, []string{"x", "y"})` | `buf` holds `"x\ny\n"` |

## Rules
- `Head`: `n <= 0` returns `""`. Return any read error (`io.EOF` is not an error).
- `Join`: `sep` goes **between** parts only. With no parts the reader is empty. Nothing is read until the caller reads.
- `Broadcast`: write `msg` once to each writer; return the first error. With no writers it does nothing and returns `nil`.
- `ReadAndLog`: return everything read from `r` and the read error; every byte read must also be written to `log`.
- `WriteLines`: each line followed by `\n`. Use a `bufio.Writer` so a few hundred short lines reach `w` in **one or two** `Write` calls. Return the error from `Flush`.
