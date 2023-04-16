package commitlog

import (
	"encoding/binary"
	"errors"
	"fmt"
	"hash"
	"hash/crc32"
	"io"
	"os"

	"google.golang.org/protobuf/proto"

	ddbv1 "github.com/danielfsousa/ddb/gen/ddb/v1"
)

// Scanner is a store scanner.
type Scanner struct {
	*os.File
	crc    hash.Hash32
	record *ddbv1.Record
	err    error
}

// NewStoreScanner returns a new instance of Scanner.
func newStoreScanner(s *store) *Scanner {
	return &Scanner{
		File:   s.File,
		crc:    crc32.New(crcTable),
		record: &ddbv1.Record{},
	}
}

// Scan advances the scanner to the next record.
func (s *Scanner) Scan() bool {
	s.record.Reset()

	var head [metaWidth]byte

	if _, err := io.ReadFull(s.File, head[:]); err != nil {
		if errors.Is(s.err, io.EOF) {
			s.err = nil
		}
		return false
	}

	checksum := binary.BigEndian.Uint32(head[:checksumWidth])
	recordLen := binary.BigEndian.Uint64(head[checksumWidth:])

	data := make([]byte, recordLen)
	if _, s.err = io.ReadFull(s.File, data); s.err != nil {
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

	return true
}

// Returns the current record.
// Only valid until the next Scan() call.
func (s *Scanner) Record() *ddbv1.Record {
	return s.record
}

// Returns last error encountered by the scanner.
func (s *Scanner) Err() error {
	return s.err
}
