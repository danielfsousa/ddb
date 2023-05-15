package commitlog

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"hash"
	"hash/crc32"
	"io"
	"os"
	"sync"

	ddbv1 "github.com/danielfsousa/ddb/gen/ddb/v1"
	"google.golang.org/protobuf/proto"
)

// =========== Store Format ===========
// +----------+--------------+--------+
// |         metadata        |  data  |
// +-------------------------+--------+
// | checksum | recordLength | record |
// +----------+--------------+--------+
// | 4 bytes  | 8 bytes      | ?      |
// +----------+--------------+--------+

const (
	checksumSize    = 4
	recLenSize      = 8
	storeHeaderSize = checksumSize + recLenSize
)

var (
	// encoding is the byte order to use for internal disk serialization.
	encoding = binary.BigEndian
	crcTable = crc32.MakeTable(crc32.Castagnoli)
)

type store struct {
	*os.File
	mu   sync.Mutex
	buf  *bufio.Writer
	hash hash.Hash32
	size uint64
}

func newStore(f *os.File) (*store, error) {
	fi, err := os.Stat(f.Name())
	if err != nil {
		return nil, err
	}
	size := uint64(fi.Size())
	return &store{
		File: f,
		size: size,
		buf:  bufio.NewWriter(f),
		hash: crc32.New(crcTable),
	}, nil
}

// Name returns the store's file path.
func (s *store) Name() string {
	return s.File.Name()
}

// Append persists the given bytes to the store.
func (s *store) Append(rec *ddbv1.Record) (n, pos uint64, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	b, err := proto.Marshal(rec)
	if err != nil {
		return 0, 0, fmt.Errorf("store failed to marshal record: %w", err)
	}

	pos = s.size
	recordLen := uint64(len(b))

	// calculate checksum
	s.hash.Reset()
	if _, err = s.hash.Write(b); err != nil {
		return 0, 0, err
	}
	checksum := s.hash.Sum32()

	// serialize header
	metadata := [storeHeaderSize]byte{}
	binary.BigEndian.PutUint32(metadata[:checksumSize], checksum)
	binary.BigEndian.PutUint64(metadata[checksumSize:], recordLen)

	// write header
	bytesMetadata, err := s.buf.Write(metadata[:])
	if err != nil {
		return 0, 0, err
	}

	// write data
	bytesRecord, err := s.buf.Write(b)
	if err != nil {
		return 0, 0, err
	}

	writtenBytes := uint64(bytesMetadata + bytesRecord)
	s.size += writtenBytes

	return writtenBytes, pos, nil
}

// Read returns the record stored at the given position.
func (s *store) Read(pos uint64) (*ddbv1.Record, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.buf.Flush(); err != nil {
		return nil, err
	}

	recordLen := make([]byte, recLenSize)
	if _, err := s.File.ReadAt(recordLen, int64(pos+checksumSize)); err != nil {
		return nil, err
	}

	b := make([]byte, encoding.Uint64(recordLen))
	if _, err := s.File.ReadAt(b, int64(pos+storeHeaderSize)); err != nil {
		return nil, err
	}

	rec := &ddbv1.Record{}
	if err := proto.Unmarshal(b, rec); err != nil {
		return nil, fmt.Errorf("store failed to marshal record: %w", err)
	}

	return rec, nil
}

// ReadAt reads len(in) bytes from the store starting at byte offset off.
// It implements the io.ReaderAt interface.
func (s *store) ReadAt(in []byte, offset int64) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.buf.Flush(); err != nil {
		return 0, err
	}

	return s.File.ReadAt(in, offset)
}

// Sync flushes the store to disk.
func (s *store) Sync() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.buf.Flush(); err != nil {
		return err
	}
	return s.File.Sync()
}

// Close persists any buffered data before closing the file.
func (s *store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.buf.Flush(); err != nil {
		return err
	}
	return s.File.Close()
}

// Scanner returns a new storeScanner for iterating over the records in the store.
func (s *store) Scanner() (*storeScanner, error) {
	f, err := os.Open(s.File.Name())
	if err != nil {
		return nil, err
	}
	return &storeScanner{
		file:   f,
		crc:    crc32.New(crcTable),
		record: &ddbv1.Record{},
	}, nil
}

// storeScanner enables iterating over the records in the store.
type storeScanner struct {
	file    *os.File
	crc     hash.Hash32
	record  *ddbv1.Record
	pos     uint64
	nextPos uint64
	err     error
}

// Scan advances the scanner to the next record.
func (s *storeScanner) Scan() bool {
	s.record.Reset()

	var header [storeHeaderSize]byte
	if _, s.err = io.ReadFull(s.file, header[:]); s.err != nil {
		if errors.Is(s.err, io.EOF) {
			s.err = nil
		}
		return false
	}

	checksum := binary.BigEndian.Uint32(header[:checksumSize])
	recordLen := binary.BigEndian.Uint64(header[checksumSize:])

	data := make([]byte, recordLen)
	if _, s.err = io.ReadFull(s.file, data); s.err != nil {
		return false
	}

	s.crc.Reset()
	if _, s.err = s.crc.Write(data); s.err != nil {
		return false
	}
	c := s.crc.Sum32()
	if c != checksum {
		s.err = fmt.Errorf("checksum mismatch. Expected %d, got %d", checksum, c)
		return false
	}

	if s.err = proto.Unmarshal(data, s.record); s.err != nil {
		return false
	}

	s.pos = s.nextPos
	s.nextPos += storeHeaderSize + recordLen

	return true
}

// Returns the current record.
func (s *storeScanner) Next() (rec *ddbv1.Record, pos uint64) {
	return s.record, s.pos
}

// Returns last error encountered by the scanner.
func (s *storeScanner) Err() error {
	return s.err
}
