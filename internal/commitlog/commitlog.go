package commitlog

import (
	"io"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"

	ddbv1 "github.com/danielfsousa/ddb/gen/ddb/v1"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type CommitLog struct {
	mu sync.RWMutex

	Dir    string
	Config Config

	activeSegment *segment
	segments      []*segment
}

// NewLog creates a log instance.
func NewCommitLog(dir string, c Config) (*CommitLog, error) {
	if c.Segment.MaxStoreBytes == 0 {
		c.Segment.MaxStoreBytes = 1024
	}
	log := &CommitLog{
		Dir:    dir,
		Config: c,
	}
	err := log.setup()
	return log, err
}

func (l *CommitLog) setup() error {
	files, err := os.ReadDir(l.Dir)
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
		if err = l.newSegment(ids[i]); err != nil {
			return err
		}
		i += 1 // skip hint files
	}
	if l.segments == nil {
		if err = l.newSegment(1); err != nil {
			return err
		}
	}

	return nil
}

func (l *CommitLog) newSegment(id uint64) error {
	s, err := newSegment(l.Dir, id, l.Config)
	if err != nil {
		return err
	}
	l.segments = append(l.segments, s)
	l.activeSegment = s
	return nil
}

// Keys returns a slice of the keys of all records stored in the log.
func (l *CommitLog) Keys() []string {
	l.mu.RLock()
	defer l.mu.RUnlock()
	keys := make(map[string]bool)
	for _, segment := range l.segments {
		for _, key := range segment.Keys() {
			keys[key] = true
		}
	}
	keysSlice := maps.Keys(keys)
	slices.Sort(keysSlice)
	return keysSlice
}

// Append appends a record to the log.
func (l *CommitLog) Append(rec *ddbv1.Record) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if err := l.activeSegment.Append(rec); err != nil {
		return err
	}
	if l.activeSegment.IsMaxed() {
		if err := l.newSegment(l.activeSegment.id + 1); err != nil {
			return err
		}
	}
	return nil
}

// Get returns a record by key.
func (l *CommitLog) Get(key string) (rec *ddbv1.Record, exists bool, err error) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	for i := len(l.segments) - 1; i >= 0; i-- {
		rec, exists, err = l.segments[i].Get(key)
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

// Has returns true if the key exists in the log.
func (l *CommitLog) Has(key string) bool {
	l.mu.RLock()
	defer l.mu.RUnlock()
	for i := len(l.segments) - 1; i >= 0; i-- {
		if l.segments[i].Has(key) {
			return true
		}
	}
	return false
}

// Close closes the log.
func (l *CommitLog) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	for _, segment := range l.segments {
		if err := segment.Close(); err != nil {
			return err
		}
	}
	return nil
}

// Remove closes the log and then removes its data from the filesystem.
func (l *CommitLog) Remove() error {
	if err := l.Close(); err != nil {
		return err
	}
	return os.RemoveAll(l.Dir)
}

// Reset removes the log and re-creates it.
func (l *CommitLog) Reset() error {
	if err := l.Remove(); err != nil {
		return err
	}
	return l.setup()
}

// Reader returns an io.Reader instance to read the whole log.
func (l *CommitLog) Reader() io.Reader {
	l.mu.RLock()
	defer l.mu.RUnlock()
	readers := make([]io.Reader, len(l.segments))
	for i, segment := range l.segments {
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
