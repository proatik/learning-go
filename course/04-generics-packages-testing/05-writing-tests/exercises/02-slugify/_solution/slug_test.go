package slug

import "testing"

// TestSlugify checks Slugify with a table of inputs and wanted slugs.
func TestSlugify(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"", ""},
		{"hello", "hello"},
		{"Hello World", "hello-world"},   // uppercase, one separator
		{"Go 1.27", "go-1-27"},           // digits
		{"a  --  b", "a-b"},              // a run of separators
		{"  Leading", "leading"},         // separators at the start
		{"Trailing!!", "trailing"},       // separators at the end
		{"café crème", "caf-cr-me"},      // non-ASCII letters are separators
		{"!!!", ""},                      // nothing but separators
		{"Hello, World!", "hello-world"}, // everything at once
		{"2024 Q1 report", "2024-q1-report"},
	}
	for _, tt := range tests {
		if got := Slugify(tt.in); got != tt.want {
			t.Errorf("Slugify(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}
