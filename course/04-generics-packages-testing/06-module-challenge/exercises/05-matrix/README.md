# Generic Matrix

Build a dense `Matrix[T]` for any numeric type — including named ones like `type Cents int64` — with shape-checked arithmetic, a row iterator, and a Go 1.27 **generic method** that converts the element type.

```go
type Number interface {
	~int | ~int32 | ~int64 | ~float32 | ~float64
}

var ErrShape = errors.New("matrix: bad shape")

type Matrix[T Number] struct {
	rows, cols int
	data       []T // element (r, c) is data[r*cols+c]
}

func New[T Number](rows, cols int) (*Matrix[T], error)
func FromRows[T Number](rows [][]T) (*Matrix[T], error)
func Identity[T Number](n int) *Matrix[T]

func (m *Matrix[T]) Dims() (rows, cols int)
func (m *Matrix[T]) At(r, c int) T
func (m *Matrix[T]) Set(r, c int, v T)
func (m *Matrix[T]) Add(o *Matrix[T]) (*Matrix[T], error)
func (m *Matrix[T]) Mul(o *Matrix[T]) (*Matrix[T], error)
func (m *Matrix[T]) Transpose() *Matrix[T]
func (m *Matrix[T]) Sum() T
func (m *Matrix[T]) Apply[U Number](f func(T) U) *Matrix[U]
func (m *Matrix[T]) Rows() iter.Seq2[int, []T]
func (m *Matrix[T]) String() string
```

## Examples

| Code | Result |
|---|---|
| `New[int](0, 3)` | `nil`, error wrapping `ErrShape` |
| `FromRows([][]int{{1, 2}, {3}})` | `nil`, error wrapping `ErrShape` (ragged) |
| `a` = `[[1 2 3] [4 5 6]]`, `b` = `[[7 8] [9 10] [11 12]]`: `a.Mul(b)` | `[[58 64] [139 154]]` |
| `a.Mul(a)` | `nil`, `"matrix: bad shape: 2x3 * 2x3"` |
| `a.Add(Identity[int](2))` | `nil`, `"matrix: bad shape: 2x3 + 2x2"` |
| `a.Transpose()` | `[[1 4] [2 5] [3 6]]` |
| `prices.Apply(func(c Cents) float64 { return float64(c) / 100 })` | a `*Matrix[float64]` |
| `fmt.Sprint(a)` | `"[1 2 3]\n[4 5 6]"` |

## Rules
- **New** needs `rows >= 1` and `cols >= 1`; otherwise return `nil` and an error wrapping `ErrShape`.
- **FromRows** **copies** the input. `nil`, no rows, an empty first row or rows of different lengths → `nil` + `ErrShape`.
- **Identity(n)** returns `nil` if `n < 1`.
- **Add** needs equal dimensions; **Mul** needs `m.cols == o.rows` and returns an `m.rows × o.cols` matrix. Their errors are exactly `fmt.Errorf("%w: %dx%d + %dx%d", ErrShape, ...)` and `"... * ..."` with both shapes.
- **Add**, **Mul**, **Transpose** and **Apply** return **new** matrices and never change `m` or `o`.
- **Rows** yields `(index, copy of the row)` and honours `break`; changing a yielded row must not change the matrix.
- **String** is each row formatted with `fmt.Sprint(row)`, joined by `"\n"`.
