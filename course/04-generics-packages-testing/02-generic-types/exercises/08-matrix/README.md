# Numeric Matrix

A `Matrix[T Number]` that works for `int`, `float64` and your own numeric types. The `Number` constraint is what lets methods use `+` and `*` on `T`.

```go
type Number interface {
	~int | ~int32 | ~int64 | ~float32 | ~float64
}

type Matrix[T Number] struct {
	rows, cols int
	data       []T // row-major: element (i, j) is data[i*cols+j]
}

var ErrShape = errors.New("matrix: shape mismatch")

func FromRows[T Number](rows [][]T) (*Matrix[T], error)
func Identity[T Number](n int) *Matrix[T]

func (m *Matrix[T]) Dims() (rows, cols int)
func (m *Matrix[T]) At(i, j int) T
func (m *Matrix[T]) Rows() [][]T
func (m *Matrix[T]) Add(o *Matrix[T]) (*Matrix[T], error)
func (m *Matrix[T]) Mul(o *Matrix[T]) (*Matrix[T], error)
func (m *Matrix[T]) Scale(k T) *Matrix[T]
func (m *Matrix[T]) Transpose() *Matrix[T]
```

## Examples

| Code | Result |
|---|---|
| `a, _ := FromRows([][]int{{1, 2}, {3, 4}})` | `a.Dims()` = `2, 2`; `a.At(1, 0)` = `3` |
| `FromRows([][]int{{1, 2}, {3}})` | `nil`, error matching `ErrShape` |
| `a.Add(a)` | `[[2 4] [6 8]]` |
| `a.Mul(Identity[int](2))` | `[[1 2] [3 4]]` |
| `A` 2×3 times `B` 3×2 | a 2×2 result |
| `A.Add(B)` with 2×3 and 3×2 | error `"add 2x3 and 3x2: matrix: shape mismatch"` |
| `a.Scale(10).Rows()` | `[[10 20] [30 40]]` |
| `FromRows([][]float64{{1, 2, 3}}).Transpose()` | `[[1] [2] [3]]` |

## Rules
- `FromRows` copies the input. It returns `ErrShape` (unwrapped is fine) if there are no rows, the first row is empty, or the rows have different lengths.
- `Identity(n)` is the n×n matrix with `1` on the diagonal; `Identity(0)` (or negative) returns a 0×0 matrix.
- `Rows` returns a **new** `[][]T`; changing it must not change the matrix.
- `Add` needs equal shapes; `Mul` needs `m.cols == o.rows` and returns an `m.rows × o.cols` matrix.
- On a shape error, `Add` / `Mul` return `nil` and an error **wrapping** `ErrShape` with exactly this text: `"add RxC and RxC: matrix: shape mismatch"` / `"mul RxC and RxC: …"` (left operand first). Use `fmt.Errorf` with `%w`.
- No method modifies its receiver or its argument: results are new matrices.
- `At` may panic on out-of-range indexes (like slice indexing).
