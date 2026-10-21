# Matrix: Transpose and Spiral Order

Work with 2D slices (`[][]int`, a slice of rows):

```go
func Transpose(m [][]int) [][]int
func Spiral(m [][]int) []int
```

`Transpose` turns an `r × c` matrix into a new `c × r` matrix: `t[j][i] == m[i][j]`.
`Spiral` lists all elements clockwise, starting at the top-left corner and spiralling inwards.

```
m = [[1 2 3]          Transpose(m) = [[1 4]
     [4 5 6]]                         [2 5]
                                      [3 6]]

m = [[ 1  2  3  4]
     [ 5  6  7  8]    Spiral(m) = [1 2 3 4 8 12 11 10 9 5 6 7]
     [ 9 10 11 12]]
```

| Call | Result |
|---|---|
| `Transpose([][]int{{1, 2, 3}})` | `[[1] [2] [3]]` |
| `Transpose([][]int{{7}})` | `[[7]]` |
| `Transpose(nil)` / `Transpose([][]int{{}, {}})` | empty (len 0) |
| `Spiral([][]int{{1, 2}, {3, 4}})` | `[1 2 4 3]` |
| `Spiral([][]int{{1}, {2}, {3}})` | `[1 2 3]` |
| `Spiral(nil)` | empty (len 0) |

## Rules
- Inputs are **rectangular**: every row has the same length (possibly 0).
- Never modify `m`. `Transpose` must return a **new** matrix: allocate every row with `make`.
- `Spiral` returns a new slice with exactly `rows*cols` elements; preallocate it.
- Hint for `Spiral`: track four boundaries (`top`, `bottom`, `left`, `right`) and shrink them after each side.
