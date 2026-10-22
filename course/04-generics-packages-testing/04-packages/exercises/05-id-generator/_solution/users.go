// Package users keeps a directory of users with generated IDs.
package users

import (
	"errors"

	"learngo/course/04-generics-packages-testing/04-packages/exercises/05-id-generator/ids"
)

// ErrEmptyName is returned by Add for an empty name.
var ErrEmptyName = errors.New("users: empty name")

// Directory maps generated IDs to user names.
type Directory struct {
	ids   *ids.Generator
	names map[string]string
}

// NewDirectory returns an empty directory whose IDs start with prefix.
func NewDirectory(prefix string) *Directory {
	return &Directory{ids: ids.New(prefix), names: map[string]string{}}
}

// Add stores name under a new ID and returns the ID, or ErrEmptyName.
func (d *Directory) Add(name string) (string, error) {
	if name == "" {
		return "", ErrEmptyName
	}
	id := d.ids.Next()
	d.names[id] = name
	return id, nil
}

// Name returns the name stored under id and whether it exists.
func (d *Directory) Name(id string) (string, bool) {
	name, ok := d.names[id]
	return name, ok
}

// Len returns the number of users.
func (d *Directory) Len() int {
	return len(d.names)
}
