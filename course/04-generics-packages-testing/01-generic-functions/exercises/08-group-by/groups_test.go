package groups_test

import (
	"slices"
	"testing"

	grp "learngo/course/04-generics-packages-testing/01-generic-functions/exercises/08-group-by"
)

type Order struct {
	ID     int
	City   string
	Amount int
}

var orders = []Order{
	{1, "Oslo", 50},
	{2, "Rome", 80},
	{3, "Oslo", 20},
	{4, "Lima", 80},
	{5, "Oslo", 70},
}

func ids(os []Order) []int {
	out := make([]int, 0, len(os))
	for _, o := range os {
		out = append(out, o.ID)
	}
	return out
}

func isEven(n int) bool { return n%2 == 0 }

func TestGroupByCity(t *testing.T) {
	got := grp.GroupBy(orders, func(o Order) string { return o.City })
	if got == nil {
		t.Fatalf("GroupBy(orders, city) = nil, want a map with 3 keys")
	}
	want := map[string][]int{"Oslo": {1, 3, 5}, "Rome": {2}, "Lima": {4}}
	if len(got) != len(want) {
		t.Errorf("GroupBy(orders, city) has %d keys, want %d", len(got), len(want))
	}
	for city, wantIDs := range want {
		if gotIDs := ids(got[city]); !slices.Equal(gotIDs, wantIDs) {
			t.Errorf("GroupBy(orders, city)[%q] = IDs %v, want %v", city, gotIDs, wantIDs)
		}
	}
}

func TestGroupByOtherKeys(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5}
	got := grp.GroupBy(nums, isEven)
	if g, w := got[true], []int{2, 4}; !slices.Equal(g, w) {
		t.Errorf("GroupBy(%v, isEven)[true] = %v, want %v", nums, g, w)
	}
	if g, w := got[false], []int{1, 3, 5}; !slices.Equal(g, w) {
		t.Errorf("GroupBy(%v, isEven)[false] = %v, want %v", nums, g, w)
	}

	words := []string{"go", "c", "zig", "js", "", "ada"}
	byLen := grp.GroupBy(words, func(s string) int { return len(s) })
	want := map[int][]string{0: {""}, 1: {"c"}, 2: {"go", "js"}, 3: {"zig", "ada"}}
	if len(byLen) != len(want) {
		t.Errorf("GroupBy(%q, len) has %d keys, want %d", words, len(byLen), len(want))
	}
	for n, w := range want {
		if g := byLen[n]; !slices.Equal(g, w) {
			t.Errorf("GroupBy(%q, len)[%d] = %q, want %q", words, n, g, w)
		}
	}
}

func TestGroupByEmpty(t *testing.T) {
	for _, in := range [][]int{nil, {}} {
		got := grp.GroupBy(in, isEven)
		if got == nil || len(got) != 0 {
			t.Errorf("GroupBy(%v, isEven) = %v (nil=%t), want an empty non-nil map", in, got, got == nil)
		}
	}
}

func TestPartition(t *testing.T) {
	tests := []struct {
		in, yes, no []int
	}{
		{[]int{1, 2, 3, 4}, []int{2, 4}, []int{1, 3}},
		{[]int{4, 3, 2, 1}, []int{4, 2}, []int{3, 1}},
		{[]int{2, 4}, []int{2, 4}, []int{}},
		{[]int{1, 3}, []int{}, []int{1, 3}},
		{[]int{0, -1, -2}, []int{0, -2}, []int{-1}},
		{nil, []int{}, []int{}},
	}
	for _, tt := range tests {
		yes, no := grp.Partition(tt.in, isEven)
		if !slices.Equal(yes, tt.yes) || !slices.Equal(no, tt.no) {
			t.Errorf("Partition(%v, isEven) = %v, %v, want %v, %v", tt.in, yes, no, tt.yes, tt.no)
		}
	}

	big, small := grp.Partition(orders, func(o Order) bool { return o.Amount >= 70 })
	if g, w := ids(big), []int{2, 4, 5}; !slices.Equal(g, w) {
		t.Errorf("Partition(orders, amount>=70) yes = IDs %v, want %v", g, w)
	}
	if g, w := ids(small), []int{1, 3}; !slices.Equal(g, w) {
		t.Errorf("Partition(orders, amount>=70) no = IDs %v, want %v", g, w)
	}
}

func TestMaxBy(t *testing.T) {
	amount := func(o Order) int { return o.Amount }
	got, ok := grp.MaxBy(orders, amount)
	if !ok || got.ID != 2 {
		t.Errorf("MaxBy(orders, amount) = %v, %t, want %v, %t (first of the tied 80s)", got, ok, orders[1], true)
	}

	city := func(o Order) string { return o.City }
	got, ok = grp.MaxBy(orders, city)
	if !ok || got.ID != 2 {
		t.Errorf("MaxBy(orders, city) = %v, %t, want %v, %t", got, ok, orders[1], true)
	}

	negAmount := func(o Order) int { return -o.Amount }
	got, ok = grp.MaxBy(orders, negAmount)
	if !ok || got.ID != 3 {
		t.Errorf("MaxBy(orders, -amount) = %v, %t, want %v, %t", got, ok, orders[2], true)
	}

	words := []string{"go", "rust", "zig", "java"}
	w, ok := grp.MaxBy(words, func(s string) int { return len(s) })
	if !ok || w != "rust" {
		t.Errorf("MaxBy(%q, len) = %q, %t, want %q, %t", words, w, ok, "rust", true)
	}

	ratio := func(n int) float64 { return 1 / float64(n) }
	n, ok := grp.MaxBy([]int{4, 2, 8}, ratio)
	if !ok || n != 2 {
		t.Errorf("MaxBy(%v, 1/n) = %d, %t, want %d, %t", []int{4, 2, 8}, n, ok, 2, true)
	}

	neg := []int{-9, -3, -5}
	n, ok = grp.MaxBy(neg, func(n int) int { return n })
	if !ok || n != -3 {
		t.Errorf("MaxBy(%v, identity) = %d, %t, want %d, %t", neg, n, ok, -3, true)
	}
}

func TestMaxByEmpty(t *testing.T) {
	amount := func(o Order) int { return o.Amount }
	if got, ok := grp.MaxBy([]Order{}, amount); ok || got != (Order{}) {
		t.Errorf("MaxBy([], amount) = %v, %t, want %v, %t", got, ok, Order{}, false)
	}
	if got, ok := grp.MaxBy(nil, amount); ok || got != (Order{}) {
		t.Errorf("MaxBy(nil, amount) = %v, %t, want %v, %t", got, ok, Order{}, false)
	}
}

func TestCallsOncePerElement(t *testing.T) {
	calls := 0
	grp.MaxBy(orders, func(o Order) int { calls++; return o.Amount })
	if calls != len(orders) {
		t.Errorf("MaxBy(orders, key) called key %d times, want %d", calls, len(orders))
	}
	calls = 0
	grp.GroupBy(orders, func(o Order) string { calls++; return o.City })
	if calls != len(orders) {
		t.Errorf("GroupBy(orders, key) called key %d times, want %d", calls, len(orders))
	}
	calls = 0
	grp.Partition(orders, func(o Order) bool { calls++; return true })
	if calls != len(orders) {
		t.Errorf("Partition(orders, pred) called pred %d times, want %d", calls, len(orders))
	}
}
