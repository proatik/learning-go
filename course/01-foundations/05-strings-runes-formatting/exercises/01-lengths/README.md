# Bytes, Runes and Vowels

Implement:

```go
func Lengths(s string) (bytes, runes int)
func CountVowels(s string) int
```

`Lengths` returns the number of **bytes** and the number of **runes** in `s`.

| Input | `bytes` | `runes` |
|---|---|---|
| `""` | `0` | `0` |
| `"Go"` | `2` | `2` |
| `"héllo"` | `6` | `5` |
| `"世界"` | `6` | `2` |
| `"🙂"` | `4` | `1` |

`CountVowels` counts the vowels `a e i o u` in either case.

| Input | Output |
|---|---|
| `"Gopher"` | `2` |
| `"AEIOU aeiou"` | `10` |
| `"rhythm"` | `0` |
| `"héllo, 世界"` | `1` |

## Rules
- Only the ASCII letters `a e i o u` (and `A E I O U`) are vowels. `é` is **not** counted.
- Use `range` over the string for `CountVowels` so multi-byte runes are handled as one unit.
