package commitlog

import (
	"bufio"
	"encoding/binary"
	"hash"
	"hash/crc32"
	"os"
	"sync"
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
	checksumWidth = 4
	lenWidth      = 8
	metaWidth     = checksumWidth + lenWidth
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
	crc  hash.Hash32
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
		crc:  crc32.New(crcTable),
	}, nil
}

// Append persists the given bytes to the store.
func (s *store) Append(in []byte) (n, pos uint64, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	pos = s.size
	recordLen := uint64(len(in))

	s.crc.Reset()
	if _, err = s.crc.Write(in); err != nil {
		return 0, 0, err
	}
	checksum := s.crc.Sum32()

	// write checksum of the record
	if err = binary.Write(s.buf, encoding, checksum); err != nil {
		return 0, 0, err
	}

	// write length of the record
	if err = binary.Write(s.buf, encoding, recordLen); err != nil {
		return 0, 0, err
	}

	w, err := s.buf.Write(in)
	if err != nil {
		return 0, 0, err
	}

	w += metaWidth
	s.size += uint64(w)

	return uint64(w), pos, nil
}

// Read returns the record stored at the given position.
func (s *store) Read(pos uint64) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.buf.Flush(); err != nil {
		return nil, err
	}

	recordSize := make([]byte, lenWidth)
	if _, err := s.File.ReadAt(recordSize, int64(pos+checksumWidth)); err != nil {
		return nil, err
	}

	bytes := make([]byte, encoding.Uint64(recordSize))
	if _, err := s.File.ReadAt(bytes, int64(pos+metaWidth)); err != nil {
		return nil, err
	}
	return bytes, nil
}

// ReadAt reads len(p) bytes from the store starting at byte offset off.
// It implements the io.ReaderAt interface.
func (s *store) ReadAt(in []byte, off int64) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.buf.Flush(); err != nil {
		return 0, err
	}

	return s.File.ReadAt(in, off)
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
