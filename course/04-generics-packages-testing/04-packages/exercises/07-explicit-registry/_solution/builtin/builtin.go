// Package builtin provides the standard render formats.
package builtin

import (
	"fmt"
	"strings"

	"learngo/course/04-generics-packages-testing/04-packages/exercises/07-explicit-registry/render"
)

// Plain renders "Title: a, b", or "Title: (none)" without items.
func Plain(title string, items []string) string {
	if len(items) == 0 {
		return title + ": (none)"
	}
	return title + ": " + strings.Join(items, ", ")
}

// Markdown renders a heading and a bullet list, or "_No items._" without items.
func Markdown(title string, items []string) string {
	if len(items) == 0 {
		return "# " + title + "\n\n_No items._"
	}
	return "# " + title + "\n\n- " + strings.Join(items, "\n- ")
}

// Numbered renders the title followed by one "N. item" line per item.
func Numbered(title string, items []string) string {
	var b strings.Builder
	b.WriteString(title)
	for i, item := range items {
		fmt.Fprintf(&b, "\n%d. %s", i+1, item)
	}
	return b.String()
}

// RegisterAll registers "plain", "markdown" and "numbered" in r and returns
// the first error.
func RegisterAll(r *render.Registry) error {
	if err := r.Register("plain", Plain); err != nil {
		return err
	}
	if err := r.Register("markdown", Markdown); err != nil {
		return err
	}
	return r.Register("numbered", Numbered)
}
