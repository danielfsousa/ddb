package server

const (
	// DefaultAddr is the host to bind to if one is not specified.
	DefaultHost = "localhost"
	// DefaultPort is the port to bind to if one is not specified.
	DefaultPort = 9191
)

type Config struct {
	Host string
	Port int
}

// NewDefaultConfig creates a new Config with default settings.
func NewDefaultConfig() *Config {
	return &Config{
		Host: DefaultHost,
		Port: DefaultPort,
	}
}
