// Package inventory tracks stock and prices orders with volume discounts and
// tax.
package inventory

import (
	"errors"
	"fmt"
	"iter"
	"maps"
	"slices"

	"learngo/course/04-generics-packages-testing/06-module-challenge/exercises/08-inventory/internal/pricing"
)

// Errors returned by Store methods, wrapped with details.
var (
	ErrInvalid      = errors.New("inventory: invalid argument")
	ErrUnknownSKU   = errors.New("inventory: unknown sku")
	ErrInsufficient = errors.New("inventory: insufficient stock")
)

// Item is one product line.
type Item struct {
	SKU   string
	Cents int64 // unit price
	Stock int
}

// config holds the settings chosen with Options.
type config struct {
	taxBP int
	tiers []pricing.Tier
}

// Option configures a Store.
type Option func(*config)

// WithTax charges basisPoints/100 percent tax (825 = 8.25%). Negative values
// are ignored.
func WithTax(basisPoints int) Option {
	return func(c *config) {
		if basisPoints >= 0 {
			c.taxBP = basisPoints
		}
	}
}

// WithVolumeDiscount takes percent off orders of at least minQty units. It
// may be given several times; the best reached tier wins. Tiers with
// minQty < 1 or percent outside 1..100 are ignored.
func WithVolumeDiscount(minQty, percent int) Option {
	return func(c *config) {
		if minQty >= 1 && percent >= 1 && percent <= 100 {
			c.tiers = append(c.tiers, pricing.Tier{MinQty: minQty, Percent: percent})
		}
	}
}

// Store is an inventory. Create one with New.
type Store struct {
	cfg   config
	items map[string]*Item
}

// New returns an empty Store configured by opts (no tax, no discounts by
// default).
func New(opts ...Option) *Store {
	s := &Store{items: make(map[string]*Item)}
	for _, o := range opts {
		o(&s.cfg)
	}
	return s
}

// Add stocks sku at unit price cents. For a known sku it adds to the stock
// and replaces the price.
func (s *Store) Add(sku string, cents int64, stock int) error {
	if sku == "" || cents <= 0 || stock < 0 {
		return fmt.Errorf("%w: Add(%q, %d, %d)", ErrInvalid, sku, cents, stock)
	}
	if it, ok := s.items[sku]; ok {
		it.Cents = cents
		it.Stock += stock
		return nil
	}
	s.items[sku] = &Item{SKU: sku, Cents: cents, Stock: stock}
	return nil
}

// Quote returns the price in cents of qty units of sku: discounted first,
// then taxed. Stock is not changed.
func (s *Store) Quote(sku string, qty int) (int64, error) {
	it, err := s.check(sku, qty)
	if err != nil {
		return 0, err
	}
	subtotal := it.Cents * int64(qty)
	discounted := pricing.Apply(subtotal, pricing.Discount(s.cfg.tiers, qty))
	return pricing.AddTax(discounted, s.cfg.taxBP), nil
}

// Buy quotes qty units of sku and removes them from stock.
func (s *Store) Buy(sku string, qty int) (int64, error) {
	total, err := s.Quote(sku, qty)
	if err != nil {
		return 0, err
	}
	s.items[sku].Stock -= qty
	return total, nil
}

// check validates an order of qty units of sku.
func (s *Store) check(sku string, qty int) (*Item, error) {
	if qty <= 0 {
		return nil, fmt.Errorf("%w: quantity %d", ErrInvalid, qty)
	}
	it, ok := s.items[sku]
	if !ok {
		return nil, fmt.Errorf("%w: %q", ErrUnknownSKU, sku)
	}
	if qty > it.Stock {
		return nil, fmt.Errorf("%w: want %d of %q, have %d", ErrInsufficient, qty, sku, it.Stock)
	}
	return it, nil
}

// Items yields a copy of every item, sorted by SKU.
func (s *Store) Items() iter.Seq[Item] {
	return func(yield func(Item) bool) {
		for _, sku := range slices.Sorted(maps.Keys(s.items)) {
			if !yield(*s.items[sku]) {
				return
			}
		}
	}
}
