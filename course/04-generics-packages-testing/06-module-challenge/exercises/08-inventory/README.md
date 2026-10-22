# Inventory with Internal Pricing

Design a small `inventory` package: a constructor with **functional options**, sentinel errors, a sorted item iterator — and all money arithmetic hidden in an **internal** `pricing` package.

```go
// package pricing — .../08-inventory/internal/pricing (importable only inside 08-inventory/)
type Tier struct {
	MinQty  int
	Percent int
}

func Discount(tiers []Tier, qty int) int
func Apply(cents int64, percent int) int64
func AddTax(cents int64, basisPoints int) int64

// package inventory — the exercise root
var (
	ErrInvalid      = errors.New("inventory: invalid argument")
	ErrUnknownSKU   = errors.New("inventory: unknown sku")
	ErrInsufficient = errors.New("inventory: insufficient stock")
)

type Item struct {
	SKU   string
	Cents int64 // unit price
	Stock int
}

type config struct {
	taxBP int
	tiers []pricing.Tier
}

type Option func(*config)

func WithTax(basisPoints int) Option
func WithVolumeDiscount(minQty, percent int) Option
func New(opts ...Option) *Store

func (s *Store) Add(sku string, cents int64, stock int) error
func (s *Store) Quote(sku string, qty int) (int64, error)
func (s *Store) Buy(sku string, qty int) (int64, error)
func (s *Store) Items() iter.Seq[Item]
```

## Examples

| Code | Result |
|---|---|
| `pricing.Discount([{10 5} {50 15} {20 10}], 20)` | `10` — best tier reached |
| `pricing.Apply(999, 15)` | `849` (849.15) |
| `pricing.AddTax(1000, 825)` | `1083` (1082.5 rounds up) |
| `New(WithTax(825), WithVolumeDiscount(10, 15))`, pen at 120¢: `Quote("pen", 10)` | `1104` = 1200 → −15% → 1020 → +8.25% tax |
| `Quote("pen", 0)` | `0`, `"inventory: invalid argument: quantity 0"` |
| `Quote("cap", 1)` | `0`, `inventory: unknown sku: "cap"` |
| `Quote("pen", 4)` with 3 in stock | `0`, `inventory: insufficient stock: want 4 of "pen", have 3` |

## Rules
- **pricing.Discount** returns the highest `Percent` among tiers with `qty >= MinQty`, else `0`.
- **pricing.Apply** = `(cents*(100-percent) + 50) / 100`; **pricing.AddTax** = `cents + (cents*bp + 5000) / 10000` (half a cent rounds up).
- **New** starts from "no tax, no tiers", then applies each option in order. **WithTax** with a negative value and **WithVolumeDiscount** with `minQty < 1` or `percent` outside `1..100` change nothing. The last valid `WithTax` wins; every valid tier is kept.
- **Add** rejects an empty `sku`, `cents <= 0` or `stock < 0` with an error wrapping `ErrInvalid`. A known SKU gets the **new price** and **extra** stock.
- **Quote** checks in this order: `qty <= 0` → `fmt.Errorf("%w: quantity %d", ErrInvalid, qty)`; unknown → `fmt.Errorf("%w: %q", ErrUnknownSKU, sku)`; too few → `fmt.Errorf("%w: want %d of %q, have %d", ErrInsufficient, ...)`. The price is `Cents*qty`, discounted, **then** taxed. Quote never changes stock.
- **Buy** is `Quote` plus removing `qty` from stock — nothing changes on error.
- **Items** yields **copies** sorted by SKU; changing one must not change the store.
- Callers outside the tree can't import `internal/pricing`, so `Store`'s API must not require it.
