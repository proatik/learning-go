# Group Anagrams

Two words are **anagrams** if they use the same letters the same number of times (`"listen"` / `"silent"`). Group a list of words by anagram class:

```go
func Key(word string) string
func GroupAnagrams(words []string) [][]string
```

## 1. `Key`

A canonical form shared by all anagrams: the word **lowercased**, with its runes **sorted**.

| Input | Output |
|---|---|
| `"listen"` | `"eilnst"` |
| `"Silent"` | `"eilnst"` |
| `"día"` | `"adí"` |
| `""` | `""` |

## 2. `GroupAnagrams`

Use `Key` as a map key to collect the words, then return the groups in a **deterministic** order:
1. each group's words are sorted ascending (byte-wise string order, so uppercase sorts before lowercase);
2. groups are ordered by their first (smallest) word.

| Input | Output |
|---|---|
| `{"eat", "tea", "tan", "ate", "nat", "bat"}` | `[[ate eat tea] [bat] [nat tan]]` |
| `{"Listen", "silent", "enlist", "google"}` | `[[Listen enlist silent] [google]]` |
| `{"a", "a"}` | `[[a a]]` (duplicates stay) |
| `{}` or `nil` | empty (`nil` is fine) |

## Rules
- A slice can't be a map key, but a `string` can: convert the sorted `[]rune` back with `string(rs)`.
- The original spelling goes into the group; only the key is lowercased.
- `words` must not be modified.
- Sort with `slices.Sort` and `slices.SortFunc` + `cmp.Compare` (lesson 02).
