package zoo

import (
	"fmt"
	"strings"
)

// Speaker is anything with a name that makes a sound.
type Speaker interface {
	Name() string
	Sound() string
}

// Dog is a dog called Nick.
type Dog struct{ Nick string }

// Cat is a cat called Nick; grumpy cats hiss.
type Cat struct {
	Nick   string
	Grumpy bool
}

// Robot is a robot identified by model and version.
type Robot struct {
	Model   string
	Version int
}

var (
	_ Speaker = Dog{}
	_ Speaker = Cat{}
	_ Speaker = Robot{}
)

// Name returns the dog's nickname.
func (d Dog) Name() string { return d.Nick }

// Sound returns "Woof".
func (d Dog) Sound() string { return "Woof" }

// Name returns the cat's nickname.
func (c Cat) Name() string { return c.Nick }

// Sound returns "Hiss" for a grumpy cat, "Meow" otherwise.
func (c Cat) Sound() string {
	if c.Grumpy {
		return "Hiss"
	}
	return "Meow"
}

// Name returns "<Model>-v<Version>".
func (r Robot) Name() string { return fmt.Sprintf("%s-v%d", r.Model, r.Version) }

// Sound returns "Beep boop".
func (r Robot) Sound() string { return "Beep boop" }

// Greet returns "<Name> says <Sound>!".
func Greet(s Speaker) string {
	return s.Name() + " says " + s.Sound() + "!"
}

// Chorus returns the greetings of all speakers joined by a space.
func Chorus(speakers []Speaker) string {
	parts := make([]string, 0, len(speakers))
	for _, s := range speakers {
		parts = append(parts, Greet(s))
	}
	return strings.Join(parts, " ")
}

// CountSound reports how many speakers make exactly the given sound.
func CountSound(speakers []Speaker, sound string) int {
	n := 0
	for _, s := range speakers {
		if s.Sound() == sound {
			n++
		}
	}
	return n
}
