# Grid Walk

A walker starts at `(0, 0)` and follows a path such as `"UURDL"`. Track the visited cells in a map keyed by a **struct**.

```go
type Pos struct {
	X, Y int
}

type Walk struct {
	Visits   map[Pos]int // how often each cell was occupied (the start counts once)
	End      Pos         // final position
	Min, Max Pos         // bounding box of all visited cells (inclusive)
}

func Move(p Pos, dir rune) (Pos, bool)
func Run(path string) (Walk, error)
func FirstRevisit(path string) (Pos, bool)
func Render(w Walk) string
```

| Direction | Effect |
|---|---|
| `'U'` | `Y + 1` |
| `'D'` | `Y - 1` |
| `'R'` | `X + 1` |
| `'L'` | `X - 1` |

| Call | Result |
|---|---|
| `Move(Pos{2, 3}, 'L')` | `Pos{1, 3}, true` |
| `Move(Pos{2, 3}, 'x')` | `Pos{2, 3}, false` |
| `Run("UURDL")` | `Visits: {0,0}:1 {0,1}:2 {0,2}:1 {1,2}:1 {1,1}:1`, `End: {0,1}`, `Min: {0,0}`, `Max: {1,2}` |
| `Run("")` | `Visits: {0,0}:1`, everything else zero |
| `Run("UUX")` | `Walk{}`, error |
| `FirstRevisit("UURDL")` | `Pos{0, 1}, true` |
| `FirstRevisit("RRR")` | `Pos{}, false` |

`Render(Run("UURDL"))` draws the bounding box, **top row = `Max.Y`**, columns from `Min.X` to `Max.X`:

```
##
E#
S.
```

## Rules
- `Run` returns `Walk{}` and an error naming the bad character and its index (e.g. `invalid direction 'X' at index 2`) for any character other than `U D L R`. Paths are ASCII.
- `FirstRevisit` returns the first cell **entered** a second time (the start cell counts as visited). If the path contains **any** invalid character, return `Pos{}, false`, even if a revisit happened before it.
- `Render`: `S` = start `(0,0)`, `E` = end (when the end is not the start), `#` = other visited cells, `.` = never visited. Rows separated by `"\n"`, no trailing newline. If `w.Visits` is empty, return `""`.
