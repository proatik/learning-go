# Rank the Shapes

Sort a slice of mixed shapes through their interface methods, and classify them with a type switch.

```go
type Shape interface {
	Area() float64
	Name() string
}

type Square struct{ Side float64 } // Name() == "square"
type Rect struct{ W, H float64 }   // Name() == "rect"
type Circle struct{ R float64 }    // Name() == "circle"

func SortByArea(shapes []Shape) []Shape
func TopN(shapes []Shape, n int) []Shape
func Kind(s Shape) string
```

## Examples

| Call | Result |
|---|---|
| `SortByArea([]Shape{Rect{2, 3}, Circle{1}, Square{1}})` | `[Square{1} Circle{1} Rect{2, 3}]` (areas 1, π, 6) |
| `SortByArea([]Shape{Square{2}, Rect{1, 4}})` | `[Rect{1, 4} Square{2}]` (tie 4 = 4 → `"rect" < "square"`) |
| `TopN([]Shape{Square{1}, Square{3}, Circle{1}}, 2)` | `[Square{3} Circle{1}]` |
| `TopN(shapes, 0)` | empty |
| `Kind(Square{1})` / `Kind(Rect{1, 2})` | `"polygon"` |
| `Kind(Circle{1})` | `"round"` |
| `Kind(nil)` | `"none"` |

## Rules
- Areas: `Side²`, `W*H`, `πR²` (`math.Pi`). Value receivers.
- `SortByArea` returns a **new** slice sorted by area ascending; equal areas are ordered by `Name()` ascending; shapes equal on both keep their input order. The input slice must not be modified. Empty or nil input → empty result (`nil` is fine).
- `TopN` returns the `n` shapes with the **largest** area, largest first (ties by `Name()` ascending, then input order). `n <= 0` → empty; `n > len` → all of them. Input unchanged.
- `Kind` uses a type switch: `Square` or `Rect` → `"polygon"`, `Circle` → `"round"`, nil → `"none"`, any other type → `"unknown"`.
