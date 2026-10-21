# Palindromes

Implement:

```go
func IsPalindrome(s string) bool
```

A palindrome reads the same forwards and backwards. Ignore **case** and every rune that is **not a letter or digit** (spaces, punctuation, symbols).

| Input | Output |
|---|---|
| `"racecar"` | `true` |
| `"A man, a plan, a canal: Panama!"` | `true` |
| `"No 'x' in Nixon"` | `true` |
| `"Gopher"` | `false` |
| `"été"` | `true` |
| `"上海自来水来自海上"` | `true` |
| `"12321"` | `true` |
| `""` | `true` |
| `"?!"` | `true` (nothing left to compare) |

## Rules
- Use `unicode.IsLetter`, `unicode.IsDigit` and `unicode.ToLower`.
- Compare **runes**, not bytes. Accented letters are distinct: `é` ≠ `e`.
