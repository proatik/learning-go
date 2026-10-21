package config

// Config holds optional settings. A nil field means "not set: use the default".
type Config struct {
	Retries *int    // default 3
	Cache   *bool   // default true
	Name    *string // default "app"
}

// Settings holds the final, effective values.
type Settings struct {
	Retries int
	Cache   bool
	Name    string
}

// Effective returns the settings described by c, using the default for
// every nil field.
func Effective(c Config) Settings {
	// TODO
	return Settings{}
}

// Merge returns a Config whose fields come from override when they are
// non-nil and from base otherwise.
func Merge(base, override Config) Config {
	// TODO
	return Config{}
}

// FromSettings returns a Config with every field set to the value in s.
func FromSettings(s Settings) Config {
	// TODO
	return Config{}
}
