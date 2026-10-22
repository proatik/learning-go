package blog_test

import (
	"errors"
	"testing"

	blog "learngo/course/04-generics-packages-testing/04-packages/exercises/02-internal-helper"
)

func TestNewPost(t *testing.T) {
	tests := []struct {
		title, slug string
	}{
		{"Hello, World!", "hello-world"},
		{"hello", "hello"},
		{"  Go 1.27 -- Released  ", "go-1-27-released"},
		{"Ünïcode Ok", "ünïcode-ok"},
		{"ÉCOLE", "école"},
		{"a_b.c/d", "a-b-c-d"},
		{"--edge--", "edge"},
		{"MixedCase42", "mixedcase42"},
		{"x", "x"},
	}
	for _, tt := range tests {
		got, err := blog.NewPost(tt.title)
		want := blog.Post{Title: tt.title, Slug: tt.slug}
		if err != nil || got != want {
			t.Errorf("NewPost(%q) = %+v, %v, want %+v, <nil>", tt.title, got, err, want)
		}
	}
}

func TestNewPostEmpty(t *testing.T) {
	for _, title := range []string{"", "   ", "!!!", "?-?", "🙂"} {
		got, err := blog.NewPost(title)
		if !errors.Is(err, blog.ErrEmptyTitle) || got != (blog.Post{}) {
			t.Errorf("NewPost(%q) = %+v, %v, want {}, %v", title, got, err, blog.ErrEmptyTitle)
		}
	}
}

func TestURL(t *testing.T) {
	tests := []struct {
		p    blog.Post
		want string
	}{
		{blog.Post{Title: "Hello", Slug: "hello-world"}, "/posts/hello-world"},
		{blog.Post{Slug: "go-2"}, "/posts/go-2"},
		{blog.Post{}, "/posts/"},
	}
	for _, tt := range tests {
		if got := tt.p.URL(); got != tt.want {
			t.Errorf("%+v.URL() = %q, want %q", tt.p, got, tt.want)
		}
	}
}

func TestNewPosts(t *testing.T) {
	tests := []struct {
		in    []string
		slugs []string
	}{
		{nil, []string{}},
		{[]string{"One", "Two"}, []string{"one", "two"}},
		{[]string{"Go", "go!", "GO?"}, []string{"go", "go-2", "go-3"}},
		{[]string{"a-2", "a", "a"}, []string{"a-2", "a", "a-3"}},
		{[]string{"x", "x", "x-2"}, []string{"x", "x-2", "x-2-2"}},
	}
	for _, tt := range tests {
		got, err := blog.NewPosts(tt.in)
		if err != nil {
			t.Errorf("NewPosts(%q) error = %v, want <nil>", tt.in, err)
			continue
		}
		slugs := []string{}
		for i, p := range got {
			slugs = append(slugs, p.Slug)
			if i < len(tt.in) && p.Title != tt.in[i] {
				t.Errorf("NewPosts(%q)[%d].Title = %q, want %q", tt.in, i, p.Title, tt.in[i])
			}
		}
		if len(slugs) != len(tt.slugs) {
			t.Errorf("NewPosts(%q) slugs = %q, want %q", tt.in, slugs, tt.slugs)
			continue
		}
		for i := range slugs {
			if slugs[i] != tt.slugs[i] {
				t.Errorf("NewPosts(%q) slugs = %q, want %q", tt.in, slugs, tt.slugs)
				break
			}
		}
	}
}

func TestNewPostsError(t *testing.T) {
	in := []string{"ok", "fine", "???"}
	got, err := blog.NewPosts(in)
	if !errors.Is(err, blog.ErrEmptyTitle) || got != nil {
		t.Fatalf("NewPosts(%q) = %v, %v, want nil, %v", in, got, err, blog.ErrEmptyTitle)
	}
	if want := "title 2: " + blog.ErrEmptyTitle.Error(); err.Error() != want {
		t.Errorf("NewPosts(%q) error = %q, want %q", in, err.Error(), want)
	}
}
