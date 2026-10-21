# Text Stats with bufio.Scanner

Count lines and words of **any** `io.Reader` and find its longest line, using `bufio.Scanner`.

```go
func LineCount(r io.Reader) (int, error)
func WordCount(r io.Reader) (int, error)
func Longest(r io.Reader) (string, error)
```

## Examples

| Input (via `strings.NewReader`) | `LineCount` | `WordCount` | `Longest` |
|---|---|---|---|
| `""` | `0` | `0` | `""` |
| `"one two\nthree"` | `2` | `3` | `"one two"` |
| `"a\n"` | `1` | `1` | `"a"` |
| `"\n\n"` | `2` | `0` | `""` |
| `"  hello   world\n\tgo "` | `2` | `3` | `"  hello   world"` |
| `"hi\r\nthere\r\n"` | `2` | `2` | `"there"` |
| `"日本語\nabcd"` | `2` | `2` | `"abcd"` |

## Rules
- Lines follow `bufio.ScanLines`: a final line without `\n` counts; a trailing `\r` before `\n` is removed; empty lines count.
- Words follow `bufio.ScanWords`: separated by any Unicode white space.
- `Longest` compares lengths in **runes** (`"日本語"` has 3); on a tie, the **first** longest line wins. The returned line has no `\n` / `\r`.
- If the reader fails, return the error from `Scanner.Err()` (the count/line returned with it doesn't matter). Never return `io.EOF`: a clean end is `nil`.
