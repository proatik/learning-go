# Shape Catalogue

Build a small shape library: three types that satisfy one interface, print themselves through `fmt.Stringer`, sort through `sort.Interface`, and are counted with a type switch.

```go
type Shape interface {
	Area() float64
	Perimeter() float64
}

type Circle struct{ R float64 }
type Rect struct{ W, H float64 }
type Triangle struct{ A, B, C float64 } // side lengths

func (Circle) Area() float64; func (Circle) Perimeter() float64; func (Circle) String() string
// ...the same three methods on Rect and Triangle (value receivers)

type ByArea []Shape // implements sort.Interface: Len, Less, Swap

func SortByArea(shapes []Shape)
func Largest(shapes []Shape) (Shape, bool)
func Summary(shapes []Shape) string
```

## Examples

| Call | Result |
|---|---|
| `Rect{2, 3}.Area()`, `.Perimeter()` | `6`, `10` |
| `Triangle{3, 4, 5}.Area()` | `6` |
| `Triangle{1, 2, 3}.Area()` | `0` (not a real triangle) |
| `fmt.Sprint(Circle{1.5})` | `"circle(r=1.5)"` |
| `fmt.Sprint(Rect{2, 3})`, `fmt.Sprint(Triangle{3, 4, 5})` | `"rect(2x3)"`, `"triangle(3,4,5)"` |
| `SortByArea([rect(10x10) circle(r=1) rect(1x2)])` | `[rect(1x2) circle(r=1) rect(10x10)]` |
| `Summary([circle(r=1) rect(2x3) square])` | `"circles=1 rects=1 triangles=0 other=1 area=..."` |

## Rules
- Use **value receivers** for every method, so both `Circle{}` and `&Circle{}` are `Shape`s.
- **Circle**: area `πr²`, perimeter `2πr` (use `math.Pi`).
- **Triangle**: area by Heron's formula `√(s(s−a)(s−b)(s−c))` with `s = (a+b+c)/2`. A triangle is **valid** only if every side is `> 0` and each pair of sides sums to **more** than the third; an invalid triangle has area **and** perimeter `0`.
- **String** formats numbers with `%g`: `circle(r=R)`, `rect(WxH)`, `triangle(A,B,C)`.
- **ByArea**: `Less` compares areas (ascending).
- **SortByArea** sorts in place by ascending area; shapes with **equal** areas keep their original order (`sort.Stable`).
- **Largest** returns the first shape with the greatest area, skipping `nil` entries; with no shape at all it returns `nil, false`.
- **Summary** returns exactly `"circles=%d rects=%d triangles=%d other=%d area=%.2f"`:
  - use a **type switch**; `Circle` and `*Circle` both count as circles (same for the others)
  - any other `Shape` (one the package does not know) counts as `other`
  - `nil` entries are skipped; `area` is the sum of all counted areas
