# Shapes

Implement one interface with three types, then write functions that work on **any** `Shape`.

```go
type Shape interface {
	Area() float64
	Perimeter() float64
}

type Rect struct{ W, H float64 }
type Circle struct{ R float64 }
type Triangle struct{ A, B, C float64 } // side lengths

func (r Rect) Area() float64
func (r Rect) Perimeter() float64
func (c Circle) Area() float64
func (c Circle) Perimeter() float64
func (t Triangle) Area() float64
func (t Triangle) Perimeter() float64

func TotalArea(shapes []Shape) float64
func Largest(shapes []Shape) (Shape, bool)
```

## Examples

| Call | Result |
|---|---|
| `Rect{3, 4}.Area()` / `.Perimeter()` | `12` / `14` |
| `Circle{1}.Area()` / `.Perimeter()` | `math.Pi` / `2*math.Pi` |
| `Triangle{3, 4, 5}.Area()` / `.Perimeter()` | `6` / `12` |
| `TotalArea([]Shape{Rect{3, 4}, Triangle{3, 4, 5}})` | `18` |
| `TotalArea(nil)` | `0` |
| `Largest([]Shape{Rect{1, 1}, Circle{1}, Rect{1, 4}})` | `Rect{1, 4}, true` |
| `Largest(nil)` | `nil, false` |

## Rules
- All methods use **value receivers**.
- Triangle area uses **Heron's formula**: `s = (A+B+C)/2`, `area = √(s(s−A)(s−B)(s−C))`. You may assume valid sides.
- `Largest` returns the shape with the biggest `Area()`; on a tie the **first** one wins. Empty or nil input → `nil, false`.
- Use `math.Pi` and `math.Sqrt`.
