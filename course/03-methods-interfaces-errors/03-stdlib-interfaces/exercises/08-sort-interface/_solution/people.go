package people

import (
	"cmp"
	"slices"
	"sort"
)

// Person has a name and an age.
type Person struct {
	Name string
	Age  int
}

// ByAge sorts people by age, youngest first.
type ByAge []Person

// Len returns the number of people.
func (a ByAge) Len() int { return len(a) }

// Less reports whether a[i] is younger than a[j].
func (a ByAge) Less(i, j int) bool { return a[i].Age < a[j].Age }

// Swap swaps a[i] and a[j].
func (a ByAge) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// ByName sorts people by name, A to Z.
type ByName []Person

// Len returns the number of people.
func (a ByName) Len() int { return len(a) }

// Less reports whether a[i]'s name sorts before a[j]'s.
func (a ByName) Less(i, j int) bool { return a[i].Name < a[j].Name }

// Swap swaps a[i] and a[j].
func (a ByName) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// SortOldestFirst sorts ps by age, oldest first; equal ages keep their order.
func SortOldestFirst(ps []Person) {
	sort.Stable(sort.Reverse(ByAge(ps)))
}

// SortAgeThenName sorts ps by age (youngest first), then by name.
func SortAgeThenName(ps []Person) {
	slices.SortFunc(ps, func(a, b Person) int {
		return cmp.Or(cmp.Compare(a.Age, b.Age), cmp.Compare(a.Name, b.Name))
	})
}
