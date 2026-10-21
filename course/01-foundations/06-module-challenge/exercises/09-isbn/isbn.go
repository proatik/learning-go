package isbn

// IsValid reports whether isbn is a valid ISBN-10. Hyphens are ignored.
func IsValid(isbn string) bool {
	// TODO
	return false
}

// ToISBN13 converts a valid ISBN-10 to its 13-digit ISBN-13 form
// (978 prefix, no hyphens).
func ToISBN13(isbn string) (string, error) {
	// TODO
	return "", nil
}
