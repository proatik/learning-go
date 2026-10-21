# Rectangle Methods

Give `Rect` read-only methods with **value** receivers and mutating methods with **pointer** receivers.

```go
type Rect struct {
	W, H float64
}

func (r Rect) Area() float64
func (r Rect) Perimeter() float64
func (r Rect) IsSquare() bool
func (r *Rect) Scale(k float64)
func (r *Rect) Grow(dw, dh float64)
```

## Examples

| Code | Result |
|---|---|
| `Rect{3, 4}.Area()` | `12` |
| `Rect{3, 4}.Perimeter()` | `14` |
| `Rect{2, 2}.IsSquare()` | `true` |
| `r := Rect{2, 3}; r.Scale(10)` | `r` is `{20 30}` |
| `r := Rect{2, 3}; r.Grow(1, -1)` | `r` is `{3 2}` |
| `r := Rect{2, 3}; r.Grow(-5, 0)` | `r` is `{0 3}` |

## Rules
- `Scale` multiplies both sides by `k` **in place**.
- `Grow` adds `dw` to `W` and `dh` to `H` **in place**; a side that would become negative becomes `0`.
- `IsSquare` is true when `W == H` (including `{0 0}`).
- `Area` and `Perimeter` must not change the receiver.
