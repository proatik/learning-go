package profile

// User is a sign-up request.
type User struct {
	Name  string
	Email string
	Age   int
}

// ValidationError reports which field of a User is invalid and why.
type ValidationError struct {
	Field  string
	Reason string
}

// Error returns "<field>: <reason>".
func (e *ValidationError) Error() string {
	// TODO
	return ""
}

// Validate returns the first problem with u as a *ValidationError, or nil.
func Validate(u User) error {
	// TODO
	return nil
}
