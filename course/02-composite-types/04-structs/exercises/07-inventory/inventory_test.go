package inventory

import (
	"slices"
	"testing"
)

func stock() []Item {
	return []Item{
		{SKU: "A1", Name: "apple", Qty: 5},
		{SKU: "B2", Name: "bread", Qty: 2},
		{SKU: "C3", Name: "cheese", Qty: 1},
	}
}

func TestFind(t *testing.T) {
	tests := []struct {
		sku    string
		want   Item
		wantOK bool
	}{
		{"A1", Item{"A1", "apple", 5}, true},
		{"C3", Item{"C3", "cheese", 1}, true},
		{"Z9", Item{}, false},
		{"", Item{}, false},
		{"a1", Item{}, false},
	}
	for _, tt := range tests {
		got, ok := Find(stock(), tt.sku)
		if got != tt.want || ok != tt.wantOK {
			t.Errorf("Find(stock, %q) = %+v, %v; want %+v, %v", tt.sku, got, ok, tt.want, tt.wantOK)
		}
	}
	if got, ok := Find(nil, "A1"); ok || got != (Item{}) {
		t.Errorf("Find(nil, \"A1\") = %+v, %v; want {}, false", got, ok)
	}
}

func TestAddItem(t *testing.T) {
	tests := []struct {
		it   Item
		want []Item
	}{
		{Item{"A1", "", 3}, []Item{{"A1", "apple", 8}, {"B2", "bread", 2}, {"C3", "cheese", 1}}},
		{Item{"C3", "CHEESE!", 10}, []Item{{"A1", "apple", 5}, {"B2", "bread", 2}, {"C3", "cheese", 11}}},
		{Item{"D4", "dates", 7}, []Item{{"A1", "apple", 5}, {"B2", "bread", 2}, {"C3", "cheese", 1}, {"D4", "dates", 7}}},
		{Item{"A1", "apple", 0}, stock()},
		{Item{"E5", "eggs", -2}, stock()},
	}
	for _, tt := range tests {
		inv := stock()
		got := AddItem(inv, tt.it)
		if !slices.Equal(got, tt.want) {
			t.Errorf("AddItem(stock, %+v) = %+v, want %+v", tt.it, got, tt.want)
		}
		if !slices.Equal(inv, stock()) {
			t.Errorf("AddItem(stock, %+v) modified its input: %+v", tt.it, inv)
		}
	}

	got := AddItem(nil, Item{"X", "x", 1})
	if want := []Item{{"X", "x", 1}}; !slices.Equal(got, want) {
		t.Errorf("AddItem(nil, {X x 1}) = %+v, want %+v", got, want)
	}
}

func TestAddItemSpareCapacity(t *testing.T) {
	// inv has room for more elements behind its length.
	inv := make([]Item, 0, 10)
	inv = append(inv, Item{"A1", "apple", 5})

	a := AddItem(inv, Item{"B2", "bread", 1})
	b := AddItem(inv, Item{"C3", "cheese", 1})
	if len(a) != 2 || len(b) != 2 {
		t.Fatalf("AddItem results have lengths %d and %d, want 2 and 2", len(a), len(b))
	}
	if a[1].SKU != "B2" {
		t.Errorf("first AddItem result changed after a second AddItem on the same inv: %+v; results must not share memory with inv", a)
	}
	if b[1].SKU != "C3" {
		t.Errorf("second AddItem result = %+v, want C3 at index 1", b)
	}
	if len(inv) != 1 || inv[0] != (Item{"A1", "apple", 5}) {
		t.Errorf("inv was modified: %+v", inv)
	}
}

func TestRemoveItem(t *testing.T) {
	tests := []struct {
		sku    string
		qty    int
		want   []Item
		wantOK bool
	}{
		{"A1", 2, []Item{{"A1", "apple", 3}, {"B2", "bread", 2}, {"C3", "cheese", 1}}, true},
		{"B2", 2, []Item{{"A1", "apple", 5}, {"C3", "cheese", 1}}, true},
		{"A1", 5, []Item{{"B2", "bread", 2}, {"C3", "cheese", 1}}, true},
		{"C3", 1, []Item{{"A1", "apple", 5}, {"B2", "bread", 2}}, true},
		{"B2", 3, stock(), false},
		{"Z9", 1, stock(), false},
		{"A1", 0, stock(), false},
		{"A1", -1, stock(), false},
	}
	for _, tt := range tests {
		inv := stock()
		got, ok := RemoveItem(inv, tt.sku, tt.qty)
		if ok != tt.wantOK || !slices.Equal(got, tt.want) {
			t.Errorf("RemoveItem(stock, %q, %d) = %+v, %v; want %+v, %v", tt.sku, tt.qty, got, ok, tt.want, tt.wantOK)
		}
		if !slices.Equal(inv, stock()) {
			t.Errorf("RemoveItem(stock, %q, %d) modified its input: %+v", tt.sku, tt.qty, inv)
		}
	}
}

func TestChainedOperations(t *testing.T) {
	var inv []Item
	inv = AddItem(inv, Item{"P", "pen", 3})
	inv = AddItem(inv, Item{"Q", "quill", 1})
	inv = AddItem(inv, Item{"P", "", 2})
	inv, ok := RemoveItem(inv, "Q", 1)
	if !ok {
		t.Fatalf("RemoveItem(%+v, \"Q\", 1) returned false, want true", inv)
	}
	want := []Item{{"P", "pen", 5}}
	if !slices.Equal(inv, want) {
		t.Errorf("after add P3, add Q1, add P2, remove Q1: %+v, want %+v", inv, want)
	}
	if got := TotalQty(inv); got != 5 {
		t.Errorf("TotalQty(%+v) = %d, want 5", inv, got)
	}
}

func TestTotalQty(t *testing.T) {
	tests := []struct {
		inv  []Item
		want int
	}{
		{stock(), 8},
		{nil, 0},
		{[]Item{{"A", "a", 100}}, 100},
	}
	for _, tt := range tests {
		if got := TotalQty(tt.inv); got != tt.want {
			t.Errorf("TotalQty(%+v) = %d, want %d", tt.inv, got, tt.want)
		}
	}
}
