package census

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
	// TODO
	return nil
}

// Stats returns the statistics of each city.
func Stats(people []Person) map[string]CityStats {
	// TODO
	return nil
}

// AverageAge returns the mean age in s, or 0 if s.Count is 0.
func AverageAge(s CityStats) float64 {
	// TODO
	return 0
}

// Cities returns the city names in stats, sorted by Count descending,
// then by name ascending.
func Cities(stats map[string]CityStats) []string {
	// TODO
	return nil
}
