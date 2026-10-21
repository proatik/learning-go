# Double: Arrays vs Slices

Arrays are **values** (passing one copies every element); slices **share** their backing array. Feel the difference by doubling numbers three ways.

```go
type Quad [4]int

func DoubleSlice(s []int)
func DoubleArray(a *Quad)
func Doubled(a Quad) Quad
```

| Function | Changes the caller's data? | How |
|---|---|---|
| `DoubleSlice(s)` | ✅ | `s[i]` writes into the shared backing array |
| `DoubleArray(&q)` | ✅ | a pointer to the array; `a[i]` auto-dereferences |
| `Doubled(q)` | ❌ | `a` is a copy; return a doubled copy, the caller's `q` stays as is |

| Example | After |
|---|---|
| `s := []int{1, 2, 3}; DoubleSlice(s)` | `s == [2 4 6]` |
| `DoubleSlice(s[1:])` on `[1 2 3]` | `s == [1 4 6]` (a sub-slice shares the same array) |
| `q := Quad{1, 2, 3, 4}; DoubleArray(&q)` | `q == [2 4 6 8]` |
| `r := Doubled(q)` on `[1 2 3 4]` | `r == [2 4 6 8]`, `q == [1 2 3 4]` |

## Rules
- Loop with `for i := range ...` and assign by index: `for _, v := range` gives copies.
- `DoubleSlice(nil)` and `DoubleSlice([]int{})` do nothing.
- `DoubleArray(nil)` does nothing (no panic).
