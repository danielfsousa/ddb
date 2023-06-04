package bitcask

import (
	"io"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"

	ddbv1 "github.com/danielfsousa/ddb/gen/ddb/v1"
	"github.com/danielfsousa/ddb/internal/backend"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type Bitcask struct {
	mu sync.RWMutex

	Dir    string
	Config Config

	activeSegment *segment
	segments      []*segment
}

var _ backend.Backend = (*Bitcask)(nil)

// NewBitcaskBackend creates a new Bitcask backend.
func NewBitcaskBackend(dir string, c Config) (*Bitcask, error) {
	if c.Segment.MaxStoreBytes == 0 {
		c.Segment.MaxStoreBytes = 2e+9 // 2GB
	}
	if c.Segment.MaxIndexBytes == 0 {
		c.Segment.MaxIndexBytes = 1e+8 // 100MB
	}
	bitcask := &Bitcask{
		Dir:    dir,
		Config: c,
	}
	err := bitcask.setup()
	return bitcask, err
}

func (b *Bitcask) setup() error {
	files, err := os.ReadDir(b.Dir)
	if err != nil {
		return err
	}

	var ids []uint64
	for _, file := range files {
		idStr := strings.TrimSuffix(file.Name(), path.Ext(file.Name()))
		id, err := strconv.ParseUint(idStr, 10, 0)
		if err != nil {
			return err
		}
		ids = append(ids, id)
	}
	slices.Sort(ids)

	for i := 0; i < len(ids); i++ {
		if err = b.newSegment(ids[i]); err != nil {
			return err
		}
		i += 1 // skip hint files
	}
	if b.segments == nil {
		if err = b.newSegment(1); err != nil {
			return err
		}
	}

	return nil
}

func (b *Bitcask) newSegment(id uint64) error {
	s, err := newSegment(b.Dir, id, b.Config)
	if err != nil {
		return err
	}
	b.segments = append(b.segments, s)
	b.activeSegment = s
	return nil
}

// Keys returns a slice of the keys of all records stored in the log.
func (b *Bitcask) Keys() []string {
	b.mu.RLock()
	defer b.mu.RUnlock()
	keys := make(map[string]bool)
	for _, segment := range b.segments {
		for _, key := range segment.Keys() {
			keys[key] = true
		}
	}
	keysSlice := maps.Keys(keys)
	slices.Sort(keysSlice)
	return keysSlice
}

// Has returns true if the key exists in the log.
func (b *Bitcask) Has(key string) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for i := len(b.segments) - 1; i >= 0; i-- {
		if b.segments[i].Has(key) {
			return true
		}
	}
	return false
}

// Get returns a record by key.
func (b *Bitcask) Get(key string) (rec *ddbv1.Record, exists bool, err error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for i := len(b.segments) - 1; i >= 0; i-- {
		rec, exists, err = b.segments[i].Get(key)
		if err != nil {
			return nil, false, err
		}
		if !exists {
			continue
		}
		return rec, true, nil
	}
	return nil, false, nil
}

// GetMetadata returns the metadata for a key.
func (b *Bitcask) GetMetadata(key string) (entry backend.RecordMetadata, exists bool) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for i := len(b.segments) - 1; i >= 0; i-- {
		meta, exists := b.segments[i].GetMetadata(key)
		if exists {
			return meta, true
		}
	}
	return backend.RecordMetadata{}, false
}

// Set appends a record to the log and updates the in-memory index.
func (b *Bitcask) Set(rec *ddbv1.Record) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if err := b.activeSegment.Append(rec); err != nil {
		return err
	}
	if b.activeSegment.IsMaxed() {
		if err := b.newSegment(b.activeSegment.id + 1); err != nil {
			return err
		}
	}
	return nil
}

// Sync flushes all pending log writes to disk.
func (b *Bitcask) Sync() error {
	return b.activeSegment.Sync()
}

// Close closes the log.
func (b *Bitcask) Close() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, segment := range b.segments {
		if err := segment.Close(); err != nil {
			return err
		}
	}
	return nil
}

// Remove closes the log and then removes its data from the filesystem.
func (b *Bitcask) Remove() error {
	if err := b.Close(); err != nil {
		return err
	}
	return os.RemoveAll(b.Dir)
}

// Reset removes the log and re-creates it.
func (b *Bitcask) Reset() error {
	if err := b.Remove(); err != nil {
		return err
	}
	return b.setup()
}

// Reader returns an io.Reader instance to read the whole log.
func (b *Bitcask) Reader() io.Reader {
	b.mu.RLock()
	defer b.mu.RUnlock()
	readers := make([]io.Reader, len(b.segments))
	for i, segment := range b.segments {
		readers[i] = &originReader{segment.store, 0}
	}
	return io.MultiReader(readers...)
}

type originReader struct {
	*store
	off int64
}

func (r *originReader) Read(p []byte) (n int, err error) {
	n, err = r.ReadAt(p, r.off)
	r.off += int64(n)
	return n, err
}
