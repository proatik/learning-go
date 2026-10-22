package pricing

import "testing"

func TestDiscount(t *testing.T) {
	tiers := []Tier{{MinQty: 10, Percent: 5}, {MinQty: 50, Percent: 15}, {MinQty: 20, Percent: 10}}
	tests := []struct {
		tiers []Tier
		qty   int
		want  int
	}{
		{tiers, 1, 0},
		{tiers, 9, 0},
		{tiers, 10, 5},
		{tiers, 19, 5},
		{tiers, 20, 10},
		{tiers, 49, 10},
		{tiers, 50, 15},
		{tiers, 1000, 15},
		{nil, 100, 0},
		{[]Tier{{MinQty: 5, Percent: 30}, {MinQty: 10, Percent: 20}}, 12, 30},
	}
	for _, tt := range tests {
		if got := Discount(tt.tiers, tt.qty); got != tt.want {
			t.Errorf("Discount(%v, %d) = %d, want %d", tt.tiers, tt.qty, got, tt.want)
		}
	}
}

func TestApply(t *testing.T) {
	tests := []struct {
		cents   int64
		percent int
		want    int64
	}{
		{1000, 0, 1000},
		{1000, 10, 900},
		{1000, 100, 0},
		{999, 15, 849},   // 849.15
		{1, 50, 1},       // 0.5 rounds up
		{3, 50, 2},       // 1.5 rounds up
		{1999, 33, 1339}, // 1339.33
	}
	for _, tt := range tests {
		if got := Apply(tt.cents, tt.percent); got != tt.want {
			t.Errorf("Apply(%d, %d) = %d, want %d", tt.cents, tt.percent, got, tt.want)
		}
	}
}

func TestAddTax(t *testing.T) {
	tests := []struct {
		cents int64
		bp    int
		want  int64
	}{
		{1000, 0, 1000},
		{1000, 825, 1083}, // 1082.5 rounds up
		{1020, 825, 1104}, // 1104.15
		{1, 4999, 1},      // 1.4999
		{1, 5000, 2},      // 1.5 rounds up
		{20000, 2000, 24000},
	}
	for _, tt := range tests {
		if got := AddTax(tt.cents, tt.bp); got != tt.want {
			t.Errorf("AddTax(%d, %d) = %d, want %d", tt.cents, tt.bp, got, tt.want)
		}
	}
}
