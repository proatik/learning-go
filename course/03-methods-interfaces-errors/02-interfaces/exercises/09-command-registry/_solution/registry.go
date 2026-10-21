package commands

import (
	"fmt"
	"maps"
	"slices"
	"strings"
	"unicode"
)

// Handler handles one command invocation.
type Handler interface {
	Handle(args []string) string
}

// HandlerFunc adapts an ordinary function to the Handler interface.
type HandlerFunc func(args []string) string

// Handle calls f(args).
func (f HandlerFunc) Handle(args []string) string { return f(args) }

// Echo joins its arguments with a space.
type Echo struct{}

// Handle returns args joined by " ".
func (Echo) Handle(args []string) string { return strings.Join(args, " ") }

// Counter counts how many times it has been called.
type Counter struct {
	n int
}

// Handle increments the counter and returns "count=N".
func (c *Counter) Handle(args []string) string {
	c.n++
	return fmt.Sprintf("count=%d", c.n)
}

var (
	_ Handler = HandlerFunc(nil)
	_ Handler = Echo{}
	_ Handler = (*Counter)(nil)
)

// Upper returns a Handler whose output is h's output in upper case.
func Upper(h Handler) Handler {
	return HandlerFunc(func(args []string) string {
		return strings.ToUpper(h.Handle(args))
	})
}

// Registry maps command names to handlers. The zero value is ready to use.
type Registry struct {
	handlers map[string]Handler
}

// isNil reports whether h is nil or holds a nil value of a known type.
func isNil(h Handler) bool {
	switch x := h.(type) {
	case nil:
		return true
	case HandlerFunc:
		return x == nil
	case *Counter:
		return x == nil
	}
	return false
}

// Register adds h under name and reports whether it was added.
func (r *Registry) Register(name string, h Handler) bool {
	if name == "" || strings.ContainsFunc(name, unicode.IsSpace) || isNil(h) {
		return false
	}
	if _, dup := r.handlers[name]; dup {
		return false
	}
	if r.handlers == nil {
		r.handlers = map[string]Handler{}
	}
	r.handlers[name] = h
	return true
}

// Names returns the registered names in ascending order.
func (r *Registry) Names() []string {
	return slices.Sorted(maps.Keys(r.handlers))
}

// Dispatch runs the command described by line and returns its output.
func (r *Registry) Dispatch(line string) string {
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return ""
	}
	name, args := fields[0], fields[1:]
	if h, ok := r.handlers[name]; ok {
		return h.Handle(args)
	}
	if name == "help" {
		if len(r.handlers) == 0 {
			return "no commands"
		}
		return "commands: " + strings.Join(r.Names(), " ")
	}
	return "unknown command: " + name
}
