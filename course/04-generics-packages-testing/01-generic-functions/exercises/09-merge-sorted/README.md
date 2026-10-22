# Merge Sorted Slices

Merge already-sorted slices in one pass, ordered by a `cmp.Compare`-style comparator, so it works for any `T`, including structs.

```go
func Merge[T any](a, b []T, compare func(x, y T) int) []T
func MergeOrdered[T cmp.Ordered](a, b []T) []T
func MergeAll[T any](compare func(x, y T) int, lists ...[]T) []T
func Desc[T any](compare func(x, y T) int) func(x, y T) int
```

A comparator returns a **negative** number when `x` comes first, **0** when they are equal, **positive** when `y` comes first. `cmp.Compare` is one.

## Examples

| Call | Result |
|---|---|
| `Merge([]int{1, 4, 9}, []int{2, 3, 10}, cmp.Compare[int])` | `[1 2 3 4 9 10]` |
| `MergeOrdered([]string{"ant", "fox"}, []string{"bee"})` | `[ant bee fox]` |
| `MergeOrdered([]int{}, []int{5})` | `[5]` |
| `MergeAll(cmp.Compare[int], []int{1, 7}, []int{3}, []int{2, 8})` | `[1 2 3 7 8]` |
| `Desc(cmp.Compare[int])(1, 2)` | positive (in descending order, `2` comes first) |
| `Merge([]int{9, 4}, []int{5, 1}, Desc(cmp.Compare[int]))` | `[9 5 4 1]` |

Stability with `byAge := func(x, y Person) int { return cmp.Compare(x.Age, y.Age) }`:

| `a` | `b` | `Merge(a, b, byAge)` |
|---|---|---|
| `[{ann 30}]` | `[{bob 30}]` | `[{ann 30} {bob 30}]` (equal: `a` first) |

## Rules
- Assume each input is already sorted by `compare`. Don't sort, and don't call `slices.Sort*`: walk both slices with two indexes, **O(len(a) + len(b))**.
- **Stable:** when `compare(a[i], b[j]) == 0`, take `a[i]` first. Elements from the same input keep their order.
- The result is a **new** slice of length `len(a) + len(b)` (length `0` for two empty inputs). Never modify the inputs.
- `MergeOrdered`: `Merge` with `cmp.Compare[T]` as the comparator (a generic function used as a value must be instantiated).
- `MergeAll`: merge every list; on ties, elements of earlier lists come first. No lists → length `0`.
- `Desc(compare)`: a comparator that orders the opposite way. `Desc(c)(x, y)` is `c(y, x)`.
