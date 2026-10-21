# Unique Values & Anagrams

Sorting puts data into a **canonical form**: once sorted, duplicates sit next to each other and two anagrams look identical.

```go
func Unique(nums []int) []int
func CountDistinct(words []string) int
func SortedLetters(s string) string
func IsAnagram(a, b string) bool
```

| Function | Returns |
|---|---|
| `Unique` | the distinct values of `nums`, ascending |
| `CountDistinct` | how many different words appear |
| `SortedLetters` | the runes of `s` in ascending order, as a string |
| `IsAnagram` | whether `a` and `b` contain exactly the same runes, each the same number of times |

| Call | Result |
|---|---|
| `Unique([]int{3, 1, 3, 2, 1})` | `[1 2 3]` |
| `Unique(nil)` | `[]` (length 0) |
| `CountDistinct([]string{"go", "zig", "go"})` | `2` |
| `SortedLetters("banana")` | `"aaabnn"` |
| `SortedLetters("héllo")` | `"hlloé"` (`é` is U+00E9, after ASCII) |
| `IsAnagram("listen", "silent")` | `true` |
| `IsAnagram("aab", "abb")` | `false` |
| `IsAnagram("Tea", "eat")` | `false` (case-sensitive) |

## Rules
- Inputs must not be modified.
- `Unique`: `slices.Clone` + `slices.Sort` + `slices.Compact`. `CountDistinct` reuses the same idea on strings.
- `SortedLetters` works on **runes**, not bytes: `[]rune(s)`, sort, convert back with `string(...)`.
- `IsAnagram` compares the sorted forms. Every rune counts, including spaces and punctuation.
