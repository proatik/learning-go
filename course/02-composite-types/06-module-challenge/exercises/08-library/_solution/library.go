package library

import (
	"cmp"
	"errors"
	"fmt"
	"slices"
)

// MaxLoans is the most books a member may hold at once.
const MaxLoans = 3

// Book is a title in the catalogue with its copy counts.
type Book struct {
	ISBN      string
	Title     string
	Author    string
	Copies    int // copies the library owns
	Available int // copies currently on the shelf
}

// Member is a library member and the ISBNs they currently hold.
type Member struct {
	ID    string
	Name  string
	Loans []string // ISBNs currently borrowed, oldest first
}

// Library stores books by ISBN and members by ID.
type Library struct {
	Books   map[string]*Book   // by ISBN
	Members map[string]*Member // by ID
}

// Errors returned by library operations.
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

// NewLibrary returns an empty library.
func NewLibrary() *Library {
	return &Library{
		Books:   make(map[string]*Book),
		Members: make(map[string]*Member),
	}
}

// AddBook adds b to the catalogue, or adds copies to an existing ISBN.
func AddBook(lib *Library, b Book) error {
	if b.ISBN == "" || b.Copies < 1 {
		return fmt.Errorf("%w: book %q with %d copies", ErrInvalid, b.ISBN, b.Copies)
	}
	if existing, ok := lib.Books[b.ISBN]; ok {
		existing.Copies += b.Copies
		existing.Available += b.Copies
		return nil
	}
	b.Available = b.Copies // b is our own copy of the caller's value
	lib.Books[b.ISBN] = &b
	return nil
}

// AddMember registers a new member.
func AddMember(lib *Library, id, name string) error {
	if id == "" {
		return fmt.Errorf("%w: empty member id", ErrInvalid)
	}
	if _, ok := lib.Members[id]; ok {
		return fmt.Errorf("%w: %q", ErrDuplicate, id)
	}
	lib.Members[id] = &Member{ID: id, Name: name}
	return nil
}

// lookup finds the member and the book, or reports which one is missing.
func lookup(lib *Library, memberID, isbn string) (*Member, *Book, error) {
	m, ok := lib.Members[memberID]
	if !ok {
		return nil, nil, fmt.Errorf("%w: %q", ErrNoMember, memberID)
	}
	b, ok := lib.Books[isbn]
	if !ok {
		return nil, nil, fmt.Errorf("%w: %q", ErrNoBook, isbn)
	}
	return m, b, nil
}

// Borrow lends one copy of isbn to member memberID.
func Borrow(lib *Library, memberID, isbn string) error {
	m, b, err := lookup(lib, memberID, isbn)
	if err != nil {
		return err
	}
	switch {
	case slices.Contains(m.Loans, isbn):
		return fmt.Errorf("%w: %q has %q", ErrAlreadyBorrowed, memberID, isbn)
	case len(m.Loans) >= MaxLoans:
		return fmt.Errorf("%w: %q", ErrLimit, memberID)
	case b.Available == 0:
		return fmt.Errorf("%w: %q", ErrUnavailable, isbn)
	}
	b.Available--
	m.Loans = append(m.Loans, isbn)
	return nil
}

// Return takes back a copy of isbn from member memberID.
func Return(lib *Library, memberID, isbn string) error {
	m, b, err := lookup(lib, memberID, isbn)
	if err != nil {
		return err
	}
	i := slices.Index(m.Loans, isbn)
	if i < 0 {
		return fmt.Errorf("%w: %q does not have %q", ErrNotBorrowed, memberID, isbn)
	}
	m.Loans = slices.Delete(m.Loans, i, i+1) // keeps the order of the rest
	b.Available++
	return nil
}

// Catalog returns copies of the books by author (all books if author is
// ""), sorted by title and then ISBN.
func Catalog(lib *Library, author string) []Book {
	out := []Book{}
	for _, b := range lib.Books {
		if author == "" || b.Author == author {
			out = append(out, *b) // *b copies the struct
		}
	}
	slices.SortFunc(out, func(x, y Book) int {
		return cmp.Or(cmp.Compare(x.Title, y.Title), cmp.Compare(x.ISBN, y.ISBN))
	})
	return out
}

// Borrowers returns the sorted IDs of members currently holding isbn.
func Borrowers(lib *Library, isbn string) []string {
	ids := []string{}
	for id, m := range lib.Members {
		if slices.Contains(m.Loans, isbn) {
			ids = append(ids, id)
		}
	}
	slices.Sort(ids)
	return ids
}
