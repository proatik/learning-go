package grid

// Pos is a cell on an infinite grid; Y grows upwards.
type Pos struct {
	X, Y int
}

// Walk is the result of following a path from (0, 0).
type Walk struct {
	Visits   map[Pos]int // how often each cell was occupied (the start counts once)
	End      Pos         // final position
	Min, Max Pos         // bounding box of all visited cells (inclusive)
}

// Move returns p moved one step in direction dir ('U', 'D', 'L' or 'R') and
// true, or p unchanged and false for any other rune.
func Move(p Pos, dir rune) (Pos, bool) {
	// TODO
	return Pos{}, false
}

// Run follows path from (0, 0). It returns Walk{} and an error if path
// contains an invalid direction.
func Run(path string) (Walk, error) {
	// TODO
	return Walk{}, nil
}

// FirstRevisit returns the first cell entered for the second time while
// following path, and true; or Pos{} and false if no cell is revisited or
// path is invalid.
func FirstRevisit(path string) (Pos, bool) {
	// TODO
	return Pos{}, false
}

// Render draws w's bounding box, top row first: 'S' start, 'E' end,
// '#' visited, '.' not visited. It returns "" if w has no visits.
func Render(w Walk) string {
	// TODO
	return ""
}
