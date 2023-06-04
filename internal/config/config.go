package config

const (
	// DefaultMaxDatafileSize is the default maximum datafile size in bytes
	DefaultMaxDatafileSize = 1 << 20 // 1MB

	// DefaultMaxKeySize is the default maximum key size in bytes
	DefaultMaxKeySize = uint64(64) // 64 bytes

	// DefaultMaxValueSize is the default value size in bytes
	DefaultMaxValueSize = uint64(1 << 16) // 65KB
)

type Config struct {
	MaxKeySize         uint64
	MaxValueSize       uint64
	MaxSegmentDataSize uint64
}

// NewDefaultConfig creates a new Config with default settings.
func NewDefaultConfig() *Config {
	return &Config{
		MaxKeySize:         DefaultMaxKeySize,
		MaxValueSize:       DefaultMaxValueSize,
		MaxSegmentDataSize: DefaultMaxDatafileSize,
	}
}
