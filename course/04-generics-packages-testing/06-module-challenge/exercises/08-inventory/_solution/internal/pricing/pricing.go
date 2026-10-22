// Package pricing holds the money arithmetic used by package inventory.
// It is internal: only code inside the inventory tree can import it.
package pricing

// Tier is a volume discount: buying at least MinQty units takes Percent off.
type Tier struct {
	MinQty  int
	Percent int
}

// Discount returns the highest Percent among the tiers that qty reaches, or
// 0 if it reaches none.
func Discount(tiers []Tier, qty int) int {
	best := 0
	for _, t := range tiers {
		if qty >= t.MinQty && t.Percent > best {
			best = t.Percent
		}
	}
	return best
}

// Apply takes percent off cents, rounding half a cent up.
func Apply(cents int64, percent int) int64 {
	return (cents*int64(100-percent) + 50) / 100
}

// AddTax adds basisPoints/100 percent tax to cents, rounding half a cent up.
func AddTax(cents int64, basisPoints int) int64 {
	return cents + (cents*int64(basisPoints)+5000)/10000
}
