package tracer

import "io"

// Tracer writes indented enter/exit lines for nested function calls.
type Tracer struct {
	W     io.Writer // where trace lines go; nil disables tracing
	depth int       // current nesting level
}

// Enter writes "<indent>> name", increases the nesting depth and returns a
// function that restores the depth and writes "<indent>< name". Use it as
// defer t.Enter("name")().
func (t *Tracer) Enter(name string) func() {
	// TODO
	return func() {}
}
