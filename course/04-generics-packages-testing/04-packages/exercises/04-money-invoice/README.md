# Money and Invoices

Build a `money` package whose `Amount` can never be invalid, then an `invoice` package that uses it.

```go
// package money  (money/money.go)
var (
	ErrNegative = errors.New("money: negative amount")
	ErrCents    = errors.New("money: cents must be between 0 and 99")
)

type Amount struct{ cents int64 } // unexported: only this package can set it

func New(units, cents int64) (Amount, error)
func FromCents(cents int64) (Amount, error)
func (a Amount) Cents() int64
func (a Amount) Add(b Amount) Amount
func (a Amount) Times(n int) Amount
func (a Amount) String() string

// package invoice  (invoice.go)
var (
	ErrNoItem = errors.New("invoice: item name is empty")
	ErrQty    = errors.New("invoice: quantity must be at least 1")
)

type Line struct {
	Item  string
	Price money.Amount
	Qty   int
}

type Invoice struct{ lines []Line } // the zero value is an empty invoice

func (l Line) Subtotal() money.Amount
func (inv *Invoice) Add(item string, price money.Amount, qty int) error
func (inv *Invoice) Lines() []Line
func (inv *Invoice) Total() money.Amount
func (inv *Invoice) String() string
```

## Examples

| Call | Result |
|---|---|
| `money.New(12, 5)` | `12.05`, `nil` |
| `money.New(-1, 0)` | error `ErrNegative` |
| `money.New(3, 100)` | error `ErrCents` |
| `money.FromCents(1999)` | `19.99`, `nil` |
| `money.Amount{}.String()` | `"0.00"` |
| `p.Times(3)` with `p` = `1.50` | `4.50` |
| `Line{"Pen", 1.50, 2}.Subtotal()` | `3.00` |

`var inv invoice.Invoice; inv.Add("Pen", 1.50, 2); inv.Add("Pad", 0.99, 1); inv.String()`:

```text
2 x Pen @ 1.50 = 3.00
1 x Pad @ 0.99 = 0.99
Total: 3.99
```

## Rules
- `New`: `units < 0` or `cents < 0` → `ErrNegative`. Otherwise `cents > 99` → `ErrCents`. On error, return `Amount{}`.
- `FromCents`: negative → `Amount{}` and `ErrNegative`.
- `Times(n)` with `n < 0` returns `Amount{}`, so an `Amount` is never negative.
- `String` is `"<units>.<cents with 2 digits>"`, e.g. `"0.05"`, `"100.00"`.
- `Add` checks, in order:
  - empty `item` → `ErrNoItem`;
  - `qty < 1` → `ErrQty`.
  - Nothing is added on error.
- `Lines` returns a **copy**, so callers can't modify the invoice. `nil` or empty is fine for an empty invoice.
- `String`: one line per item, `"<qty> x <item> @ <price> = <subtotal>"`, then `"Total: <total>"`, joined by `"\n"` with no trailing newline. An empty invoice is just `"Total: 0.00"`.
