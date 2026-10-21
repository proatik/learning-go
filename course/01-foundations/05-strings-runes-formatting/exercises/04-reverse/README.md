# Reverse Runes and Words

Implement:

```go
func Reverse(s string) string
func ReverseWords(s string) string
```

`Reverse` reverses the **runes** of `s` (not the bytes — reversing bytes would scramble multi-byte characters).

| Input | Output |
|---|---|
| `"Hello"` | `"olleH"` |
| `"héllo, 世界"` | `"界世 ,olléh"` |
| `"🙂👍"` | `"👍🙂"` |
| `""` | `""` |

`ReverseWords` reverses the order of the **words** (split on whitespace) and joins them with single spaces.

| Input | Output |
|---|---|
| `"one two three"` | `"three two one"` |
| `"  Go   is\tfun "` | `"fun is Go"` |
| `"世界 héllo"` | `"héllo 世界"` |
| `""` | `""` |

## Rules
- `Reverse`: convert to `[]rune`, swap from both ends, convert back.
- `ReverseWords`: use `strings.Fields`; walk the words from last to first.
