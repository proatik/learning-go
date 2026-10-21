package profile

func ExamplePrintProfile() {
	PrintProfile("Gopher", 2009, true)
	// Output:
	// Name:   Gopher
	// Born:   2009
	// Active: true
}

func ExamplePrintProfile_inactive() {
	PrintProfile("Ada", 1815, false)
	// Output:
	// Name:   Ada
	// Born:   1815
	// Active: false
}
