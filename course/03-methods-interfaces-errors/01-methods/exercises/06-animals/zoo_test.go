package zoo

import "testing"

func TestAnimal(t *testing.T) {
	tests := []struct {
		a            Animal
		speak, descr string
	}{
		{Animal{"Tom", "meow"}, "Tom says meow", "Tom is an animal"},
		{Animal{"Daisy", "moo"}, "Daisy says moo", "Daisy is an animal"},
		{Animal{"Ëlo", "…"}, "Ëlo says …", "Ëlo is an animal"},
	}
	for _, tt := range tests {
		if got := tt.a.Speak(); got != tt.speak {
			t.Errorf("%+v.Speak() = %q, want %q", tt.a, got, tt.speak)
		}
		if got := tt.a.Describe(); got != tt.descr {
			t.Errorf("%+v.Describe() = %q, want %q", tt.a, got, tt.descr)
		}
	}
}

func TestNewDog(t *testing.T) {
	d := NewDog("Rex", "beagle")
	want := Dog{Animal: Animal{Name: "Rex", Sound: "woof"}, Breed: "beagle"}
	if d != want {
		t.Errorf("NewDog(%q, %q) = %+v, want %+v", "Rex", "beagle", d, want)
	}
}

func TestDog(t *testing.T) {
	tests := []struct {
		name, breed               string
		descr, speak, fetch, base string
	}{
		{"Rex", "beagle", "Rex is a beagle", "Rex says woof", "Rex fetches the ball", "Rex is an animal"},
		{"Rex", "", "Rex is a dog", "Rex says woof", "Rex fetches the ball", "Rex is an animal"},
		{"Lola", "poodle", "Lola is a poodle", "Lola says woof", "Lola fetches the ball", "Lola is an animal"},
	}
	for _, tt := range tests {
		d := NewDog(tt.name, tt.breed)
		if got := d.Describe(); got != tt.descr {
			t.Errorf("NewDog(%q, %q).Describe() = %q, want %q", tt.name, tt.breed, got, tt.descr)
		}
		if got := d.Speak(); got != tt.speak {
			t.Errorf("NewDog(%q, %q).Speak() = %q, want %q", tt.name, tt.breed, got, tt.speak)
		}
		if got := d.Fetch(); got != tt.fetch {
			t.Errorf("NewDog(%q, %q).Fetch() = %q, want %q", tt.name, tt.breed, got, tt.fetch)
		}
		if got := d.Animal.Describe(); got != tt.base {
			t.Errorf("NewDog(%q, %q).Animal.Describe() = %q, want %q", tt.name, tt.breed, got, tt.base)
		}
	}
}

func TestGuideDog(t *testing.T) {
	tests := []struct {
		g     GuideDog
		descr string
	}{
		{GuideDog{NewDog("Max", "labrador"), "Sam"}, "Max is a labrador guiding Sam"},
		{GuideDog{NewDog("Max", ""), "Sam"}, "Max is a dog guiding Sam"},
		{GuideDog{NewDog("Max", "labrador"), ""}, "Max is a labrador"},
	}
	for _, tt := range tests {
		if got := tt.g.Describe(); got != tt.descr {
			t.Errorf("%+v.Describe() = %q, want %q", tt.g, got, tt.descr)
		}
	}
	g := GuideDog{NewDog("Max", "labrador"), "Sam"}
	if got, want := g.Speak(), "Max says woof"; got != want {
		t.Errorf("GuideDog.Speak() = %q, want %q", got, want)
	}
	if got, want := g.Fetch(), "Max fetches the ball"; got != want {
		t.Errorf("GuideDog.Fetch() = %q, want %q", got, want)
	}
}

func TestRename(t *testing.T) {
	a := Animal{"Tom", "meow"}
	a.Rename("Tim")
	if a.Name != "Tim" {
		t.Errorf("Animal.Rename(%q): Name = %q, want %q", "Tim", a.Name, "Tim")
	}

	d := NewDog("Rex", "beagle")
	d.Rename("Rexy")
	if got, want := d.Describe(), "Rexy is a beagle"; got != want {
		t.Errorf("Dog.Rename(%q) then Describe() = %q, want %q", "Rexy", got, want)
	}

	g := GuideDog{NewDog("Max", "labrador"), "Sam"}
	g.Rename("Buddy")
	if got, want := g.Describe(), "Buddy is a labrador guiding Sam"; got != want {
		t.Errorf("GuideDog.Rename(%q) then Describe() = %q, want %q", "Buddy", got, want)
	}
}
