package cut

import "testing"

func TestParseKeyValue(t *testing.T) {
	tests := []struct {
		s          string
		key, value string
	}{
		{"name=Go", "name", "Go"},
		{"  lang = Go 1.27 ", "lang", "Go 1.27"},
		{"url=a=b", "url", "a=b"},
		{"empty=", "empty", ""},
		{"ville=Montréal", "ville", "Montréal"},
		{"城市=東京", "城市", "東京"},
		{"\tk\t=\tv\n", "k", "v"},
	}
	for _, tt := range tests {
		key, value, err := ParseKeyValue(tt.s)
		if err != nil || key != tt.key || value != tt.value {
			t.Errorf("ParseKeyValue(%q) = (%q, %q, %v), want (%q, %q, nil)",
				tt.s, key, value, err, tt.key, tt.value)
		}
	}
}

func TestParseKeyValueErrors(t *testing.T) {
	tests := []struct {
		s    string
		want string
	}{
		{"no separator", `missing '=' in "no separator"`},
		{"", `missing '=' in ""`},
		{" = value", `empty key in " = value"`},
		{"=", `empty key in "="`},
	}
	for _, tt := range tests {
		key, value, err := ParseKeyValue(tt.s)
		if err == nil {
			t.Errorf("ParseKeyValue(%q) = (%q, %q, nil), want error %q", tt.s, key, value, tt.want)
			continue
		}
		if err.Error() != tt.want {
			t.Errorf("ParseKeyValue(%q) error = %q, want %q", tt.s, err.Error(), tt.want)
		}
		if key != "" || value != "" {
			t.Errorf("ParseKeyValue(%q) = (%q, %q, err), want empty key and value on error", tt.s, key, value)
		}
	}
}

func TestStripScheme(t *testing.T) {
	tests := []struct {
		url    string
		rest   string
		secure bool
	}{
		{"https://go.dev", "go.dev", true},
		{"https://go.dev/doc/", "go.dev/doc/", true},
		{"http://example.com", "example.com", false},
		{"go.dev", "go.dev", false},
		{"ftp://files", "ftp://files", false},
		{"HTTPS://GO.DEV", "HTTPS://GO.DEV", false},
		{"", "", false},
	}
	for _, tt := range tests {
		rest, secure := StripScheme(tt.url)
		if rest != tt.rest || secure != tt.secure {
			t.Errorf("StripScheme(%q) = (%q, %t), want (%q, %t)", tt.url, rest, secure, tt.rest, tt.secure)
		}
	}
}

func TestSplitExt(t *testing.T) {
	tests := []struct {
		filename  string
		name, ext string
	}{
		{"archive.tar.gz", "archive.tar", "gz"},
		{"main.go", "main", "go"},
		{"README", "README", ""},
		{"photo.JPG", "photo", "JPG"},
		{"trailing.", "trailing", ""},
		{".bashrc", "", "bashrc"},
		{"résumé.pdf", "résumé", "pdf"},
		{"", "", ""},
	}
	for _, tt := range tests {
		name, ext := SplitExt(tt.filename)
		if name != tt.name || ext != tt.ext {
			t.Errorf("SplitExt(%q) = (%q, %q), want (%q, %q)", tt.filename, name, ext, tt.name, tt.ext)
		}
	}
}
