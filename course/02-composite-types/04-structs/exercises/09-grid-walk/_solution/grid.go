package grid

import (
	"fmt"
	"strings"
)

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
	switch dir {
	case 'U':
		p.Y++
	case 'D':
		p.Y--
	case 'R':
		p.X++
	case 'L':
		p.X--
	default:
		return p, false
	}
	return p, true // p is our own copy
}

// Run follows path from (0, 0). It returns Walk{} and an error if path
// contains an invalid direction.
func Run(path string) (Walk, error) {
	var p Pos // the zero Pos is the start (0, 0)
	w := Walk{Visits: map[Pos]int{p: 1}}
	for i, r := range path {
		next, ok := Move(p, r)
		if !ok {
			return Walk{}, fmt.Errorf("invalid direction %q at index %d", r, i)
		}
		p = next
		w.Visits[p]++ // a struct key: no string formatting needed
		w.Min = Pos{X: min(w.Min.X, p.X), Y: min(w.Min.Y, p.Y)}
		w.Max = Pos{X: max(w.Max.X, p.X), Y: max(w.Max.Y, p.Y)}
	}
	w.End = p
	return w, nil
}

// FirstRevisit returns the first cell entered for the second time while
// following path, and true; or Pos{} and false if no cell is revisited or
// path is invalid.
func FirstRevisit(path string) (Pos, bool) {
	if _, err := Run(path); err != nil {
		return Pos{}, false // validate the whole path first
	}
	var p Pos
	seen := map[Pos]bool{p: true}
	for _, r := range path {
		p, _ = Move(p, r) // already validated

		if seen[p] {
			return p, true
		}
		seen[p] = true
	}
	return Pos{}, false
}

// Render draws w's bounding box, top row first: 'S' start, 'E' end,
// '#' visited, '.' not visited. It returns "" if w has no visits.
func Render(w Walk) string {
	if len(w.Visits) == 0 {
		return ""
	}
	var rows []string
	for y := w.Max.Y; y >= w.Min.Y; y-- {
		var sb strings.Builder
		for x := w.Min.X; x <= w.Max.X; x++ {
			p := Pos{X: x, Y: y}
			switch {
			case p == Pos{}: // comparing structs with ==
				sb.WriteByte('S')
			case p == w.End:
				sb.WriteByte('E')
			case w.Visits[p] > 0:
				sb.WriteByte('#')
			default:
				sb.WriteByte('.')
			}
		}
		rows = append(rows, sb.String())
	}
	return strings.Join(rows, "\n")
}
