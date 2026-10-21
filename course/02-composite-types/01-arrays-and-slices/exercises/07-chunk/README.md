# Chunk a Slice

Split a slice into consecutive chunks of at most `size` elements:

```go
func Chunk(s []int, size int) [][]int
```

| Call | Result |
|---|---|
| `Chunk([]int{1, 2, 3, 4, 5}, 2)` | `[[1 2] [3 4] [5]]` |
| `Chunk([]int{1, 2, 3, 4}, 2)` | `[[1 2] [3 4]]` |
| `Chunk([]int{1, 2, 3}, 5)` | `[[1 2 3]]` |
| `Chunk([]int{1, 2, 3}, 1)` | `[[1] [2] [3]]` |
| `Chunk(nil, 3)` / `Chunk([]int{}, 3)` | `[]` (no chunks) |
| `Chunk([]int{1, 2}, 0)` | `nil` (invalid size) |

## Rules
- Each chunk is a **sub-slice of `s`**, not a copy: don't copy elements.
- Each chunk's **capacity must equal its length** (use a full slice expression `s[lo:hi:max]`). Then a caller appending to one chunk can't overwrite the next chunk's elements; the tests check this.
- Return `nil` if `size <= 0`.
- Preallocate the outer slice: you can compute the number of chunks up front.
