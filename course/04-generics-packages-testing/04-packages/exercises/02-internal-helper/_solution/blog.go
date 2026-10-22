// Package blog creates blog posts with unique URL slugs.
package blog

import (
	"errors"
	"fmt"

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
	s := slug.Make(title)
	if s == "" {
		return Post{}, ErrEmptyTitle
	}
	return Post{Title: title, Slug: s}, nil
}

// URL returns the post's path, "/posts/<slug>".
func (p Post) URL() string { return "/posts/" + p.Slug }

// NewPosts returns one post per title. Duplicate slugs get the suffix "-N"
// with the smallest free N >= 2. An empty slug stops with an error that
// wraps ErrEmptyTitle and names the title's index.
func NewPosts(titles []string) ([]Post, error) {
	posts := make([]Post, 0, len(titles))
	taken := make(map[string]bool, len(titles))
	for i, title := range titles {
		p, err := NewPost(title)
		if err != nil {
			return nil, fmt.Errorf("title %d: %w", i, err)
		}
		if taken[p.Slug] {
			base := p.Slug
			for n := 2; taken[p.Slug]; n++ {
				p.Slug = fmt.Sprintf("%s-%d", base, n)
			}
		}
		taken[p.Slug] = true
		posts = append(posts, p)
	}
	return posts, nil
}
