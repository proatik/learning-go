package library

import "errors"

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
	// TODO
	return nil
}

// AddBook adds b to the catalogue, or adds copies to an existing ISBN.
func AddBook(lib *Library, b Book) error {
	// TODO
	return nil
}

// AddMember registers a new member.
func AddMember(lib *Library, id, name string) error {
	// TODO
	return nil
}

// Borrow lends one copy of isbn to member memberID.
func Borrow(lib *Library, memberID, isbn string) error {
	// TODO
	return nil
}

// Return takes back a copy of isbn from member memberID.
func Return(lib *Library, memberID, isbn string) error {
	// TODO
	return nil
}

// Catalog returns copies of the books by author (all books if author is
// ""), sorted by title and then ISBN.
func Catalog(lib *Library, author string) []Book {
	// TODO
	return nil
}

// Borrowers returns the sorted IDs of members currently holding isbn.
func Borrowers(lib *Library, isbn string) []string {
	// TODO
	return nil
}
