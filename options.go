package ddb

import "github.com/danielfsousa/ddb/internal/config"

// Option is a function that takes a config and modifies it.
type Option func(*config.Config) error

func withConfig(src *config.Config) Option {
	return func(dst *config.Config) error {
		*dst = *src
		return nil
	}
}

// WithMaxKeySize sets the maximum key size.
func WithMaxKeySize(size uint64) Option {
	return withConfig(&config.Config{
		MaxKeySize: size,
	})
}

// WithMaxValueSize sets the maximum value size.
func WithMaxValueSize(size uint64) Option {
	return withConfig(&config.Config{
		MaxValueSize: size,
	})
}

// WithMaxSegmentDataSize sets the maximum datafile size option
func WithMaxSegmentDataSize(size uint64) Option {
	return func(cfg *config.Config) error {
		cfg.MaxSegmentDataSize = size
		return nil
	}
}
