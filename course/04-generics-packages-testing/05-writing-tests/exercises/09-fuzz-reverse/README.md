# A fuzz test

Write `FuzzReverse`: a seed corpus plus **properties** that every correct `Reverse` satisfies. Its seeds must catch every broken version.

```go
// Reverse returns s with its runes in reverse order.
func Reverse(s string) string
```

`Reverse` is already written (in the read-only checker). You write the fuzz test in `rev_test.go`.

## Examples

| Call | Result |
|---|---|
| `Reverse("")` | `""` |
| `Reverse("hello world")` | `"dlrow olleh"` |
| `Reverse("héllo, 世界")` | `"界世 ,olléh"` |

## Rules

- **Seeds** (`f.Add`): at least `""`, a word, a sentence with spaces, and non-ASCII text (`é`, `世`, an emoji).
- **Skip** input that is not valid UTF-8: `if !utf8.ValidString(s) { t.Skip(...) }`.
- **Check three properties** for every input `s` (you don't know the expected output, so you can't compare to a fixed `want`):
  1. **Round trip:** `Reverse(Reverse(s)) == s`.
  2. **Shape:** `Reverse(s)` is valid UTF-8 and has the same byte length as `s`.
  3. **Ends swap:** for non-empty `s`, the first rune of `Reverse(s)` is the last rune of `s` (`utf8.DecodeRuneInString`, `utf8.DecodeLastRuneInString`).
- The round trip alone is not enough: some broken versions survive it.
- **The checker** runs only your seeds (no `-fuzz`): they must pass on the correct `Reverse` and fail for every planted bug.
