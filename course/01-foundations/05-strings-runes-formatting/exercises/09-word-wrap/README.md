# Word Wrap

Implement:

```go
func Wrap(text string, width int) string
```

Re-flow `text` so that no line is longer than `width` **runes**, breaking only between words. Lines are joined with `"\n"`.

```go
Wrap("The quick brown fox jumps over the lazy dog", 10)
```
```
The quick
brown fox
jumps over
the lazy
dog
```

| Call | Result |
|---|---|
| `Wrap("a b c d", 3)` | `"a b\nc d"` |
| `Wrap("  lots   of\n\tspace  ", 20)` | `"lots of space"` |
| `Wrap("supercalifragilistic is long", 10)` | `"supercalifragilistic\nis long"` |
| `Wrap("héllo wörld", 5)` | `"héllo\nwörld"` |
| `Wrap("", 10)` | `""` |

## Rules
- Words are what `strings.Fields` returns; any run of whitespace (including newlines) becomes a single space or a line break.
- **Greedy**: put as many words on a line as fit. A line holding several words uses single spaces between them.
- Width is measured in **runes** (`utf8.RuneCountInString`), not bytes — `"héllo"` is 5 wide.
- A word longer than `width` is never split; it goes on a line by itself.
- No leading or trailing spaces on any line, and no trailing `"\n"`.
- `width` is always at least 1.
