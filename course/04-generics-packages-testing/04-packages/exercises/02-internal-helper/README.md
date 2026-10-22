# Internal Helper

Write a slug helper in an `internal/` package and build the public `blog` API on top of it.

```go
// package slug  (internal/slug/slug.go): importable only from this exercise's tree
func Make(title string) string

// package blog  (blog.go): the public API
var ErrEmptyTitle = errors.New("blog: title has no letters or digits")

type Post struct {
	Title string
	Slug  string
}

func NewPost(title string) (Post, error)
func (p Post) URL() string
func NewPosts(titles []string) ([]Post, error)
```

## Examples

| Call | Result |
|---|---|
| `slug.Make("Hello, World!")` | `"hello-world"` |
| `slug.Make("  Go 1.27 -- Released  ")` | `"go-1-27-released"` |
| `slug.Make("Ünïcode Ok")` | `"ünïcode-ok"` |
| `slug.Make("!!!")` | `""` |
| `NewPost("Hello, World!")` | `Post{"Hello, World!", "hello-world"}, nil` |
| `NewPost("?!")` | error `ErrEmptyTitle` |
| `Post{Slug: "hello-world"}.URL()` | `"/posts/hello-world"` |
| `NewPosts([]string{"Go", "go!", "GO?"})` | slugs `"go"`, `"go-2"`, `"go-3"` |

## Rules
- `Make`:
  - lowercases every rune (`unicode.ToLower`) and keeps letters and digits (`unicode.IsLetter`, `unicode.IsDigit`);
  - replaces every **run** of other characters with a single `-`;
  - never starts or ends with `-`.
- `NewPost` keeps the title as given and fills `Slug` with `slug.Make(title)`. An empty slug → `Post{}` and `ErrEmptyTitle`.
- `URL` returns `"/posts/" + Slug`.
- `NewPosts` builds one post per title, in order. A slug that's already taken gets the suffix `-N`, with the smallest `N ≥ 2` that is still free.
  - Example: `"a-2"`, `"a"`, `"a"` → `"a-2"`, `"a"`, `"a-3"`.
- On an empty slug, `NewPosts` returns `nil` and an error that wraps `ErrEmptyTitle` and names the index: `fmt.Errorf("title %d: %w", i, ErrEmptyTitle)`.
- The tests import only `blog`: `internal/slug` is an implementation detail.
