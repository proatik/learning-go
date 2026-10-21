package cut

import (
	"fmt"
	"strings"
)

// ParseKeyValue splits s at its first '=' and trims whitespace around the
// key and the value. It returns an error if s has no '=' or the key is empty.
func ParseKeyValue(s string) (key, value string, err error) {
	k, v, found := strings.Cut(s, "=") // only the first '=' splits
	if !found {
		return "", "", fmt.Errorf("missing '=' in %q", s)
	}
	k = strings.TrimSpace(k)
	if k == "" {
		return "", "", fmt.Errorf("empty key in %q", s)
	}
	return k, strings.TrimSpace(v), nil
}

// StripScheme removes a leading "https://" or "http://" from url and
// reports whether the scheme was https.
func StripScheme(url string) (rest string, secure bool) {
	if rest, ok := strings.CutPrefix(url, "https://"); ok {
		return rest, true
	}
	if rest, ok := strings.CutPrefix(url, "http://"); ok {
		return rest, false
	}
	return url, false
}

// SplitExt splits filename at its last '.' into name and extension.
// If there is no '.', ext is empty.
func SplitExt(filename string) (name, ext string) {
	// CutLast returns (filename, "", false) when there is no '.',
	// which is exactly the result we want.
	name, ext, _ = strings.CutLast(filename, ".")
	return name, ext
}
