package appconf

// Config is the application's configuration.
type Config struct {
	Host string
	Port int
}

// App loads its Config lazily: on the first call to Config, not in New.
type App struct {
	config func() (Config, error) // set by New; runs the loader at most once
}

// New returns an App that will get its configuration from load. New must not
// call load.
func New(load func() (Config, error)) *App {
	// TODO
	return &App{}
}

// Config returns the configuration, calling the loader on the first call
// only. Every call, from any goroutine, gets the same Config and error.
func (a *App) Config() (Config, error) {
	// TODO
	return Config{}, nil
}

// Addr returns "host:port" from the configuration, or the loader's error
// wrapped as "appconf: <err>".
func (a *App) Addr() (string, error) {
	// TODO
	return "", nil
}
