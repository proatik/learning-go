package islands

import "slices"

// cell is a grid coordinate.
type cell struct{ r, c int }

// dirs lists the four neighbour offsets.
var dirs = [4]cell{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

// CountIslands returns the number of 4-connected groups of '#' cells.
func CountIslands(grid []string) int {
	return len(IslandSizes(grid))
}

// IslandSizes returns the size of every island, largest first.
func IslandSizes(grid []string) []int {
	land := func(p cell) bool {
		return p.r >= 0 && p.r < len(grid) && p.c >= 0 && p.c < len(grid[p.r]) &&
			grid[p.r][p.c] == '#'
	}
	visited := make(map[cell]bool)
	sizes := []int{}
	for r, row := range grid {
		for c := range len(row) {
			start := cell{r, c}
			if !land(start) || visited[start] {
				continue
			}
			// Depth-first search with a slice as the stack.
			size := 0
			stack := []cell{start}
			visited[start] = true
			for len(stack) > 0 {
				p := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				size++
				for _, d := range dirs {
					n := cell{p.r + d.r, p.c + d.c}
					if land(n) && !visited[n] {
						visited[n] = true // mark when pushed, so each cell is pushed once
						stack = append(stack, n)
					}
				}
			}
			sizes = append(sizes, size)
		}
	}
	slices.Sort(sizes)
	slices.Reverse(sizes)
	return sizes
}

// FloodFill replaces the region containing (row, col) with fill and
// returns the number of cells changed.
func FloodFill(grid [][]byte, row, col int, fill byte) int {
	inside := func(p cell) bool {
		return p.r >= 0 && p.r < len(grid) && p.c >= 0 && p.c < len(grid[p.r])
	}
	start := cell{row, col}
	if !inside(start) {
		return 0
	}
	target := grid[row][col]
	if target == fill {
		return 0
	}
	// Painting a cell changes it from target, which doubles as "visited".
	changed := 0
	grid[row][col] = fill
	stack := []cell{start}
	for len(stack) > 0 {
		p := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		changed++
		for _, d := range dirs {
			n := cell{p.r + d.r, p.c + d.c}
			if inside(n) && grid[n.r][n.c] == target {
				grid[n.r][n.c] = fill
				stack = append(stack, n)
			}
		}
	}
	return changed
}
