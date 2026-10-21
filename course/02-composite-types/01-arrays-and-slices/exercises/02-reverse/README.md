# Reverse: In Place vs Copy

Reverse a slice **in place**, and reverse an array **by value**:

```go
func Reverse(s []int)
func ReversedArray(a [5]int) [5]int
```

| Call | Effect / result |
|---|---|
| `s := []int{1, 2, 3}; Reverse(s)` | `s` is now `[3 2 1]` |
| `s := []int{1, 2, 3, 4}; Reverse(s)` | `s` is now `[4 3 2 1]` |
| `Reverse(nil)` | no-op, no panic |
| `ReversedArray([5]int{1, 2, 3, 4, 5})` | `[5 4 3 2 1]` |

## Rules
- `Reverse` returns nothing: it modifies the caller's elements through the shared backing array. Don't allocate a new slice.
- `ReversedArray` receives a **copy** of the array; return the reversed copy. The caller's array must be unchanged (it is, automatically — arrays are values).
- Swap with Go's tuple assignment: `x, y = y, x`.
