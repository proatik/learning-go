# Map of Slices: Positions and Invert

Practise the grouping idiom `m[k] = append(m[k], v)`.

```go
func Positions(words []string) map[string][]int
func Invert(m map[string]int) map[int][]string
```

## 1. `Positions`

For each word, the indices where it occurs, in increasing order.

| Input | Output |
|---|---|
| `{"to", "be", "or", "not", "to", "be"}` | `map[be:[1 5] not:[3] or:[2] to:[0 4]]` |
| `{"a"}` | `map[a:[0]]` |
| `{}` or `nil` | empty, non-nil map |

## 2. `Invert`

Swap keys and values. Several keys can share a value, so each new value is a **slice** of the old keys, **sorted** ascending.

| Input | Output |
|---|---|
| `map[alice:90 bob:85 carol:90]` | `map[85:[bob] 90:[alice carol]]` |
| `map[x:1]` | `map[1:[x]]` |
| `map[]` or `nil` | empty, non-nil map |

## Rules
- Never initialize the slices yourself: `append` to the (possibly nil) `m[k]` and assign it back.
- `Invert` ranges over a map, so the slices come out in random order: sort each one (`slices.Sort`) before returning.
- Neither function may modify its input.
