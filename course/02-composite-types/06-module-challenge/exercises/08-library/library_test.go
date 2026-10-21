package library

import (
	"errors"
	"slices"
	"testing"
)

func newLib(t *testing.T) *Library {
	t.Helper()
	lib := NewLibrary()
	if lib == nil || lib.Books == nil || lib.Members == nil {
		t.Fatalf("NewLibrary() = %+v, want a library with both maps created", lib)
	}
	return lib
}

// setup builds a library with four books and three members.
func setup(t *testing.T) *Library {
	t.Helper()
	lib := newLib(t)
	books := []Book{
		{ISBN: "111", Title: "The Go Programming Language", Author: "Donovan", Copies: 1},
		{ISBN: "222", Title: "Learning Go", Author: "Bodner", Copies: 2},
		{ISBN: "333", Title: "Concurrency in Go", Author: "Cox-Buday", Copies: 1},
		{ISBN: "444", Title: "100 Go Mistakes", Author: "Harsanyi", Copies: 3},
		{ISBN: "555", Title: "Efficient Go", Author: "Plotka", Copies: 1},
	}
	for _, b := range books {
		if err := AddBook(lib, b); err != nil {
			t.Fatalf("AddBook(%q) = %v, want nil", b.ISBN, err)
		}
	}
	for _, m := range [][2]string{{"m1", "Ada"}, {"m2", "Bob"}, {"m3", "Cy"}} {
		if err := AddMember(lib, m[0], m[1]); err != nil {
			t.Fatalf("AddMember(%q) = %v, want nil", m[0], err)
		}
	}
	return lib
}

func mustBook(t *testing.T, lib *Library, isbn string) *Book {
	t.Helper()
	b, ok := lib.Books[isbn]
	if !ok || b == nil {
		t.Fatalf("Books[%q] missing", isbn)
	}
	return b
}

func mustMember(t *testing.T, lib *Library, id string) *Member {
	t.Helper()
	m, ok := lib.Members[id]
	if !ok || m == nil {
		t.Fatalf("Members[%q] missing", id)
	}
	return m
}

func TestAddBook(t *testing.T) {
	lib := newLib(t)
	b := Book{ISBN: "111", Title: "Go", Author: "Rob", Copies: 2, Available: 99}
	if err := AddBook(lib, b); err != nil {
		t.Fatalf("AddBook = %v", err)
	}
	got := mustBook(t, lib, "111")
	if want := (Book{"111", "Go", "Rob", 2, 2}); *got != want {
		t.Fatalf("Books[111] = %+v, want %+v", *got, want)
	}
	b.Title = "changed" // the caller's copy must not be stored
	if got.Title != "Go" {
		t.Fatalf("AddBook stored a pointer to the caller's value")
	}

	if err := AddBook(lib, Book{ISBN: "111", Title: "Other", Author: "X", Copies: 3}); err != nil {
		t.Fatalf("AddBook(existing) = %v", err)
	}
	got = mustBook(t, lib, "111")
	if want := (Book{"111", "Go", "Rob", 5, 5}); *got != want {
		t.Fatalf("after adding copies Books[111] = %+v, want %+v", *got, want)
	}

	for _, bad := range []Book{
		{ISBN: "", Title: "No ISBN", Copies: 1},
		{ISBN: "999", Title: "Zero", Copies: 0},
		{ISBN: "999", Title: "Negative", Copies: -2},
		{ISBN: "111", Copies: 0},
	} {
		if err := AddBook(lib, bad); !errors.Is(err, ErrInvalid) {
			t.Errorf("AddBook(%+v) = %v, want ErrInvalid", bad, err)
		}
	}
	if len(lib.Books) != 1 || lib.Books["111"].Copies != 5 {
		t.Errorf("failed AddBook calls changed the library")
	}
}

func TestAddMember(t *testing.T) {
	lib := newLib(t)
	if err := AddMember(lib, "m1", "Ada"); err != nil {
		t.Fatalf("AddMember(m1) = %v", err)
	}
	m := mustMember(t, lib, "m1")
	if m.ID != "m1" || m.Name != "Ada" || len(m.Loans) != 0 {
		t.Fatalf("Members[m1] = %+v, want {m1 Ada []}", *m)
	}
	if err := AddMember(lib, "m1", "Imposter"); !errors.Is(err, ErrDuplicate) {
		t.Errorf("AddMember(duplicate) = %v, want ErrDuplicate", err)
	}
	if lib.Members["m1"].Name != "Ada" {
		t.Errorf("duplicate AddMember overwrote the member")
	}
	if err := AddMember(lib, "", "Nobody"); !errors.Is(err, ErrInvalid) {
		t.Errorf("AddMember(\"\") = %v, want ErrInvalid", err)
	}
}

func TestBorrowReturn(t *testing.T) {
	lib := setup(t)
	if err := Borrow(lib, "m1", "111"); err != nil {
		t.Fatalf("Borrow(m1, 111) = %v", err)
	}
	if got := mustBook(t, lib, "111").Available; got != 0 {
		t.Errorf("Available(111) = %d, want 0", got)
	}
	if got := mustMember(t, lib, "m1").Loans; !slices.Equal(got, []string{"111"}) {
		t.Errorf("m1.Loans = %v, want [111]", got)
	}
	if err := Borrow(lib, "m2", "111"); !errors.Is(err, ErrUnavailable) {
		t.Errorf("Borrow(m2, 111) = %v, want ErrUnavailable", err)
	}
	if err := Return(lib, "m1", "111"); err != nil {
		t.Fatalf("Return(m1, 111) = %v", err)
	}
	if got := mustBook(t, lib, "111").Available; got != 1 {
		t.Errorf("Available(111) after return = %d, want 1", got)
	}
	if got := mustMember(t, lib, "m1").Loans; len(got) != 0 {
		t.Errorf("m1.Loans after return = %v, want empty", got)
	}
	if err := Borrow(lib, "m2", "111"); err != nil {
		t.Errorf("Borrow(m2, 111) after return = %v, want nil", err)
	}
}

func TestReturnKeepsOrder(t *testing.T) {
	lib := setup(t)
	for _, isbn := range []string{"222", "333", "444"} {
		if err := Borrow(lib, "m3", isbn); err != nil {
			t.Fatalf("Borrow(m3, %s) = %v", isbn, err)
		}
	}
	if err := Return(lib, "m3", "333"); err != nil {
		t.Fatalf("Return(m3, 333) = %v", err)
	}
	if got := mustMember(t, lib, "m3").Loans; !slices.Equal(got, []string{"222", "444"}) {
		t.Fatalf("m3.Loans = %v, want [222 444]", got)
	}
	if err := Borrow(lib, "m3", "555"); err != nil {
		t.Fatalf("Borrow(m3, 555) = %v", err)
	}
	if got := mustMember(t, lib, "m3").Loans; !slices.Equal(got, []string{"222", "444", "555"}) {
		t.Fatalf("m3.Loans = %v, want [222 444 555]", got)
	}
}

func TestBorrowErrors(t *testing.T) {
	lib := setup(t)
	for _, isbn := range []string{"222", "333", "444"} {
		if err := Borrow(lib, "m1", isbn); err != nil {
			t.Fatalf("Borrow(m1, %s) = %v", isbn, err)
		}
	}
	if err := Borrow(lib, "m2", "333"); !errors.Is(err, ErrUnavailable) { // 333 has 1 copy, m1 has it
		t.Fatalf("Borrow(m2, 333) = %v, want ErrUnavailable", err)
	}
	snapshot := Catalog(lib, "")
	tests := []struct {
		member, isbn string
		want         error
	}{
		{"zz", "111", ErrNoMember},
		{"zz", "zz", ErrNoMember},
		{"m2", "zz", ErrNoBook},
		{"m1", "222", ErrAlreadyBorrowed},
		{"m1", "333", ErrAlreadyBorrowed}, // before ErrLimit and ErrUnavailable
		{"m1", "111", ErrLimit},
		{"m1", "555", ErrLimit},
		{"m2", "333", ErrUnavailable},
	}
	for _, tt := range tests {
		if err := Borrow(lib, tt.member, tt.isbn); !errors.Is(err, tt.want) {
			t.Errorf("Borrow(%q, %q) = %v, want %v", tt.member, tt.isbn, err, tt.want)
		}
	}
	if got := Catalog(lib, ""); !slices.Equal(got, snapshot) {
		t.Errorf("failed Borrow calls changed the books:\n got %v\nwant %v", got, snapshot)
	}
	if got := mustMember(t, lib, "m1").Loans; !slices.Equal(got, []string{"222", "333", "444"}) {
		t.Errorf("failed Borrow calls changed m1.Loans: %v", got)
	}
	if got := mustMember(t, lib, "m2").Loans; len(got) != 0 {
		t.Errorf("failed Borrow calls changed m2.Loans: %v", got)
	}
}

func TestReturnErrors(t *testing.T) {
	lib := setup(t)
	if err := Borrow(lib, "m1", "222"); err != nil {
		t.Fatalf("Borrow(m1, 222) = %v", err)
	}
	tests := []struct {
		member, isbn string
		want         error
	}{
		{"zz", "222", ErrNoMember},
		{"m1", "zz", ErrNoBook},
		{"m2", "222", ErrNotBorrowed},
		{"m1", "111", ErrNotBorrowed},
	}
	for _, tt := range tests {
		if err := Return(lib, tt.member, tt.isbn); !errors.Is(err, tt.want) {
			t.Errorf("Return(%q, %q) = %v, want %v", tt.member, tt.isbn, err, tt.want)
		}
	}
	if got := mustBook(t, lib, "222").Available; got != 1 {
		t.Errorf("failed Return calls changed Available(222) to %d, want 1", got)
	}
	if err := Return(lib, "m1", "222"); err != nil {
		t.Fatalf("Return(m1, 222) = %v", err)
	}
	if err := Return(lib, "m1", "222"); !errors.Is(err, ErrNotBorrowed) {
		t.Errorf("second Return(m1, 222) = %v, want ErrNotBorrowed", err)
	}
	if got := mustBook(t, lib, "222").Available; got != 2 {
		t.Errorf("Available(222) = %d, want 2 (never above Copies)", got)
	}
}

func TestCatalog(t *testing.T) {
	lib := setup(t)
	_ = AddBook(lib, Book{ISBN: "000", Title: "Learning Go", Author: "Someone", Copies: 1}) // same title
	_ = Borrow(lib, "m1", "444")

	got := Catalog(lib, "")
	want := []Book{
		{"444", "100 Go Mistakes", "Harsanyi", 3, 2},
		{"333", "Concurrency in Go", "Cox-Buday", 1, 1},
		{"555", "Efficient Go", "Plotka", 1, 1},
		{"000", "Learning Go", "Someone", 1, 1},
		{"222", "Learning Go", "Bodner", 2, 2},
		{"111", "The Go Programming Language", "Donovan", 1, 1},
	}
	if !slices.Equal(got, want) {
		t.Fatalf("Catalog(\"\") =\n %v\nwant\n %v", got, want)
	}

	got[0].Available = 100 // must be a copy
	if lib.Books["444"].Available != 2 {
		t.Fatalf("modifying the Catalog result changed the library")
	}

	if got := Catalog(lib, "Bodner"); !slices.Equal(got, []Book{{"222", "Learning Go", "Bodner", 2, 2}}) {
		t.Errorf("Catalog(Bodner) = %v", got)
	}
	if got := Catalog(lib, "bodner"); len(got) != 0 {
		t.Errorf("Catalog(bodner) = %v, want empty (exact match)", got)
	}
	if got := Catalog(newLib(t), ""); len(got) != 0 {
		t.Errorf("Catalog(empty library) = %v, want empty", got)
	}
}

func TestBorrowers(t *testing.T) {
	lib := setup(t)
	for _, m := range []string{"m3", "m1", "m2"} {
		if err := Borrow(lib, m, "444"); err != nil {
			t.Fatalf("Borrow(%s, 444) = %v", m, err)
		}
	}
	_ = Borrow(lib, "m2", "222")
	tests := []struct {
		isbn string
		want []string
	}{
		{"444", []string{"m1", "m2", "m3"}},
		{"222", []string{"m2"}},
		{"111", []string{}},
		{"nope", []string{}},
	}
	for _, tt := range tests {
		if got := Borrowers(lib, tt.isbn); !slices.Equal(got, tt.want) {
			t.Errorf("Borrowers(%q) = %v, want %v", tt.isbn, got, tt.want)
		}
	}
	_ = Return(lib, "m2", "444")
	if got := Borrowers(lib, "444"); !slices.Equal(got, []string{"m1", "m3"}) {
		t.Errorf("Borrowers(444) after return = %v, want [m1 m3]", got)
	}
}
