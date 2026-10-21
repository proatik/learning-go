# Word Count

Count how many times each word appears:

```go
func WordCount(s string) map[string]int
```

| Input | Output |
|---|---|
| `"the cat the hat"` | `map[cat:1 hat:1 the:2]` |
| `"Go go GO"` | `map[go:3]` |
| `"  spaced\tout\n words  "` | `map[out:1 spaced:1 words:1]` |
| `""` | `map[]` (empty, **not** nil) |

## Rules
- Words are separated by any whitespace: split with `strings.Fields`.
- Counting is case-insensitive: store every word in lowercase (`strings.ToLower`).
- Punctuation is part of the word (`"go!"` and `"go"` are different words).
- Always return an initialized map, even for empty input.
- Use the counting idiom `m[k]++`: no "does the key exist?" check.
