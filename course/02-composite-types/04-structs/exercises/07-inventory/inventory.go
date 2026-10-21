package inventory

// Item is a stock entry identified by its SKU.
type Item struct {
	SKU  string
	Name string
	Qty  int
}

// Find returns the item with the given SKU and true, or Item{} and false.
func Find(inv []Item, sku string) (Item, bool) {
	// TODO
	return Item{}, false
}

// AddItem returns a copy of inv with it added: an existing SKU gets its Qty
// increased, a new SKU is appended. inv itself is never modified.
func AddItem(inv []Item, it Item) []Item {
	// TODO
	return nil
}

// RemoveItem returns a copy of inv with qty units of sku removed and true.
// Items whose Qty reaches 0 are dropped. It returns the items unchanged and
// false if sku is missing, qty <= 0, or there isn't enough stock.
// inv itself is never modified.
func RemoveItem(inv []Item, sku string, qty int) ([]Item, bool) {
	// TODO
	return nil, false
}

// TotalQty returns the sum of all quantities.
func TotalQty(inv []Item) int {
	// TODO
	return 0
}
