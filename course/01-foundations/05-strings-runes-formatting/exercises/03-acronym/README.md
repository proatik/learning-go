# Acronym

Implement:

```go
func Acronym(phrase string) string
```

Build an acronym from the **first letter of each word**, upper-cased.

| Input | Output |
|---|---|
| `"Portable Network Graphics"` | `"PNG"` |
| `"as soon as possible"` | `"ASAP"` |
| `"Complementary metal-oxide semiconductor"` | `"CMOS"` |
| `"  Ruby   on\tRails "` | `"ROR"` |
| `"Rock & Roll"` | `"RR"` |
| `"The Road _Not_ Taken"` | `"TRNT"` |
| `"élan vital"` | `"ÉV"` |
| `""` | `""` |

## Rules
- Words are separated by any whitespace **or** a hyphen `-`.
- A word's "first letter" is its first rune for which `unicode.IsLetter` is true; words with no letters (like `&`) contribute nothing.
- Upper-case with `unicode.ToUpper`, and build the result with a `strings.Builder`.
