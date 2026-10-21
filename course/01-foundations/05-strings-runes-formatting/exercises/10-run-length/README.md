# Run-Length Encoding

Run-length encoding replaces each run of the same rune by the rune followed by the run's length. Implement:

```go
func Encode(s string) string
func Decode(s string) (string, error)
```

| `Encode` input | Output |
|---|---|
| `"aaabcc"` | `"a3b1c2"` |
| `"abc"` | `"a1b1c1"` |
| `"wwwwwwwwwwwwb"` | `"w12b1"` |
| `"ééé世世"` | `"é3世2"` |
| `"aabbaa"` | `"a2b2a2"` |
| `""` | `""` |

| `Decode` input | Output |
|---|---|
| `"a3b1c2"` | `"aaabcc", nil` |
| `"w12b1"` | `"wwwwwwwwwwwwb", nil` |
| `" 2!3"` | `"  !!!", nil` |
| `""` | `"", nil` |
| `"3a"` | `"", error` — starts with a count |
| `"a3b"` | `"", error` — `b` has no count |
| `"a0"` | `"", error` — count must be at least 1 |

## Rules
- Work with **runes**: `é` and `世` are single symbols.
- `Encode` input never contains the digits `0`–`9` (they would be ambiguous). Spaces and punctuation are ordinary symbols.
- In `Decode`, a count is one or more digits `0`–`9`, and every symbol must be followed by a count ≥ 1. On any error return `""` and an error from `fmt.Errorf` that mentions the problem.
- Use `strings.Builder` for both functions. `strconv.Itoa` and `strings.Repeat` are useful.
- `Decode(Encode(s))` must give back `s`.
