package zoo

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

// Name returns the dog's nickname.
func (d Dog) Name() string {
	// TODO
	return ""
}

// Sound returns "Woof".
func (d Dog) Sound() string {
	// TODO
	return ""
}

// Name returns the cat's nickname.
func (c Cat) Name() string {
	// TODO
	return ""
}

// Sound returns "Hiss" for a grumpy cat, "Meow" otherwise.
func (c Cat) Sound() string {
	// TODO
	return ""
}

// Name returns "<Model>-v<Version>".
func (r Robot) Name() string {
	// TODO
	return ""
}

// Sound returns "Beep boop".
func (r Robot) Sound() string {
	// TODO
	return ""
}

// Greet returns "<Name> says <Sound>!".
func Greet(s Speaker) string {
	// TODO
	return ""
}

// Chorus returns the greetings of all speakers joined by a space.
func Chorus(speakers []Speaker) string {
	// TODO
	return ""
}

// CountSound reports how many speakers make exactly the given sound.
func CountSound(speakers []Speaker, sound string) int {
	// TODO
	return 0
}
