package census

import (
	"cmp"
	"slices"
)

// Person is someone living in a city.
type Person struct {
	Name string
	Age  int
	City string
}

// CityStats summarizes the people of one city.
type CityStats struct {
	Count    int
	TotalAge int
	Oldest   Person // first person with the highest age
}

// GroupByCity returns the names of the people in each city, in input order.
func GroupByCity(people []Person) map[string][]string {
	groups := make(map[string][]string)
	for _, p := range people {
		groups[p.City] = append(groups[p.City], p.Name) // nil slice for a new city is fine
	}
	return groups
}

// Stats returns the statistics of each city.
func Stats(people []Person) map[string]CityStats {
	stats := make(map[string]CityStats)
	for _, p := range people {
		s := stats[p.City] // copy (zero CityStats for a new city)
		if s.Count == 0 || p.Age > s.Oldest.Age {
			s.Oldest = p
		}
		s.Count++
		s.TotalAge += p.Age
		stats[p.City] = s // store the modified copy back
	}
	return stats
}

// AverageAge returns the mean age in s, or 0 if s.Count is 0.
func AverageAge(s CityStats) float64 {
	if s.Count == 0 {
		return 0
	}
	return float64(s.TotalAge) / float64(s.Count)
}

// Cities returns the city names in stats, sorted by Count descending,
// then by name ascending.
func Cities(stats map[string]CityStats) []string {
	cities := make([]string, 0, len(stats))
	for city := range stats {
		cities = append(cities, city)
	}
	slices.SortFunc(cities, func(a, b string) int {
		return cmp.Or(
			cmp.Compare(stats[b].Count, stats[a].Count),
			cmp.Compare(a, b),
		)
	})
	return cities
}
