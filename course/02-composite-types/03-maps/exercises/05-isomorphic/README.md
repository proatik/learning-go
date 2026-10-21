# Isomorphic Strings and Word Patterns

Two strings are **isomorphic** if you can replace every occurrence of each character of `s` with another character to get `t`, where:
- the same character always maps to the same character, and
- two different characters never map to the same one (the mapping is one-to-one).

```go
func IsIsomorphic(s, t string) bool
func WordPattern(pattern, s string) bool
```

## 1. `IsIsomorphic` (compare runes)

| `s` | `t` | Result | Why |
|---|---|---|---|
| `"egg"` | `"add"` | `true` | e→a, g→d |
| `"paper"` | `"title"` | `true` | p→t, a→i, e→l, r→e |
| `"foo"` | `"bar"` | `false` | o would map to both a and r |
| `"ab"` | `"aa"` | `false` | a and b would both map to a |
| `"日本日"` | `"xyx"` | `true` | runes, not bytes |
| `"abc"` | `"ab"` | `false` | different number of runes |
| `""` | `""` | `true` | |

## 2. `WordPattern`

Same idea, but each rune of `pattern` maps to a **word** of `s` (words split with `strings.Fields`).

| `pattern` | `s` | Result |
|---|---|---|
| `"abba"` | `"dog cat cat dog"` | `true` |
| `"abba"` | `"dog cat cat fish"` | `false` |
| `"aaaa"` | `"dog cat cat dog"` | `false` |
| `"abba"` | `"dog dog dog dog"` | `false` |
| `"abc"` | `"dog cat"` | `false` (3 letters, 2 words) |

## Rules
- Use **two maps**, one per direction (`s`→`t` and `t`→`s`), and check both at every position.
- Convert to `[]rune` to compare lengths and index runes side by side.
- Return `false` as soon as a conflict appears.
