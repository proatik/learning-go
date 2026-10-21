package census

import (
	"maps"
	"math"
	"slices"
	"testing"
)

func people() []Person {
	return []Person{
		{Name: "Ann", Age: 31, City: "Oslo"},
		{Name: "Bob", Age: 45, City: "Rome"},
		{Name: "Cy", Age: 45, City: "Oslo"},
		{Name: "Di", Age: 45, City: "Oslo"},
		{Name: "Ed", Age: 20, City: "Lima"},
		{Name: "Flo", Age: 62, City: "Lima"},
		{Name: "Gus", Age: 9, City: "Rome"},
		{Name: "Hal", Age: 0, City: "Kiev"},
	}
}

func TestGroupByCity(t *testing.T) {
	got := GroupByCity(people())
	want := map[string][]string{
		"Oslo": {"Ann", "Cy", "Di"},
		"Rome": {"Bob", "Gus"},
		"Lima": {"Ed", "Flo"},
		"Kiev": {"Hal"},
	}
	if !maps.EqualFunc(got, want, slices.Equal[[]string]) {
		t.Errorf("GroupByCity(people) = %v, want %v", got, want)
	}
	for _, in := range [][]Person{nil, {}} {
		if got := GroupByCity(in); got == nil || len(got) != 0 {
			t.Errorf("GroupByCity(%#v) = %#v, want a non-nil empty map", in, got)
		}
	}
}

func TestStats(t *testing.T) {
	got := Stats(people())
	want := map[string]CityStats{
		"Oslo": {Count: 3, TotalAge: 121, Oldest: Person{"Cy", 45, "Oslo"}},
		"Rome": {Count: 2, TotalAge: 54, Oldest: Person{"Bob", 45, "Rome"}},
		"Lima": {Count: 2, TotalAge: 82, Oldest: Person{"Flo", 62, "Lima"}},
		"Kiev": {Count: 1, TotalAge: 0, Oldest: Person{"Hal", 0, "Kiev"}},
	}
	if len(got) != len(want) {
		t.Errorf("Stats(people) has %d cities, want %d: %+v", len(got), len(want), got)
	}
	for city, w := range want {
		if g, ok := got[city]; !ok {
			t.Errorf("Stats(people) is missing city %q", city)
		} else if g != w {
			t.Errorf("Stats(people)[%q] = %+v, want %+v", city, g, w)
		}
	}
	for _, in := range [][]Person{nil, {}} {
		if got := Stats(in); got == nil || len(got) != 0 {
			t.Errorf("Stats(%#v) = %#v, want a non-nil empty map", in, got)
		}
	}
}

func TestAverageAge(t *testing.T) {
	tests := []struct {
		s    CityStats
		want float64
	}{
		{CityStats{Count: 2, TotalAge: 54}, 27},
		{CityStats{Count: 3, TotalAge: 121}, 121.0 / 3},
		{CityStats{Count: 1, TotalAge: 0}, 0},
		{CityStats{Count: 4, TotalAge: 10}, 2.5},
		{CityStats{}, 0},
	}
	for _, tt := range tests {
		if got := AverageAge(tt.s); math.Abs(got-tt.want) > 1e-9 {
			t.Errorf("AverageAge(%+v) = %v, want %v", tt.s, got, tt.want)
		}
	}
}

func TestCities(t *testing.T) {
	tests := []struct {
		stats map[string]CityStats
		want  []string
	}{
		{Stats(people()), []string{"Oslo", "Lima", "Rome", "Kiev"}},
		{
			map[string]CityStats{"b": {Count: 1}, "a": {Count: 1}, "c": {Count: 5}},
			[]string{"c", "a", "b"},
		},
		{map[string]CityStats{"only": {Count: 2}}, []string{"only"}},
		{map[string]CityStats{}, []string{}},
	}
	for _, tt := range tests {
		got := Cities(tt.stats)
		if !slices.Equal(got, tt.want) {
			t.Errorf("Cities(%+v) = %q, want %q", tt.stats, got, tt.want)
		}
	}
}
