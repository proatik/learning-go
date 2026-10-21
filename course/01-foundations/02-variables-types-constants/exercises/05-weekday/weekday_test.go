package weekday

import "testing"

func TestConstants(t *testing.T) {
	days := []struct {
		name string
		d    Weekday
		want int
	}{
		{"Sunday", Sunday, 0}, {"Monday", Monday, 1}, {"Tuesday", Tuesday, 2},
		{"Wednesday", Wednesday, 3}, {"Thursday", Thursday, 4}, {"Friday", Friday, 5},
		{"Saturday", Saturday, 6},
	}
	for _, tt := range days {
		if int(tt.d) != tt.want {
			t.Errorf("%s = %d, want %d", tt.name, tt.d, tt.want)
		}
	}
}

func TestIsWeekend(t *testing.T) {
	want := []bool{true, false, false, false, false, false, true}
	for d := range Weekday(7) {
		if got := IsWeekend(d); got != want[d] {
			t.Errorf("IsWeekend(%d) = %v, want %v", d, got, want[d])
		}
	}
}

func TestNext(t *testing.T) {
	tests := []struct{ d, want Weekday }{
		{Sunday, Monday}, {Monday, Tuesday}, {Wednesday, Thursday},
		{Friday, Saturday}, {Saturday, Sunday},
	}
	for _, tt := range tests {
		if got := Next(tt.d); got != tt.want {
			t.Errorf("Next(%d) = %d, want %d", tt.d, got, tt.want)
		}
	}
}

func TestPrev(t *testing.T) {
	tests := []struct{ d, want Weekday }{
		{Sunday, Saturday}, {Monday, Sunday}, {Thursday, Wednesday},
		{Saturday, Friday},
	}
	for _, tt := range tests {
		if got := Prev(tt.d); got != tt.want {
			t.Errorf("Prev(%d) = %d, want %d", tt.d, got, tt.want)
		}
	}
}

func TestNextPrevRoundTrip(t *testing.T) {
	for d := range Weekday(7) {
		if got := Prev(Next(d)); got != d {
			t.Errorf("Prev(Next(%d)) = %d, want %d", d, got, d)
		}
	}
}
