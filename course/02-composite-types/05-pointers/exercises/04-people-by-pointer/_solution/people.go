package people

// Person is a named person with an age in years.
type Person struct {
	Name string
	Age  int
}

// Birthday makes the person p points to one year older.
// It does nothing if p is nil.
func Birthday(p *Person) {
	if p == nil {
		return
	}
	p.Age++ // p.Age is shorthand for (*p).Age
}

// BirthdayAll makes every person in people one year older.
func BirthdayAll(people []Person) {
	for i := range people {
		Birthday(&people[i]) // address of the element itself, not of a copy
	}
}

// Find returns a pointer to the first element of people whose Name is name,
// or nil if there is none.
func Find(people []Person, name string) *Person {
	for i := range people {
		if people[i].Name == name {
			return &people[i]
		}
	}
	return nil
}

// Oldest returns a pointer to the element of people with the highest Age
// (the first one on ties), or nil if people is empty.
func Oldest(people []Person) *Person {
	var best *Person // nil until we have seen someone
	for i := range people {
		if best == nil || people[i].Age > best.Age {
			best = &people[i]
		}
	}
	return best
}
