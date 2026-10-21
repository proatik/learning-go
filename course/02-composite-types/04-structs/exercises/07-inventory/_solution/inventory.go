package inventory

import "slices"

// Item is a stock entry identified by its SKU.
type Item struct {
	SKU  string
	Name string
	Qty  int
}

// index returns the position of sku in inv, or -1.
func index(inv []Item, sku string) int {
	return slices.IndexFunc(inv, func(it Item) bool { return it.SKU == sku })
}

// Find returns the item with the given SKU and true, or Item{} and false.
func Find(inv []Item, sku string) (Item, bool) {
	if i := index(inv, sku); i >= 0 {
		return inv[i], true // returns a copy of the element
	}
	return Item{}, false
}

// AddItem returns a copy of inv with it added: an existing SKU gets its Qty
// increased, a new SKU is appended. inv itself is never modified.
func AddItem(inv []Item, it Item) []Item {
	out := slices.Clone(inv) // our own backing array: safe to modify and append
	if it.Qty <= 0 {
		return out
	}
	if i := index(out, it.SKU); i >= 0 {
		out[i].Qty += it.Qty // modify through the index, not a range copy
		return out
	}
	return append(out, it)
}

// RemoveItem returns a copy of inv with qty units of sku removed and true.
// Items whose Qty reaches 0 are dropped. It returns the items unchanged and
// false if sku is missing, qty <= 0, or there isn't enough stock.
// inv itself is never modified.
func RemoveItem(inv []Item, sku string, qty int) ([]Item, bool) {
	out := slices.Clone(inv)
	i := index(out, sku)
	if i < 0 || qty <= 0 || qty > out[i].Qty {
		return out, false
	}
	out[i].Qty -= qty
	if out[i].Qty == 0 {
		out = slices.Delete(out, i, i+1) // keeps the order of the others
	}
	return out, true
}

// TotalQty returns the sum of all quantities.
func TotalQty(inv []Item) int {
	total := 0
	for _, it := range inv {
		total += it.Qty
	}
	return total
}
