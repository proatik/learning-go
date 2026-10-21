# Sort by Length, then Alphabetically

Write multi-key comparison functions with `cmp.Compare`, `strings.Compare` and `cmp.Or`.

```go
func SortByLength(words []string)
func SortByLengthDesc(words []string)
func Shortest(words []string, n int) []string
```

**Length** means the number of **runes** (`utf8.RuneCountInString`), so `"héllo"` has length 5.

| Function | Order |
|---|---|
| `SortByLength` | sorts `words` **in place**: shortest first; equal lengths alphabetically (`strings.Compare`) |
| `SortByLengthDesc` | sorts `words` **in place**: longest first; equal lengths still alphabetically **ascending** |
| `Shortest` | returns a **new** slice with the `n` first words in `SortByLength` order; `words` is not modified |

| Input | After `SortByLength` | After `SortByLengthDesc` |
|---|---|---|
| `[pear fig banana kiwi apple]` | `[fig kiwi pear apple banana]` | `[banana apple kiwi pear fig]` |
| `[bb a ccc aa]` | `[a aa bb ccc]` | `[ccc aa bb a]` |

| Call | Result |
|---|---|
| `Shortest([]string{"pear", "fig", "banana", "kiwi"}, 2)` | `[fig kiwi]` |
| `Shortest([]string{"b", "a"}, 5)` | `[a b]` (fewer than `n` words: return all) |
| `Shortest([]string{"b", "a"}, 0)` | `[]` (length 0; also for negative `n`) |

## Rules
- Use `slices.SortFunc` and **one** `cmp.Or(...)` per comparator: no `if` chains.
- Descending on one key only: swap the arguments of **that** key.
- `Shortest` must not reuse `SortByLength` on the caller's slice: `Clone` first. Return a slice of the copy.
