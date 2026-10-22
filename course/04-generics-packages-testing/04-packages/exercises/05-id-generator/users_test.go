package users_test

import (
	"errors"
	"testing"

	users "learngo/course/04-generics-packages-testing/04-packages/exercises/05-id-generator"
)

func TestDirectory(t *testing.T) {
	d := users.NewDirectory("u-")
	for i, name := range []string{"Ann", "Bob", "Ünal"} {
		want := []string{"u-1", "u-2", "u-3"}[i]
		if got, err := d.Add(name); got != want || err != nil {
			t.Errorf("Add(%q) = %q, %v, want %q, <nil>", name, got, err, want)
		}
	}
	tests := []struct {
		id     string
		want   string
		wantOK bool
	}{
		{"u-1", "Ann", true},
		{"u-3", "Ünal", true},
		{"u-4", "", false},
		{"1", "", false},
	}
	for _, tt := range tests {
		if got, ok := d.Name(tt.id); got != tt.want || ok != tt.wantOK {
			t.Errorf("Name(%q) = %q, %v, want %q, %v", tt.id, got, ok, tt.want, tt.wantOK)
		}
	}
	if got := d.Len(); got != 3 {
		t.Errorf("Len() = %d, want 3", got)
	}
}

func TestEmptyNameKeepsNumbering(t *testing.T) {
	d := users.NewDirectory("e-")
	if got, err := d.Add(""); got != "" || !errors.Is(err, users.ErrEmptyName) {
		t.Errorf("Add(%q) = %q, %v, want %q, %v", "", got, err, "", users.ErrEmptyName)
	}
	if got, err := d.Add("Cy"); got != "e-1" || err != nil {
		t.Errorf("Add(%q) after a rejected name = %q, %v, want %q, <nil>", "Cy", got, err, "e-1")
	}
}

func TestDirectoriesAreIndependent(t *testing.T) {
	a := users.NewDirectory("u-")
	b := users.NewDirectory("u-")
	a.Add("Ann")
	a.Add("Bob")
	if got, _ := b.Add("Cy"); got != "u-1" {
		t.Errorf("second directory first Add(%q) = %q, want %q", "Cy", got, "u-1")
	}
	if got, _ := a.Add("Dee"); got != "u-3" {
		t.Errorf("first directory third Add(%q) = %q, want %q", "Dee", got, "u-3")
	}
}
