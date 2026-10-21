package zoo

// Animal is the base type embedded by more specific animals.
type Animal struct {
	Name  string
	Sound string
}

// Speak returns "<Name> says <Sound>".
func (a Animal) Speak() string {
	return a.Name + " says " + a.Sound
}

// Describe returns "<Name> is an animal".
func (a Animal) Describe() string {
	return a.Name + " is an animal"
}

// Rename changes the animal's name in place.
func (a *Animal) Rename(name string) {
	a.Name = name
}

// Dog is an Animal with a breed.
type Dog struct {
	Animal
	Breed string
}

// NewDog returns a Dog that says "woof".
func NewDog(name, breed string) Dog {
	return Dog{Animal: Animal{Name: name, Sound: "woof"}, Breed: breed}
}

// Describe returns "<Name> is a <Breed>", or "<Name> is a dog" without a breed.
func (d Dog) Describe() string {
	breed := d.Breed
	if breed == "" {
		breed = "dog"
	}
	return d.Name + " is a " + breed
}

// Fetch returns "<Name> fetches the ball".
func (d Dog) Fetch() string {
	return d.Name + " fetches the ball"
}

// GuideDog is a Dog working with a handler.
type GuideDog struct {
	Dog
	Handler string
}

// Describe extends Dog.Describe with " guiding <Handler>".
func (g GuideDog) Describe() string {
	s := g.Dog.Describe()
	if g.Handler != "" {
		s += " guiding " + g.Handler
	}
	return s
}
