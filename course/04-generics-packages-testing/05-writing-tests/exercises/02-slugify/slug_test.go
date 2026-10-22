package slug

import "testing"

// TestSlugify checks Slugify with a table of inputs and wanted slugs.
func TestSlugify(t *testing.T) {
	// TODO: a table of {in, want string} cases.
	// Think about each rule separately: letters, digits, one separator,
	// a run of separators, and separators at the start and at the end.
	// Report mismatches as: "Slugify(%q) = %q, want %q"
	_ = Slugify
}
