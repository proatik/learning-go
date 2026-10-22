package inventory_test

import (
	"errors"
	"fmt"
	"slices"
	"testing"

	inventory "learngo/course/04-generics-packages-testing/06-module-challenge/exercises/08-inventory"
)

func newStore(t *testing.T, opts ...inventory.Option) *inventory.Store {
	t.Helper()
	s := inventory.New(opts...)
	if s == nil {
		t.Fatalf("New(...) = nil, want a store")
	}
	return s
}

func add(t *testing.T, s *inventory.Store, sku string, cents int64, stock int) {
	t.Helper()
	if err := s.Add(sku, cents, stock); err != nil {
		t.Fatalf("Add(%q, %d, %d) error = %v, want <nil>", sku, cents, stock, err)
	}
}

func items(s *inventory.Store) []inventory.Item {
	return slices.Collect(s.Items())
}

func TestAddAndItems(t *testing.T) {
	s := newStore(t)
	add(t, s, "pen", 120, 10)
	add(t, s, "ink", 450, 3)
	add(t, s, "pen", 125, 5)
	want := []inventory.Item{{SKU: "ink", Cents: 450, Stock: 3}, {SKU: "pen", Cents: 125, Stock: 15}}
	if got := items(s); !slices.Equal(got, want) {
		t.Errorf("Items() = %+v, want %+v", got, want)
	}

	for it := range s.Items() {
		it.Stock = 0
		break
	}
	if got := items(s); !slices.Equal(got, want) {
		t.Errorf("Items() after changing a yielded copy = %+v, want %+v", got, want)
	}
}

func TestAddInvalid(t *testing.T) {
	s := newStore(t)
	tests := []struct {
		sku   string
		cents int64
		stock int
	}{
		{"", 100, 1},
		{"pen", 0, 1},
		{"pen", -5, 1},
		{"pen", 100, -1},
	}
	for _, tt := range tests {
		if err := s.Add(tt.sku, tt.cents, tt.stock); !errors.Is(err, inventory.ErrInvalid) {
			t.Errorf("Add(%q, %d, %d) error = %v, want %v", tt.sku, tt.cents, tt.stock, err, inventory.ErrInvalid)
		}
	}
	if got := items(s); len(got) != 0 {
		t.Errorf("Items() after invalid Adds = %+v, want none", got)
	}
	add(t, s, "zero-stock", 100, 0)
}

func TestQuote(t *testing.T) {
	tests := []struct {
		name string
		opts []inventory.Option
		qty  int
		want int64
	}{
		{"no options", nil, 3, 360},
		{"tax only", []inventory.Option{inventory.WithTax(825)}, 3, 390},
		{"below the tier", []inventory.Option{inventory.WithVolumeDiscount(10, 15)}, 9, 1080},
		{"at the tier", []inventory.Option{inventory.WithVolumeDiscount(10, 15)}, 10, 1020},
		{"discount then tax", []inventory.Option{inventory.WithTax(825), inventory.WithVolumeDiscount(10, 15)}, 10, 1104},
		{"best tier wins", []inventory.Option{inventory.WithVolumeDiscount(5, 5), inventory.WithVolumeDiscount(20, 20), inventory.WithVolumeDiscount(10, 10)}, 25, 2400},
		{"middle tier", []inventory.Option{inventory.WithVolumeDiscount(5, 5), inventory.WithVolumeDiscount(20, 20), inventory.WithVolumeDiscount(10, 10)}, 12, 1296},
		{"invalid options ignored", []inventory.Option{inventory.WithTax(-100), inventory.WithVolumeDiscount(0, 50), inventory.WithVolumeDiscount(2, 0), inventory.WithVolumeDiscount(2, 101)}, 10, 1200},
		{"last tax wins", []inventory.Option{inventory.WithTax(1000), inventory.WithTax(0)}, 1, 120},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newStore(t, tt.opts...)
			add(t, s, "pen", 120, 100)
			got, err := s.Quote("pen", tt.qty)
			if err != nil || got != tt.want {
				t.Errorf("Quote(%q, %d) = %d, %v, want %d, <nil>", "pen", tt.qty, got, err, tt.want)
			}
			if it := items(s); len(it) != 1 || it[0].Stock != 100 {
				t.Errorf("Items() after Quote = %+v, want stock 100 (Quote must not change stock)", it)
			}
		})
	}
}

func TestQuoteErrors(t *testing.T) {
	s := newStore(t)
	add(t, s, "pen", 120, 3)
	tests := []struct {
		sku     string
		qty     int
		want    error
		wantMsg string
	}{
		{"pen", 0, inventory.ErrInvalid, "inventory: invalid argument: quantity 0"},
		{"pen", -2, inventory.ErrInvalid, "inventory: invalid argument: quantity -2"},
		{"cap", 1, inventory.ErrUnknownSKU, `inventory: unknown sku: "cap"`},
		{"pen", 4, inventory.ErrInsufficient, `inventory: insufficient stock: want 4 of "pen", have 3`},
	}
	for _, tt := range tests {
		got, err := s.Quote(tt.sku, tt.qty)
		if got != 0 || !errors.Is(err, tt.want) {
			t.Errorf("Quote(%q, %d) = %d, %v, want 0, %v", tt.sku, tt.qty, got, err, tt.want)
			continue
		}
		if err.Error() != tt.wantMsg {
			t.Errorf("Quote(%q, %d) error = %q, want %q", tt.sku, tt.qty, err.Error(), tt.wantMsg)
		}
	}
}

func TestBuy(t *testing.T) {
	s := newStore(t, inventory.WithVolumeDiscount(2, 50))
	add(t, s, "mug", 800, 5)
	if got, err := s.Buy("mug", 2); got != 800 || err != nil {
		t.Errorf("Buy(%q, 2) = %d, %v, want 800, <nil>", "mug", got, err)
	}
	if got, err := s.Buy("mug", 1); got != 800 || err != nil {
		t.Errorf("Buy(%q, 1) = %d, %v, want 800, <nil>", "mug", got, err)
	}
	if got, err := s.Buy("mug", 3); got != 0 || !errors.Is(err, inventory.ErrInsufficient) {
		t.Errorf("Buy(%q, 3) with 2 left = %d, %v, want 0, %v", "mug", got, err, inventory.ErrInsufficient)
	}
	if got, err := s.Buy("mug", 2); got != 800 || err != nil {
		t.Errorf("Buy(%q, 2) = %d, %v, want 800, <nil>", "mug", got, err)
	}
	if want := []inventory.Item{{SKU: "mug", Cents: 800, Stock: 0}}; !slices.Equal(items(s), want) {
		t.Errorf("Items() = %+v, want %+v", items(s), want)
	}
}

func Example() {
	shop := inventory.New(
		inventory.WithTax(825),
		inventory.WithVolumeDiscount(10, 15),
	)
	shop.Add("pen", 120, 50)
	shop.Add("ink", 450, 2)

	for _, order := range []struct {
		sku string
		qty int
	}{{"pen", 10}, {"ink", 5}} {
		cents, err := shop.Buy(order.sku, order.qty)
		if err != nil {
			fmt.Println("error:", err)
			continue
		}
		fmt.Printf("%d × %s = %d cents\n", order.qty, order.sku, cents)
	}
	for it := range shop.Items() {
		fmt.Println(it.SKU, it.Stock)
	}
	// Output:
	// 10 × pen = 1104 cents
	// error: inventory: insufficient stock: want 5 of "ink", have 2
	// ink 2
	// pen 40
}
