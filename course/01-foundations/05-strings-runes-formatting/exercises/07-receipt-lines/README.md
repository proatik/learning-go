# Receipt Lines

Format the lines of a shop receipt with `fmt.Sprintf` width and precision. Implement:

```go
func Header() string
func Separator() string
func Line(item string, qty int, unit float64) string
func Total(total float64) string
```

Every line is exactly **37 runes** wide:

```
Item         Qty     Price      Total
-------------------------------------
Coffee         2 x    3.50 =     7.00
Café           1 x    4.25 =     4.25
Croissant au   3 x    2.10 =     6.30
TOTAL                           17.55
```

| Function | Layout |
|---|---|
| `Header()` | `"Item"` left-aligned in 12, space, `"Qty"` right-aligned in 3, 3 spaces, `"Price"` right in 7, 3 spaces, `"Total"` right in 8 |
| `Separator()` | 37 dashes `-` |
| `Line(item, qty, unit)` | item left-aligned in 12 and **cut to 12 runes**, space, qty right in 3, `" x "`, unit right in 7 with 2 decimals, `" = "`, `qty*unit` right in 8 with 2 decimals |
| `Total(total)` | `"TOTAL"` left-aligned in 29, then total right in 8 with 2 decimals |

## Rules
- Use one `fmt.Sprintf` per function (and `strings.Repeat` for `Separator`).
- `%-12.12s` = left-aligned, width 12, at most 12 runes (precision truncates strings).
- Remember to convert: `float64(qty) * unit`.
