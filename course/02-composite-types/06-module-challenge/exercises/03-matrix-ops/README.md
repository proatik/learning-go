# Matrix Operations

Implement three classic operations on `[][]int` matrices (`m[row][col]`).

```go
func Multiply(a, b [][]int) ([][]int, error)
func Rotate(m [][]int) ([][]int, error)
func Spiral(m [][]int) ([]int, error)
```

## Examples

```go
Multiply([][]int{{1, 2}, {3, 4}}, [][]int{{5, 6}, {7, 8}})
// [[19 22] [43 50]], nil

Rotate([][]int{{1, 2, 3}, {4, 5, 6}}) // 90° clockwise
// [[4 1] [5 2] [6 3]], nil

Spiral([][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}})
// [1 2 3 6 9 8 7 4 5], nil
```

## Rules
- A matrix is **rectangular** when every row has the same length. A ragged matrix is an error in all three functions.
- **Multiply**: `a` is `n×k`, `b` must be `k×p`; the result is `n×p` with `c[i][j] = Σ a[i][x]*b[x][j]`. Return `nil` and an error if either matrix is empty (no rows, or rows of length 0), ragged, or the inner dimensions differ.
- **Rotate**: rotate 90° **clockwise**; an `r×c` matrix becomes `c×r`. An empty matrix (no rows, or rows of length 0) returns an empty, non-nil `[][]int{}`.
- **Spiral**: read clockwise from the top-left corner, spiralling inwards. An empty matrix returns an empty, non-nil `[]int{}`.
- On error, Rotate and Spiral return `nil`.
- **Never modify the input**, and the result must not share memory with it — allocate every row of the result with `make`.
