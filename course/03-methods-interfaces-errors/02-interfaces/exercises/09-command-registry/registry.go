package commands

// Handler handles one command invocation.
type Handler interface {
	Handle(args []string) string
}

// HandlerFunc adapts an ordinary function to the Handler interface.
type HandlerFunc func(args []string) string

// Handle calls f(args).
func (f HandlerFunc) Handle(args []string) string {
	// TODO
	return ""
}

// Echo joins its arguments with a space.
type Echo struct{}

// Handle returns args joined by " ".
func (Echo) Handle(args []string) string {
	// TODO
	return ""
}

// Counter counts how many times it has been called.
type Counter struct {
	n int
}

// Handle increments the counter and returns "count=N".
func (c *Counter) Handle(args []string) string {
	// TODO
	return ""
}

// Upper returns a Handler whose output is h's output in upper case.
func Upper(h Handler) Handler {
	// TODO
	return nil
}

// Registry maps command names to handlers. The zero value is ready to use.
type Registry struct {
	handlers map[string]Handler
}

// Register adds h under name and reports whether it was added.
func (r *Registry) Register(name string, h Handler) bool {
	// TODO
	return false
}

// Names returns the registered names in ascending order.
func (r *Registry) Names() []string {
	// TODO
	return nil
}

// Dispatch runs the command described by line and returns its output.
func (r *Registry) Dispatch(line string) string {
	// TODO
	return ""
}
