# Inventory with Value Semantics

An inventory is a `[]Item`. Every function returns an **updated copy** and leaves the caller's slice exactly as it was, like a struct value.

```go
type Item struct {
	SKU  string
	Name string
	Qty  int
}

func Find(inv []Item, sku string) (Item, bool)
func AddItem(inv []Item, it Item) []Item
func RemoveItem(inv []Item, sku string, qty int) ([]Item, bool)
func TotalQty(inv []Item) int
```

`inv := []Item{{"A1", "apple", 5}, {"B2", "bread", 2}}`

| Call | Result |
|---|---|
| `Find(inv, "B2")` | `{B2 bread 2}, true` |
| `Find(inv, "Z9")` | `Item{}, false` |
| `AddItem(inv, Item{"A1", "", 3})` | `[{A1 apple 8} {B2 bread 2}]` |
| `AddItem(inv, Item{"C3", "cheese", 1})` | `[{A1 apple 5} {B2 bread 2} {C3 cheese 1}]` |
| `RemoveItem(inv, "A1", 2)` | `[{A1 apple 3} {B2 bread 2}], true` |
| `RemoveItem(inv, "B2", 2)` | `[{A1 apple 5}], true` (reached 0 → removed) |
| `RemoveItem(inv, "B2", 3)` | `inv unchanged, false` (not enough) |
| `TotalQty(inv)` | `7` |

## Rules
- **Never modify `inv`**: not its elements, not its order, and not the spare capacity behind it. Two calls on the same `inv` must give independent results.
- `AddItem`: if the SKU exists, increase its `Qty` and keep the existing `Name`; otherwise append `it` at the end. If `it.Qty <= 0`, return the items unchanged.
- `RemoveItem`: returns `false` and the items unchanged if the SKU is missing, `qty <= 0`, or `qty` is more than in stock. An item whose `Qty` reaches 0 is removed; other items keep their order.
- Watch out: `append(inv, it)` may write into `inv`'s backing array when it has spare capacity. Start from `slices.Clone(inv)`.
