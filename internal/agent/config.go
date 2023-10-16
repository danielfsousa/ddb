package agent

const (
	// DefaultBindAddr is the address to bind Serf on if one is not specified.
	DefaultBindAddr = "localhost:8401"
	// DefaultRPCPort is the port for RPC clients (and Raft) connections to bind to if one is not specified.
	DefaultRPCPort = 9191
)

type Config struct {
	DataDir        string
	BindAddr       string
	RPCPort        int
	NodeName       string
	StartJoinAddrs []string
	Bootstrap      bool
}

// NewDefaultConfig creates a new Config with default settings.
func NewDefaultConfig() *Config {
	return &Config{
		BindAddr: DefaultBindAddr,
		RPCPort:  DefaultRPCPort,
	}
}
