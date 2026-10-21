package zoo_test

import (
	"testing"

	zoo "learngo/course/03-methods-interfaces-errors/02-interfaces/exercises/02-speakers"
)

func TestMethods(t *testing.T) {
	tests := []struct {
		s           zoo.Speaker
		name, sound string
	}{
		{zoo.Dog{Nick: "Rex"}, "Rex", "Woof"},
		{zoo.Dog{}, "", "Woof"},
		{zoo.Cat{Nick: "Tom"}, "Tom", "Meow"},
		{zoo.Cat{Nick: "Tom", Grumpy: true}, "Tom", "Hiss"},
		{zoo.Cat{Nick: "Ünal"}, "Ünal", "Meow"},
		{zoo.Robot{Model: "RX", Version: 2}, "RX-v2", "Beep boop"},
		{zoo.Robot{Model: "T", Version: 1000}, "T-v1000", "Beep boop"},
		{zoo.Robot{}, "-v0", "Beep boop"},
	}
	for _, tt := range tests {
		if got := tt.s.Name(); got != tt.name {
			t.Errorf("%+v.Name() = %q, want %q", tt.s, got, tt.name)
		}
		if got := tt.s.Sound(); got != tt.sound {
			t.Errorf("%+v.Sound() = %q, want %q", tt.s, got, tt.sound)
		}
	}
}

func TestGreet(t *testing.T) {
	tests := []struct {
		s    zoo.Speaker
		want string
	}{
		{zoo.Dog{Nick: "Rex"}, "Rex says Woof!"},
		{zoo.Cat{Nick: "Tom", Grumpy: true}, "Tom says Hiss!"},
		{zoo.Robot{Model: "RX", Version: 2}, "RX-v2 says Beep boop!"},
	}
	for _, tt := range tests {
		if got := zoo.Greet(tt.s); got != tt.want {
			t.Errorf("Greet(%+v) = %q, want %q", tt.s, got, tt.want)
		}
	}
}

func TestChorus(t *testing.T) {
	tests := []struct {
		in   []zoo.Speaker
		want string
	}{
		{nil, ""},
		{[]zoo.Speaker{}, ""},
		{[]zoo.Speaker{zoo.Dog{Nick: "Rex"}}, "Rex says Woof!"},
		{[]zoo.Speaker{zoo.Dog{Nick: "Rex"}, zoo.Robot{Model: "RX", Version: 2}}, "Rex says Woof! RX-v2 says Beep boop!"},
		{[]zoo.Speaker{zoo.Cat{Nick: "A"}, zoo.Cat{Nick: "B", Grumpy: true}, zoo.Dog{Nick: "C"}}, "A says Meow! B says Hiss! C says Woof!"},
	}
	for _, tt := range tests {
		if got := zoo.Chorus(tt.in); got != tt.want {
			t.Errorf("Chorus(%+v) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestCountSound(t *testing.T) {
	mixed := []zoo.Speaker{zoo.Dog{Nick: "A"}, zoo.Cat{Nick: "B"}, zoo.Dog{Nick: "C"}, zoo.Cat{Nick: "D", Grumpy: true}, zoo.Robot{Model: "R"}}
	tests := []struct {
		in    []zoo.Speaker
		sound string
		want  int
	}{
		{nil, "Woof", 0},
		{mixed, "Woof", 2},
		{mixed, "Meow", 1},
		{mixed, "Hiss", 1},
		{mixed, "Beep boop", 1},
		{mixed, "woof", 0},
		{mixed, "", 0},
	}
	for _, tt := range tests {
		if got := zoo.CountSound(tt.in, tt.sound); got != tt.want {
			t.Errorf("CountSound(%+v, %q) = %d, want %d", tt.in, tt.sound, got, tt.want)
		}
	}
}
