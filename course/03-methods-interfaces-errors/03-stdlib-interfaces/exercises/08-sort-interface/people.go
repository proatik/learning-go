package people

// Person has a name and an age.
type Person struct {
	Name string
	Age  int
}

// ByAge sorts people by age, youngest first.
type ByAge []Person

// Len returns the number of people.
func (a ByAge) Len() int {
	// TODO
	return 0
}

// Less reports whether a[i] is younger than a[j].
func (a ByAge) Less(i, j int) bool {
	// TODO
	return false
}

// Swap swaps a[i] and a[j].
func (a ByAge) Swap(i, j int) {
	// TODO
}

// ByName sorts people by name, A to Z.
type ByName []Person

// Len returns the number of people.
func (a ByName) Len() int {
	// TODO
	return 0
}

// Less reports whether a[i]'s name sorts before a[j]'s.
func (a ByName) Less(i, j int) bool {
	// TODO
	return false
}

// Swap swaps a[i] and a[j].
func (a ByName) Swap(i, j int) {
	// TODO
}

// SortOldestFirst sorts ps by age, oldest first; equal ages keep their order.
func SortOldestFirst(ps []Person) {
	// TODO
}

// SortAgeThenName sorts ps by age (youngest first), then by name.
func SortAgeThenName(ps []Person) {
	// TODO
}
