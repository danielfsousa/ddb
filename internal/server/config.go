package server

import "github.com/danielfsousa/ddb"

const (
	// DefaultAddr is the host to bind to if one is not specified.
	DefaultHost = "localhost"
	// DefaultPort is the port to bind to if one is not specified.
	DefaultPort = 9191
)

type Config struct {
	Host string
	Port int
	Ddb  *ddb.Ddb
}

// NewDefaultConfig creates a new Config with default settings.
func NewDefaultConfig() *Config {
	return &Config{
		Host: DefaultHost,
		Port: DefaultPort,
	}
}
