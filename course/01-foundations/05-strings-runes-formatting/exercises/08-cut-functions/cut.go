package cut

// ParseKeyValue splits s at its first '=' and trims whitespace around the
// key and the value. It returns an error if s has no '=' or the key is empty.
func ParseKeyValue(s string) (key, value string, err error) {
	// TODO: strings.Cut, strings.TrimSpace, fmt.Errorf
	return "", "", nil
}

// StripScheme removes a leading "https://" or "http://" from url and
// reports whether the scheme was https.
func StripScheme(url string) (rest string, secure bool) {
	// TODO: strings.CutPrefix
	return "", false
}

// SplitExt splits filename at its last '.' into name and extension.
// If there is no '.', ext is empty.
func SplitExt(filename string) (name, ext string) {
	// TODO: strings.CutLast
	return "", ""
}
