package zoo

// Animal is the base type embedded by more specific animals.
type Animal struct {
	Name  string
	Sound string
}

// Speak returns "<Name> says <Sound>".
func (a Animal) Speak() string {
	// TODO
	return ""
}

// Describe returns "<Name> is an animal".
func (a Animal) Describe() string {
	// TODO
	return ""
}

// Rename changes the animal's name in place.
func (a *Animal) Rename(name string) {
	// TODO
}

// Dog is an Animal with a breed.
type Dog struct {
	Animal
	Breed string
}

// NewDog returns a Dog that says "woof".
func NewDog(name, breed string) Dog {
	// TODO
	return Dog{}
}

// Describe returns "<Name> is a <Breed>", or "<Name> is a dog" without a breed.
func (d Dog) Describe() string {
	// TODO
	return ""
}

// Fetch returns "<Name> fetches the ball".
func (d Dog) Fetch() string {
	// TODO
	return ""
}

// GuideDog is a Dog working with a handler.
type GuideDog struct {
	Dog
	Handler string
}

// Describe extends Dog.Describe with " guiding <Handler>".
func (g GuideDog) Describe() string {
	// TODO
	return ""
}
