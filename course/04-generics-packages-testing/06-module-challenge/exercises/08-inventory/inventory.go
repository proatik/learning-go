// Package inventory tracks stock and prices orders with volume discounts and
// tax.
package inventory

import (
	"errors"
	"iter"

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
		// TODO
	}
}

// WithVolumeDiscount takes percent off orders of at least minQty units. It
// may be given several times; the best reached tier wins. Tiers with
// minQty < 1 or percent outside 1..100 are ignored.
func WithVolumeDiscount(minQty, percent int) Option {
	return func(c *config) {
		// TODO
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
	// TODO
	return nil
}

// Add stocks sku at unit price cents. For a known sku it adds to the stock
// and replaces the price.
func (s *Store) Add(sku string, cents int64, stock int) error {
	// TODO
	return nil
}

// Quote returns the price in cents of qty units of sku: discounted first,
// then taxed. Stock is not changed.
func (s *Store) Quote(sku string, qty int) (int64, error) {
	// TODO: use pricing.Discount, pricing.Apply and pricing.AddTax.
	return 0, nil
}

// Buy quotes qty units of sku and removes them from stock.
func (s *Store) Buy(sku string, qty int) (int64, error) {
	// TODO
	return 0, nil
}

// Items yields a copy of every item, sorted by SKU.
func (s *Store) Items() iter.Seq[Item] {
	return func(yield func(Item) bool) {
		// TODO
	}
}
