package matrix

// Transpose returns a new matrix t with t[j][i] == m[i][j].
// m must be rectangular. An empty matrix (no rows or no columns) gives an
// empty result.
func Transpose(m [][]int) [][]int {
	if len(m) == 0 || len(m[0]) == 0 {
		return nil
	}
	rows, cols := len(m), len(m[0])
	t := make([][]int, cols)
	for j := range t {
		t[j] = make([]int, rows) // every row needs its own allocation
		for i := range rows {
			t[j][i] = m[i][j]
		}
	}
	return t
}

// Spiral returns the elements of the rectangular matrix m in clockwise
// spiral order, starting at m[0][0].
//
// top/bottom/left/right are the inclusive bounds of the ring not yet
// visited. Each pass walks the ring's four sides and then shrinks it.
func Spiral(m [][]int) []int {
	if len(m) == 0 || len(m[0]) == 0 {
		return nil
	}
	out := make([]int, 0, len(m)*len(m[0]))
	top, bottom := 0, len(m)-1
	left, right := 0, len(m[0])-1
	for top <= bottom && left <= right {
		for j := left; j <= right; j++ { // top row, left → right
			out = append(out, m[top][j])
		}
		top++
		for i := top; i <= bottom; i++ { // right column, top → bottom
			out = append(out, m[i][right])
		}
		right--
		if top <= bottom { // bottom row, right → left (if a row is left)
			for j := right; j >= left; j-- {
				out = append(out, m[bottom][j])
			}
			bottom--
		}
		if left <= right { // left column, bottom → top (if a column is left)
			for i := bottom; i >= top; i-- {
				out = append(out, m[i][left])
			}
			left++
		}
	}
	return out
}
