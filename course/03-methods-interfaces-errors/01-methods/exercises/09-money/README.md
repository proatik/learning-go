# Money Type

A `Money` type built on `int64` cents, with arithmetic methods, a `String` method with thousands separators, and a strict parser.

```go
type Money int64 // cents

func New(dollars, cents int64) Money // dollars*100 + cents
func (m Money) Add(o Money) Money
func (m Money) Sub(o Money) Money
func (m Money) Times(n int64) Money
func (m Money) Abs() Money
func (m Money) String() string       // "$1,234.56", "-$0.50"
func Parse(s string) (Money, bool)
```

## Examples

| Code | Result |
|---|---|
| `New(12, 34)` | `Money(1234)` |
| `New(12, 34).Add(New(0, 66))` | `$13.00` |
| `New(5, 0).Sub(New(7, 50))` | `-$2.50` |
| `New(1, 25).Times(3)` | `$3.75` |
| `Money(-50).Abs()` | `$0.50` |
| `Money(7).String()` | `"$0.07"` |
| `Money(123456789).String()` | `"$1,234,567.89"` |
| `Money(-100000).String()` | `"-$1,000.00"` |
| `fmt.Println(New(12, 34))` | `$12.34` |
| `Parse("$1,234.56")` | `123456, true` |
| `Parse("-$0.50")` | `-50, true` |
| `Parse("$5")` | `500, true` |
| `Parse("$1,23.00")` | `0, false` |

## Rules
- `String`: optional `-` **before** the `$`, dollars with a `,` every three digits from the right, then `.` and exactly two cent digits. Zero is `"$0.00"`.
- `Parse` accepts exactly: optional `-`, then `$`, then the dollars, then optionally `.` followed by **exactly two** digits.
  - Dollars are one or more digits, either with no commas at all, or correctly grouped: a first group of 1–3 digits, then groups of exactly 3 separated by single commas (`1,234,567`).
  - Anything else (empty string, missing `$`, spaces, `+`, `$-5`, `$1.2`, `$1.234`, `$1,,000`, `$12,34`, `$.50`) returns `0, false`.
  - `Parse(m.String())` must give back `m`.
- No `float64` anywhere.
