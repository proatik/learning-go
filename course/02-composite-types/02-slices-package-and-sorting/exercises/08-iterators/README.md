# Ranging over Iterators

Practise `for ... range` over `iter.Seq` / `iter.Seq2` values from `slices` and `strings`.

```go
func Evens(nums []int) []int
func LastIndex(words []string, w string) int
func ChunkSums(nums []int, n int) []int
func SortedFields(s string) []string
```

| Function | Returns | Use |
|---|---|---|
| `Evens` | the even numbers of `nums`, in order | `for v := range slices.Values(nums)` |
| `LastIndex` | index of the **last** `w` in `words`, or `-1` | `for i, v := range slices.Backward(words)` |
| `ChunkSums` | the sum of each consecutive group of `n` numbers (the last group may be shorter) | `for chunk := range slices.Chunk(nums, n)` |
| `SortedFields` | the whitespace-separated words of `s`, sorted | `slices.Sorted(strings.FieldsSeq(s))` |

| Call | Result |
|---|---|
| `Evens([]int{1, 2, 3, 4, -6})` | `[2 4 -6]` |
| `LastIndex([]string{"a", "b", "a"}, "a")` | `2` |
| `LastIndex([]string{"a"}, "z")` | `-1` |
| `ChunkSums([]int{1, 2, 3, 4, 5}, 2)` | `[3 7 5]` |
| `ChunkSums([]int{1, 2, 3}, 5)` | `[6]` |
| `ChunkSums([]int{1, 2, 3}, 0)` | `[]` (`n < 1`: no chunks) |
| `SortedFields("  the quick\tbrown  fox ")` | `[brown fox quick the]` |

## Rules
- Use the iterator named in the table for each function; no index-based `for i := 0; ...` loops.
- `LastIndex` returns from **inside** the loop as soon as it finds a match: the iterator stops.
- `slices.Chunk` **panics** if `n < 1`: check first.
- Inputs must not be modified.
