# Character Frequency

Two functions that count **runes** (not bytes):

```go
func CharFrequency(s string) map[rune]int
func FirstUnique(s string) (r rune, ok bool)
```

## 1. `CharFrequency`

Count every rune in `s` **except whitespace**. Counting is case-sensitive.

| Input | Output |
|---|---|
| `"hello"` | `map['e':1 'h':1 'l':2 'o':1]` |
| `"a a A"` | `map['A':1 'a':2]` |
| `"héé"` | `map['h':1 'é':2]` |
| `""` or `"  \n"` | empty, non-nil map |

## 2. `FirstUnique`

Return the **first** rune of `s` (in string order) that occurs exactly once, with `ok == true`. Whitespace counts like any other rune here. If there is none, return `0, false`.

| Input | `r` | `ok` |
|---|---|---|
| `"swiss"` | `'w'` | `true` |
| `"aabbc"` | `'c'` | `true` |
| `"héllo hé"` | `'o'` | `true` |
| `"aabb"` | `0` | `false` |
| `""` | `0` | `false` |

## Rules
- `range` over a string yields runes; use `unicode.IsSpace` to skip whitespace in `CharFrequency`.
- `FirstUnique` needs **two passes**: count everything, then walk `s` again in order. Do not range over the map to find the answer: its order is random.
