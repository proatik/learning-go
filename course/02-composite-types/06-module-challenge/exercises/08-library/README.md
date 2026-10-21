# Library System

Model a lending library: books and members are stored by key as pointers, so every function sees and updates the same records.

```go
const MaxLoans = 3

type Book struct {
	ISBN      string
	Title     string
	Author    string
	Copies    int // copies the library owns
	Available int // copies currently on the shelf
}

type Member struct {
	ID    string
	Name  string
	Loans []string // ISBNs currently borrowed, oldest first
}

type Library struct {
	Books   map[string]*Book   // by ISBN
	Members map[string]*Member // by ID
}

var (
	ErrInvalid         = errors.New("library: invalid argument")
	ErrDuplicate       = errors.New("library: duplicate member")
	ErrNoBook          = errors.New("library: no such book")
	ErrNoMember        = errors.New("library: no such member")
	ErrAlreadyBorrowed = errors.New("library: book already borrowed by member")
	ErrLimit           = errors.New("library: loan limit reached")
	ErrUnavailable     = errors.New("library: no copies available")
	ErrNotBorrowed     = errors.New("library: book not borrowed by member")
)

func NewLibrary() *Library
func AddBook(lib *Library, b Book) error
func AddMember(lib *Library, id, name string) error
func Borrow(lib *Library, memberID, isbn string) error
func Return(lib *Library, memberID, isbn string) error
func Catalog(lib *Library, author string) []Book
func Borrowers(lib *Library, isbn string) []string
```

## Example

```go
lib := NewLibrary()
AddBook(lib, Book{ISBN: "111", Title: "The Go Programming Language", Author: "Donovan", Copies: 1})
AddMember(lib, "m1", "Ada")
AddMember(lib, "m2", "Bob")
Borrow(lib, "m1", "111")   // nil — Available 0, m1.Loans = [111]
Borrow(lib, "m2", "111")   // ErrUnavailable
Borrowers(lib, "111")      // [m1]
Return(lib, "m1", "111")   // nil — Available 1, m1.Loans = []
```

## Rules
- `NewLibrary` returns a library with both maps created. `lib` is never nil afterwards.
- **AddBook**: `b.ISBN == ""` or `b.Copies < 1` → `ErrInvalid`. A **new** ISBN is stored with `Available = Copies` (the caller's `Available` is ignored). An **existing** ISBN gains `b.Copies` more copies (both `Copies` and `Available` increase); its title and author stay unchanged. Store your own `*Book`, not a pointer to the caller's value.
- **AddMember**: `id == ""` → `ErrInvalid`; existing ID → `ErrDuplicate`.
- **Borrow** checks, in order: unknown member → `ErrNoMember`; unknown book → `ErrNoBook`; member already has this ISBN → `ErrAlreadyBorrowed`; member has `MaxLoans` loans → `ErrLimit`; `Available == 0` → `ErrUnavailable`. On success decrement `Available` and append the ISBN to `Loans`.
- **Return** checks, in order: `ErrNoMember`, `ErrNoBook`, ISBN not in the member's loans → `ErrNotBorrowed`. On success increment `Available` and remove the ISBN from `Loans`, keeping the order of the others.
- A failed operation changes nothing.
- **Catalog** returns **copies** (`Book` values) of the books by `author` (exact match), or of every book when `author == ""`, sorted by `Title`, then `ISBN`. Modifying the result must not affect the library. No matches → empty slice.
- **Borrowers** returns the IDs of members currently holding `isbn`, sorted ascending (empty if none or unknown).
- Return the sentinels (wrapping with `%w` is fine); tests use `errors.Is`.
