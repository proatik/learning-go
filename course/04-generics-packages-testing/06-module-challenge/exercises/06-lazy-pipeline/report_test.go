package report

import (
	"iter"
	"slices"
	"testing"
)

var orders = []Order{
	{ID: 1, Customer: "ana", Cents: 1250},
	{ID: 2, Customer: "bo", Cents: 99},
	{ID: 3, Customer: "cy", Cents: 50000, Refunded: true},
	{ID: 4, Customer: "di", Cents: 2000},
	{ID: 5, Customer: "ed", Cents: 1000},
	{ID: 6, Customer: "fay", Cents: 7},
	{ID: 7, Customer: "gus", Cents: 3005},
}

// source yields list and counts how many were read in *pulled.
func source(list []Order, pulled *int) iter.Seq[Order] {
	return func(yield func(Order) bool) {
		for _, o := range list {
			*pulled++
			if !yield(o) {
				return
			}
		}
	}
}

func TestRevenue(t *testing.T) {
	if got := Revenue(slices.Values(orders)); got != 7361 {
		t.Errorf("Revenue(orders) = %d, want 7361 (refunds excluded)", got)
	}
	if got := Revenue(slices.Values([]Order{})); got != 0 {
		t.Errorf("Revenue(none) = %d, want 0", got)
	}
}

func TestLargeOrderIDs(t *testing.T) {
	tests := []struct {
		min        int64
		n          int
		want       []int
		wantPulled int
	}{
		{1000, 2, []int{1, 4}, 4},
		{1000, 3, []int{1, 4, 5}, 5},
		{1000, 10, []int{1, 4, 5, 7}, 7},
		{2000, 1, []int{4}, 4},
		{100000, 1, []int{}, 7},
		{0, 0, []int{}, 0},
	}
	for _, tt := range tests {
		pulled := 0
		got := LargeOrderIDs(source(orders, &pulled), tt.min, tt.n)
		if got == nil || !slices.Equal(got, tt.want) {
			t.Errorf("LargeOrderIDs(orders, %d, %d) = %#v, want %v", tt.min, tt.n, got, tt.want)
		}
		if pulled != tt.wantPulled {
			t.Errorf("LargeOrderIDs(orders, %d, %d) read %d orders, want %d (stop once you have n)", tt.min, tt.n, pulled, tt.wantPulled)
		}
	}
}

func TestLine(t *testing.T) {
	tests := []struct {
		o    Order
		want string
	}{
		{Order{ID: 1, Customer: "ana", Cents: 1250}, "#1 ana $12.50"},
		{Order{ID: 42, Customer: "bo", Cents: 7}, "#42 bo $0.07"},
		{Order{ID: 3, Customer: "Zoë", Cents: 0}, "#3 Zoë $0.00"},
		{Order{ID: 9, Customer: "cy", Cents: 100000, Refunded: true}, "#9 cy $1000.00 (refunded)"},
	}
	for _, tt := range tests {
		if got := Line(tt.o); got != tt.want {
			t.Errorf("Line(%+v) = %q, want %q", tt.o, got, tt.want)
		}
	}
}

func TestPages(t *testing.T) {
	got := slices.Collect(Pages(slices.Values(orders[:5]), 2))
	want := [][]string{
		{"#1 ana $12.50", "#2 bo $0.99"},
		{"#3 cy $500.00 (refunded)", "#4 di $20.00"},
		{"#5 ed $10.00"},
	}
	if !slices.EqualFunc(got, want, slices.Equal) {
		t.Errorf("Pages(orders[:5], 2) = %q, want %q", got, want)
	}

	pulled := 0
	for page := range Pages(source(orders, &pulled), 3) {
		_ = page
		break
	}
	if pulled != 3 {
		t.Errorf("reading the first page of 3 read %d orders, want 3", pulled)
	}
}
