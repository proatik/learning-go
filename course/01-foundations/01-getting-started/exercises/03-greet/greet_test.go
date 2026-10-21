package greet

import "testing"

func TestGreet(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"Gopher", "Hello, Gopher!"},
		{"Rob", "Hello, Rob!"},
		{"", "Hello, World!"},
		{"Ken Thompson", "Hello, Ken Thompson!"},
	}
	for _, tt := range tests {
		if got := Greet(tt.name); got != tt.want {
			t.Errorf("Greet(%q) = %q, want %q", tt.name, got, tt.want)
		}
	}
}
