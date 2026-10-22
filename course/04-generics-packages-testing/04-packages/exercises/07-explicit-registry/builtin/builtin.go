// Package builtin provides the standard render formats.
package builtin

import "learngo/course/04-generics-packages-testing/04-packages/exercises/07-explicit-registry/render"

// init registers the formats as a side effect of importing this package.
// TODO: delete it; callers will call RegisterAll explicitly instead.
func init() {
	render.Register("plain", Plain)
	render.Register("markdown", Markdown)
	render.Register("numbered", Numbered)
}

// Plain renders "Title: a, b", or "Title: (none)" without items.
func Plain(title string, items []string) string {
	// TODO
	return ""
}

// Markdown renders a heading and a bullet list, or "_No items._" without items.
func Markdown(title string, items []string) string {
	// TODO
	return ""
}

// Numbered renders the title followed by one "N. item" line per item.
func Numbered(title string, items []string) string {
	// TODO
	return ""
}

// RegisterAll registers "plain", "markdown" and "numbered" in r and returns
// the first error.
func RegisterAll(r *render.Registry) error {
	// TODO
	return nil
}
