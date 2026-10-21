# City Statistics

Group people by city and keep per-city statistics in a **map of structs**.

```go
type Person struct {
	Name string
	Age  int
	City string
}

type CityStats struct {
	Count    int
	TotalAge int
	Oldest   Person
}

func GroupByCity(people []Person) map[string][]string   // city → names
func Stats(people []Person) map[string]CityStats        // city → stats
func AverageAge(s CityStats) float64
func Cities(stats map[string]CityStats) []string        // ordered city names
```

```go
people := []Person{
	{Name: "Ann", Age: 31, City: "Oslo"},
	{Name: "Bob", Age: 45, City: "Rome"},
	{Name: "Cy", Age: 45, City: "Oslo"},
	{Name: "Di", Age: 45, City: "Oslo"},
}
```

| Call | Result |
|---|---|
| `GroupByCity(people)` | `map[Oslo:[Ann Cy Di] Rome:[Bob]]` |
| `Stats(people)["Oslo"]` | `CityStats{Count: 3, TotalAge: 121, Oldest: {Cy 45 Oslo}}` |
| `AverageAge(Stats(people)["Rome"])` | `45` |
| `AverageAge(CityStats{})` | `0` |
| `Cities(Stats(people))` | `[Oslo Rome]` |

## Rules
- Names in `GroupByCity` keep the input order.
- `Oldest` is the **first** person (in input order) with the highest age in that city.
- You can't write `stats[city].Count++`: read the struct, modify the copy, store it back.
- `Cities` sorts by `Count` **descending**, then city name ascending.
- Both maps are non-nil (empty for no people). `AverageAge` returns `0` when `Count` is 0.
