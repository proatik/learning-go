# Matrix

An encapsulated matrix: hidden storage in one flat slice, constructors that validate, and methods for access, arithmetic and printing.

```go
type Matrix struct {
	rows, cols int
	data       []float64 // row-major: element (i, j) is data[i*cols+j]
}

func New(rows, cols int) *Matrix       // all zeros
func FromRows(rows [][]float64) *Matrix // copies the values
func Identity(n int) *Matrix

func (m *Matrix) Rows() int
func (m *Matrix) Cols() int
func (m *Matrix) At(i, j int) (float64, bool)
func (m *Matrix) Set(i, j int, v float64) bool
func (m *Matrix) Equal(o *Matrix) bool
func (m *Matrix) Scale(k float64)          // in place
func (m *Matrix) Transpose() *Matrix
func (m *Matrix) Add(o *Matrix) *Matrix
func (m *Matrix) Mul(o *Matrix) *Matrix
func (m *Matrix) String() string
```

## Examples

| Code | Result |
|---|---|
| `New(2, 3)` | 2×3 of zeros |
| `New(0, 3)` | `nil` |
| `FromRows([][]float64{{1, 2}, {3}})` | `nil` (ragged) |
| `Identity(2).String()` | `"[1 0]\n[0 1]"` |
| `a := FromRows([][]float64{{1, 2}, {3, 4}}); a.At(1, 0)` | `3, true` |
| `a.At(2, 0)` | `0, false` |
| `a.Transpose().String()` | `"[1 3]\n[2 4]"` |
| `a.Mul(Identity(2)).Equal(a)` | `true` |
| `a.Mul(New(3, 1))` | `nil` (2 cols vs 3 rows) |
| `FromRows([][]float64{{1.5, -2}}).String()` | `"[1.5 -2]"` |

## Rules
- `New` and `Identity` return `nil` for a non-positive size. `FromRows` returns `nil` for no rows, an empty row, or rows of different lengths, and **copies** the values (changing the input later must not affect the matrix).
- `At` / `Set` return `false` (and `0` / change nothing) for out-of-range indices.
- `Equal`: same shape and same values. Two nil matrices are equal; nil and non-nil are not.
- `Rows` and `Cols` return `0` for a nil matrix.
- `Transpose`, `Add` and `Mul` return a **new** matrix and leave the operands unchanged. `Add` needs equal shapes; `Mul` needs `m.Cols() == o.Rows()`; otherwise (or when `o` is nil) they return `nil`.
- `String`: one line per row, `[` values separated by single spaces `]`, lines joined by `\n` (no trailing newline). Format each value with `strconv.FormatFloat(v, 'g', -1, 64)`.
- The tests live in package `matrix_test`, so only the exported API is available.
