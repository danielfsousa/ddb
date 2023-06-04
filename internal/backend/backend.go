package backend

import (
	"io"

	ddbv1 "github.com/danielfsousa/ddb/gen/ddb/v1"
)

// RecordMetadata contains metadata about a record.
type RecordMetadata struct {
	Pos       uint64
	DeletedAt *int64
}

// Backend is an interface for a key-value store backend.
type Backend interface {
	Has(key string) bool
	Get(key string) (rec *ddbv1.Record, exists bool, err error)
	GetMetadata(key string) (RecordMetadata, bool)
	Set(rec *ddbv1.Record) error
	Reader() io.Reader
	Sync() error
	Close() error
}
