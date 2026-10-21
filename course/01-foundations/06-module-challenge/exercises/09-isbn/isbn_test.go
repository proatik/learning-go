package isbn

import "testing"

func TestIsValid(t *testing.T) {
	tests := []struct {
		isbn string
		want bool
	}{
		{"3-598-21508-8", true},
		{"3598215088", true},
		{"3-598-21507-X", true},
		{"359821507X", true},
		{"0-306-40615-2", true},
		{"0306406152", true},
		{"0-19-853453-1", true},
		{"--3598215088--", true}, // hyphens anywhere are ignored
		{"0000000000", true},
		{"3-598-21508-9", false},   // wrong check digit
		{"3-598-21507-A", false},   // invalid check character
		{"3-598-21507-x", false},   // lowercase x
		{"3-598-P1581-X", false},   // letter in the body
		{"3-598-2X507-9", false},   // X not in last position
		{"X598215088", false},      // X not in last position
		{"359821507", false},       // too short
		{"3598215078X", false},     // too long
		{"3-598-21507", false},     // missing check character
		{"98245726788", false},     // 11 digits
		{"3 598 21508 8", false},   // spaces are not ignored
		{"3132P34035", false},      // letter in the body
		{"", false},                // empty
		{"----------", false},      // only hyphens
		{"3-598-21508-8\n", false}, // trailing newline
		{"３598215088", false},      // full-width digit
	}
	for _, tt := range tests {
		if got := IsValid(tt.isbn); got != tt.want {
			t.Errorf("IsValid(%q) = %t, want %t", tt.isbn, got, tt.want)
		}
	}
}

func TestToISBN13(t *testing.T) {
	tests := []struct {
		isbn, want string
	}{
		{"0-306-40615-2", "9780306406157"},
		{"3-598-21508-8", "9783598215087"},
		{"359821507X", "9783598215070"},
		{"0-19-853453-1", "9780198534532"},
		{"0000000000", "9780000000002"},
	}
	for _, tt := range tests {
		got, err := ToISBN13(tt.isbn)
		if err != nil || got != tt.want {
			t.Errorf("ToISBN13(%q) = %q, %v, want %q, nil", tt.isbn, got, err, tt.want)
		}
	}

	for _, s := range []string{"", "3-598-21508-9", "3-598-21507-x", "359821507", "abcdefghij"} {
		if got, err := ToISBN13(s); err == nil || got != "" {
			t.Errorf("ToISBN13(%q) = %q, %v, want \"\", error", s, got, err)
		}
	}
}
