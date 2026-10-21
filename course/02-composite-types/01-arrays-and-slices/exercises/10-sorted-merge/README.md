# Sorted Slices: Merge, Dedupe, Intersect

All inputs are sorted in ascending order (duplicates allowed). Use the **two-pointer** technique: walk both slices with one index each, never sorting.

```go
func MergeSorted(a, b []int) []int
func Dedupe(s []int) []int
func Intersect(a, b []int) []int
```

| Call | Result |
|---|---|
| `MergeSorted([]int{1, 4, 9}, []int{2, 3, 10})` | `[1 2 3 4 9 10]` |
| `MergeSorted([]int{1, 1}, []int{1})` | `[1 1 1]` |
| `MergeSorted(nil, []int{5})` | `[5]` |
| `Dedupe([]int{1, 1, 2, 3, 3, 3})` | `[1 2 3]` |
| `Dedupe([]int{4})` | `[4]` |
| `Intersect([]int{1, 2, 2, 3, 5}, []int{2, 2, 5, 6})` | `[2 5]` |
| `Intersect([]int{1, 2}, []int{3, 4})` | `[]` |

## Rules
- `MergeSorted`: returns a **new** slice of length `len(a)+len(b)` (allocate once); never modifies `a` or `b`; the result must not share memory with them. Keeps duplicates.
- `Dedupe`: removes repeated values **in place** (the `s[:0]` trick), with **zero allocations**. The result shares `s`'s backing array. Leftover tail contents don't matter.
- `Intersect`: returns a new slice of the values present in **both**, each value **once**, ascending. Never modifies the inputs.
- O(len(a)+len(b)) time: no nested loops over both inputs, no sorting.
- Don't use the `slices` or `sort` packages.
