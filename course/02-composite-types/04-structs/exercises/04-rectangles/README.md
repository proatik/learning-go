# Rectangles

A `Rect` is built from two nested `Point` structs. Like `image.Rectangle` in the standard library, it is **half-open**: `Min` is inside, `Max` is just outside.

```go
type Point struct {
	X, Y int
}

// Rect contains the points p with Min.X <= p.X < Max.X and Min.Y <= p.Y < Max.Y.
type Rect struct {
	Min, Max Point
}

func Area(r Rect) int
func Contains(r Rect, p Point) bool
func Intersect(a, b Rect) (Rect, bool)
func Bounds(rs []Rect) Rect
```

`r := Rect{Min: Point{0, 0}, Max: Point{4, 3}}` (4 wide, 3 high):

| Call | Result |
|---|---|
| `Area(r)` | `12` |
| `Contains(r, Point{0, 0})` | `true` (Min is inside) |
| `Contains(r, Point{4, 1})` | `false` (Max.X is outside) |
| `Intersect(r, Rect{Point{2, 1}, Point{6, 5}})` | `Rect{Min: {2 1}, Max: {4 3}}, true` |
| `Intersect(r, Rect{Point{4, 0}, Point{6, 3}})` | `Rect{}, false` (only touching edges) |
| `Bounds([]Rect{r, {Point{-1, 2}, Point{1, 5}}})` | `Rect{Min: {-1 0}, Max: {4 5}}` |
| `Bounds(nil)` | `Rect{}` |

## Rules
- You may assume every input `Rect` is well-formed: `Min.X <= Max.X` and `Min.Y <= Max.Y`.
- A rect with zero width or height has area 0 and contains no points.
- `Intersect` returns `Rect{}, false` when the overlap has zero area.
- `Bounds` is the smallest rect containing all of `rs` (empty rects included); `Rect{}` for no rects.
- Build results with keyed literals: `Rect{Min: Point{X: ..., Y: ...}, Max: ...}`.
