# Ranked Students

Two **parallel slices** describe a class: `names[i]` scored `scores[i]`. Rank them without structs by sorting **indices**.

```go
func Rank(names []string, scores []int) ([]string, error)
func Standings(names []string, scores []int) ([]string, error)
```

**Order:** score **descending**, then name **ascending** (`strings.Compare`) for equal scores.

- `Rank` returns the names in that order.
- `Standings` returns one line per student, `"<rank>. <name> (<score>)"`, using **competition ranking**: equal scores share a rank, and the next rank skips accordingly (1, 1, 3, ...).
- Both return `nil` and an error if `len(names) != len(scores)`.

```go
names := []string{"cat", "amy", "bob", "dan"}
scores := []int{80, 95, 95, 70}

Rank(names, scores)
// [amy bob cat dan], nil

Standings(names, scores)
// ["1. amy (95)", "1. bob (95)", "3. cat (80)", "4. dan (70)"], nil
```

| Call | Result |
|---|---|
| `Rank([]string{"x", "y"}, []int{1, 2})` | `[y x], nil` |
| `Rank(nil, nil)` | `[], nil` |
| `Rank([]string{"x"}, []int{1, 2})` | `nil, <error>` |
| `Standings([]string{"a", "b", "c"}, []int{5, 5, 5})` | `["1. a (5)", "1. b (5)", "1. c (5)"], nil` |

## Rules
- **Do not modify** `names` or `scores`: sorting one of them would break the pairing with the other.
- Build `order := []int{0, 1, ..., n-1}`, sort it with `slices.SortFunc` and a comparator that looks up `scores[a]`, `names[a]`, and so on. Then map indices to results.
- One `cmp.Or(...)` expresses the whole order.
- Format lines with `fmt.Sprintf`.
