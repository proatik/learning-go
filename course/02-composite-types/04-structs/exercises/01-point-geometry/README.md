# Point Geometry

Write plain functions that take and return `Point` **values**.

```go
type Point struct {
	X, Y float64
}

func Distance(p, q Point) float64              // straight-line distance
func Midpoint(p, q Point) Point                // point halfway between p and q
func Translate(p Point, dx, dy float64) Point  // p moved by (dx, dy)
func Centroid(points []Point) (Point, bool)    // average of all points
```

| Call | Result |
|---|---|
| `Distance(Point{0, 0}, Point{3, 4})` | `5` |
| `Distance(Point{1, 1}, Point{1, 1})` | `0` |
| `Midpoint(Point{0, 0}, Point{4, -2})` | `Point{X: 2, Y: -1}` |
| `Translate(Point{1, 2}, 10, -5)` | `Point{X: 11, Y: -3}` |
| `Centroid([]Point{{0, 0}, {4, 0}, {2, 6}})` | `Point{X: 2, Y: 2}, true` |
| `Centroid(nil)` | `Point{}, false` |

## Rules
- Use **keyed** literals when building results: `Point{X: ..., Y: ...}`.
- `Translate` returns a new value; the argument is a copy anyway, so the caller's point never changes.
- `Centroid` returns `false` (and the zero `Point`) for an empty slice.
- `math.Hypot` or `math.Sqrt` are both fine for `Distance`.
