# Islands and Flood Fill

Explore 2D grids: count the islands on a map and implement a paint-bucket flood fill.

```go
func CountIslands(grid []string) int
func IslandSizes(grid []string) []int
func FloodFill(grid [][]byte, row, col int, fill byte) int
```

Neighbours are the 4 cells **up, down, left, right** (no diagonals).

## Examples

```go
grid := []string{
	"##..#",
	"#...#",
	"..#..",
	".....",
	"###.#",
}
CountIslands(grid) // 5
IslandSizes(grid)  // [3 3 2 1 1]

img := [][]byte{
	[]byte("aab"),
	[]byte("abb"),
	[]byte("bba"),
}
FloodFill(img, 0, 2, 'x') // 5 — img is now "aax", "axx", "xxa"
```

## Rules

**CountIslands / IslandSizes**
- `'#'` is land; **any other byte** is water. An island is a maximal 4-connected group of land cells.
- Rows may have **different lengths**; a cell beyond the end of its row is water. An empty or nil grid has no islands.
- `IslandSizes` returns the number of cells of each island, sorted **descending**. No islands → empty slice.
- The input is not modified (strings are immutable, so you need your own `visited` record).

**FloodFill**
- Replace the cell `(row, col)` and every cell 4-connected to it that has the **same original byte** with `fill`, in place. Return how many cells changed.
- Rows may have different lengths (as above).
- Return `0` and change nothing if `(row, col)` is outside the grid, or the cell already equals `fill`.
- Must handle large regions (e.g. a 300×300 grid) — an explicit stack or queue (a slice) is recommended over deep recursion.
