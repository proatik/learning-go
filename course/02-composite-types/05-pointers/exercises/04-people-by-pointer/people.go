package people

// Person is a named person with an age in years.
type Person struct {
	Name string
	Age  int
}

// Birthday makes the person p points to one year older.
// It does nothing if p is nil.
func Birthday(p *Person) {
	// TODO
}

// BirthdayAll makes every person in people one year older.
func BirthdayAll(people []Person) {
	// TODO
}

// Find returns a pointer to the first element of people whose Name is name,
// or nil if there is none.
func Find(people []Person, name string) *Person {
	// TODO
	return nil
}

// Oldest returns a pointer to the element of people with the highest Age
// (the first one on ties), or nil if people is empty.
func Oldest(people []Person) *Person {
	// TODO
	return nil
}
