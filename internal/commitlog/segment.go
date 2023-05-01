package commitlog

import (
	"fmt"
	"io"
	"os"
	"path"

	ddbv1 "github.com/danielfsousa/ddb/gen/ddb/v1"
	"github.com/danielfsousa/ddb/pkg/fmode"
)

type segment struct {
	id     uint64
	store  *store
	index  *index
	hint   *hint
	config Config
}

func newSegment(dir string, id uint64, c Config) (*segment, error) {
	s := &segment{
		id:     id,
		config: c,
	}

	var err error
	s.store, err = buildStore(id, dir)
	if err != nil {
		return nil, err
	}

	s.hint, err = buildHint(id, dir)
	if err != nil {
		return nil, err
	}

	s.index, err = buildIndex(s.hint, s.store)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func buildStore(id uint64, dir string) (*store, error) {
	storeFile, err := os.OpenFile(
		path.Join(dir, fmt.Sprintf("%d%s", id, ".store")),
		os.O_RDWR|os.O_CREATE|os.O_APPEND,
		fmode.USER_RW|fmode.GROUP_R|fmode.OTHER_R,
	)
	if err != nil {
		return nil, err
	}
	return newStore(storeFile)
}

func buildHint(id uint64, dir string) (*hint, error) {
	hintFile, err := os.OpenFile(
		path.Join(dir, fmt.Sprintf("%d%s", id, ".hint")),
		os.O_RDWR|os.O_CREATE|os.O_APPEND,
		fmode.USER_RW|fmode.GROUP_R|fmode.OTHER_R,
	)
	if err != nil {
		return nil, err
	}
	return newHint(hintFile)
}

func buildIndex(hint *hint, store *store) (idx *index, err error) {
	idx = newIndex()
	if hint.size > 0 {
		scanner, err := hint.Scanner()
		if err != nil {
			return nil, err
		}
		for scanner.Scan() {
			if err := scanner.Err(); err != nil {
				return nil, err
			}
			key, pos := scanner.Next()
			idx.Set(key, pos)
		}
	} else {
		scanner, err := store.Scanner()
		if err != nil {
			return nil, err
		}
		for scanner.Scan() {
			if err := scanner.Err(); err != nil {
				return nil, err
			}
			rec, pos := scanner.Next()
			idx.Set(rec.Key, pos)
		}
	}
	return idx, nil
}

// Keys returns the keys of all records stored in the segment.
func (s *segment) Keys() []string {
	return s.index.Keys()
}

// Append writes the record to the segment.
func (s *segment) Append(record *ddbv1.Record) error {
	if s.IsMaxed() {
		return io.EOF
	}
	_, pos, err := s.store.Append(record)
	if err != nil {
		return err
	}
	s.index.Set(record.Key, pos)
	return nil
}

// Get returns the record for the given key.
func (s *segment) Get(key string) (rec *ddbv1.Record, exists bool, err error) {
	pos, exists := s.index.Get(key)
	if !exists {
		return nil, false, nil
	}
	rec, err = s.store.Read(pos)
	if err != nil {
		return nil, false, err
	}
	return rec, true, nil
}

// Has returns true if the segment has the given key.
func (s *segment) Has(key string) bool {
	_, exists := s.index.Get(key)
	return exists
}

// IsMaxed returns true if the segment is at its max size.
func (s *segment) IsMaxed() bool {
	return s.store.size >= s.config.Segment.MaxStoreBytes
}

// Close  and closes the store and hint files.
func (s *segment) Close() error {
	if err := s.store.Close(); err != nil {
		return err
	}
	if err := s.hint.Close(); err != nil {
		return err
	}
	return nil
}

// Remove clears the index, closes the segment and removes the store and hint files.
func (s *segment) Remove() error {
	s.index.Clear()
	if err := s.Close(); err != nil {
		return err
	}
	if err := os.Remove(s.store.Name()); err != nil {
		return err
	}
	if err := os.Remove(s.hint.Name()); err != nil {
		return err
	}
	return nil
}
