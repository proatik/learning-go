# Line Reader

Wrap a `bufio.Scanner` in a **single-use** iterator that yields numbered lines from any `io.Reader`, skipping blanks and comments, and reports errors through an `Err` method.

```go
type LineReader struct {
	sc   *bufio.Scanner
	line int // number of the last line read from sc
}

func NewLineReader(r io.Reader) *LineReader // provided
func (lr *LineReader) Err() error           // provided: the scanner's error, nil at a clean EOF

// All returns a single-use iterator over (line number, trimmed line) pairs.
func (lr *LineReader) All() iter.Seq2[int, string]
```

## Examples

Input `"host = a\n\n  # comment\n  port = 80  \n"`:

| Code | Pairs |
|---|---|
| `for n, l := range NewLineReader(r).All()` | `(1, "host = a")`, `(4, "port = 80")` |

| Situation | Result |
|---|---|
| Empty input | No pairs; `Err() == nil` |
| Input `"a\r\nb"` (Windows line ending, no final `\n`) | `(1, "a")`, `(2, "b")` |
| `break` after `(1, …)`, then range over `All()` again | Continues with the next line: `(4, "port = 80")` |
| The reader fails after `"x\n"` | `(1, "x")`; then `Err()` returns the reader's error |

## Rules
- Line numbers start at `1` and count **every** line of the input, including skipped ones.
- Each line is trimmed with `strings.TrimSpace`. Skip a line if it is then empty or starts with `#`.
- **Single-use:** the scanner can't rewind. Ranging over `All()` again continues after the last line read, with correct numbers (keep the count in `lr.line`, not in a local variable).
- Read **lazily**: one `Scan` per input line, only while the loop wants more. After a `break` no more input is read.
- Stop as soon as `yield` returns `false`.
