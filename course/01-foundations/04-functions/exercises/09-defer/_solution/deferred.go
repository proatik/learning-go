package deferred

import (
	"errors"
	"fmt"
	"strconv"
)

// Stack prints "push 1" .. "push n", then "pop n" .. "pop 1",
// using defer to schedule the pops.
func Stack(n int) {
	for i := range n {
		fmt.Println("push", i+1)
		// The argument i+1 is evaluated now; the call runs when Stack returns.
		defer fmt.Println("pop", i+1)
	}
}

// ParseAge parses s as an age between 0 and 150 inclusive.
// Every error it returns is prefixed with `ParseAge("<s>"): `.
func ParseAge(s string) (age int, err error) {
	// Runs after every return statement below has set age and err.
	defer func() {
		if err != nil {
			age = 0
			err = fmt.Errorf("ParseAge(%q): %v", s, err)
		}
	}()

	age, err = strconv.Atoi(s)
	if err != nil {
		return age, err
	}
	if age < 0 || age > 150 {
		return age, errors.New("age out of range")
	}
	return age, nil
}
