# Sudoku Checker

Work with a Sudoku grid stored as a fixed-size array. `0` means an empty cell; `g[row][col]`, rows and columns `0`–`8`. The nine **boxes** are the 3×3 squares starting at rows/columns 0, 3 and 6.

```go
func Validate(g [9][9]int) error
func IsSolved(g [9][9]int) bool
func Candidates(g [9][9]int, row, col int) []int
func Place(g *[9][9]int, row, col, v int) error
```

## Examples

```go
var g [9][9]int
Place(&g, 0, 0, 5)     // nil — g[0][0] is now 5
Place(&g, 0, 8, 5)     // error: 5 already in row 0
Place(&g, 4, 0, 5)     // error: 5 already in column 0
Place(&g, 2, 2, 5)     // error: 5 already in the top-left box
Candidates(g, 0, 1)    // [1 2 3 4 6 7 8 9]
Validate(g)            // nil
```

## Rules
- **Validate** returns an error if any value is outside `0`–`9`, or a digit `1`–`9` appears more than once in a row, a column or a box. Zeros never conflict. Otherwise `nil`.
- **IsSolved** reports whether the grid is valid **and** has no zeros.
- **Candidates** returns, in ascending order, the digits that could be written into the empty cell `(row, col)` without a row/column/box conflict. Return an empty (possibly nil) slice if the cell is out of range or not empty.
- **Place** writes `v` into `g[row][col]` through the pointer. Return an error, leaving the grid **unchanged**, if `row` or `col` is outside `0`–`8`, `v` is outside `1`–`9`, the cell is not empty, or `v` already appears in that row, column or box. A `nil` grid pointer is also an error.
- Validate, IsSolved and Candidates receive a **copy** of the grid (arrays are values) and must not need a pointer.
