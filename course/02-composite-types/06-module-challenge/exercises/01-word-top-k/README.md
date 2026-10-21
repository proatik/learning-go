# Word Frequency Top-K

Count the words in a text and report the `k` most frequent ones.

```go
type WordCount struct {
	Word  string
	Count int
}

func Counts(text string) map[string]int
func TopK(text string, k int) []WordCount
```

## Examples

| Call | Result |
|---|---|
| `Counts("Go go GO! gopher")` | `map[go:3 gopher:1]` |
| `TopK("b a b c a b", 2)` | `[{b 3} {a 2}]` |
| `TopK("x y z", 2)` | `[{x 1} {y 1}]` (tie → alphabetical) |
| `TopK("one", 5)` | `[{one 1}]` |
| `TopK("anything", 0)` | `[]` (empty) |

## Rules
- A **word** is a maximal run of letters or digits (`unicode.IsLetter` / `unicode.IsDigit`). Everything else separates words: `"don't"` is `"don"` and `"t"`.
- Words are compared case-insensitively and reported in lowercase (`strings.ToLower`).
- `Counts` of a text with no words returns an empty, non-nil map.
- `TopK` orders by **count descending**, then **word ascending** for ties, and returns the first `k`.
- If `k` exceeds the number of distinct words, return them all. If `k <= 0`, return an empty slice.
- Do not modify anything shared between calls; the result is a fresh slice.
