package ddb

import (
	"errors"
	"time"

	ddbv1 "github.com/danielfsousa/ddb/gen/ddb/v1"
	"github.com/danielfsousa/ddb/internal/commitlog"
	"github.com/danielfsousa/ddb/internal/config"
)

var (
	// ErrKeyNotFound is the error returned when a key is not found.
	ErrKeyNotFound = errors.New("key not found")

	// ErrKeyEmpty is the error returned when a key is empty.
	ErrKeyEmpty = errors.New("key cannot be empty")

	// ErrKeyTooLarge is the error returned when a key is too large.
	ErrKeyTooLarge = errors.New("key is too large")

	// ErrValueTooLarge is the error returned when a value is too large.
	ErrValueTooLarge = errors.New("value is too large")
)

// Ddb is a distributed key-value store consisting of a commit log and an in-memory index hash map.
type Ddb struct {
	config *config.Config
	log    *commitlog.CommitLog
	dir    string
}

// Open opens a new Ddb instance at the given directory.
func Open(dir string, options ...Option) (*Ddb, error) {
	cfg := config.NewDefaultConfig()
	for _, option := range options {
		if err := option(cfg); err != nil {
			return nil, err
		}
	}

	log, err := commitlog.NewCommitLog(dir, commitlog.Config{})
	if err != nil {
		return nil, err
	}

	return &Ddb{
		config: cfg,
		log:    log,
		dir:    dir,
	}, nil
}

// Has returns true if the given key exists in the database.
func (d *Ddb) Has(key string) bool {
	meta, exists := d.log.GetMetadata(key)
	return exists && meta.DeletedAt == nil
}

// Get retrieves the value for the given key.
func (d *Ddb) Get(key string) ([]byte, error) {
	if !d.Has(key) {
		return nil, ErrKeyNotFound
	}
	rec, exists, err := d.log.Get(key)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrKeyNotFound
	}
	return rec.Value, nil
}

// Set sets the value for the given key.
func (d *Ddb) Set(key string, val []byte) error {
	if key == "" {
		return ErrKeyEmpty
	}
	if uint64(len(key)) > d.config.MaxKeySize {
		return ErrKeyTooLarge
	}
	if uint64(len(val)) > d.config.MaxValueSize {
		return ErrValueTooLarge
	}

	rec := &ddbv1.Record{
		Key:   key,
		Value: val,
	}
	return d.log.Append(rec)
}

// Delete deletes the value for the given key.
func (d *Ddb) Delete(key string) error {
	if !d.Has(key) {
		return ErrKeyNotFound
	}
	t := time.Now().Unix()
	rec := &ddbv1.Record{
		Key:       key,
		DeletedAt: &t,
	}
	return d.log.Append(rec)
}

// Sync flushes all buffers to disk, ensuring that all writes persisted.
// func (s *Ddb) Sync() error {
// 	return s.log.Sync()
// }

// Close closes the Ddb instance.
func (d *Ddb) Close() error {
	return d.log.Close()
}

// Statistics represents statistics about the Ddb instance.
type Statistics struct {
	Segments int
	Keys     int
	Size     int
}

// func (d *Ddb) Stats() (stats *Statistics, err error) {
// 	// TODO: implement function to return the size of the datafiles directory
// 	stats.Segments =
// 	stats.Keys =
// }
