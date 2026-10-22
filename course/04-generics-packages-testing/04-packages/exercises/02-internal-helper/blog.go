// Package blog creates blog posts with unique URL slugs.
package blog

import (
	"errors"

	// Allowed: blog.go sits in the parent directory of internal/.
	"learngo/course/04-generics-packages-testing/04-packages/exercises/02-internal-helper/internal/slug"
)

// ErrEmptyTitle is returned when a title has no letters or digits.
var ErrEmptyTitle = errors.New("blog: title has no letters or digits")

// Post is a blog post.
type Post struct {
	Title string
	Slug  string
}

// NewPost returns a post whose slug is made from title, or ErrEmptyTitle
// when the slug would be empty.
func NewPost(title string) (Post, error) {
	_ = slug.Make // TODO
	return Post{}, nil
}

// URL returns the post's path, "/posts/<slug>".
func (p Post) URL() string {
	// TODO
	return ""
}

// NewPosts returns one post per title. Duplicate slugs get the suffix "-N"
// with the smallest free N >= 2. An empty slug stops with an error that
// wraps ErrEmptyTitle and names the title's index.
func NewPosts(titles []string) ([]Post, error) {
	// TODO: track taken slugs in a map[string]bool
	return nil, nil
}
