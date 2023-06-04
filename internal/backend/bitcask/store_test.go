package bitcask

import (
	"os"
	"testing"
	"time"

	ddbv1 "github.com/danielfsousa/ddb/gen/ddb/v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

var expectedWrites = []*ddbv1.Record{
	{
		Timestamp: time.Now().UnixNano(),
		Key:       "1",
		Value:     []byte("hello world 1"),
	},
	{
		Timestamp: time.Now().UnixNano(),
		Key:       "2",
		Value:     []byte("hello world 2"),
	},
	{
		Timestamp: time.Now().UnixNano(),
		Key:       "3",
		Value:     []byte("hello world 3"),
	},
}

func TestStoreAppendReadScan(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "store_append_read_scan_test")
	require.NoError(t, err)

	store, err := newStore(f)
	require.NoError(t, err)

	testAppend(t, store)
	testRead(t, store)
	testReadAt(t, store)
	testScanner(t, store)

	store, err = newStore(f)
	require.NoError(t, err)
	testRead(t, store)
}

func testAppend(t *testing.T, store *store) {
	t.Helper()
	for i := uint64(1); i < 4; i++ {
		write := expectedWrites[i-1]
		width := uint64(proto.Size(write)) + storeHeaderSize
		n, pos, err := store.Append(write)
		require.NoError(t, err)
		require.Equal(t, pos+n, width*i)
	}
}

func testRead(t *testing.T, store *store) {
	t.Helper()
	var pos uint64
	for i := uint64(1); i < 4; i++ {
		read, err := store.Read(pos)
		require.NoError(t, err)
		require.True(t, proto.Equal(expectedWrites[i-1], read))
		pos += uint64(proto.Size(read)) + storeHeaderSize
	}
}

func testScanner(t *testing.T, store *store) {
	t.Helper()
	scanner, err := store.Scanner()
	require.NoError(t, err)
	for i := 0; scanner.Scan(); i++ {
		rec, _ := scanner.Next()
		require.True(t, proto.Equal(expectedWrites[i], rec))
	}
	require.NoError(t, scanner.Err())
}

func testReadAt(t *testing.T, store *store) {
	t.Helper()
	for i, off := uint64(1), int64(0); i < 4; i++ {
		b := make([]byte, storeHeaderSize)
		n, err := store.ReadAt(b, off)
		require.NoError(t, err)
		require.Equal(t, storeHeaderSize, n)

		off += int64(storeHeaderSize)
		write := expectedWrites[i-1]
		recordLen := encoding.Uint64(b[checksumSize:])

		b = make([]byte, recordLen)
		n, err = store.ReadAt(b, off)
		require.NoError(t, err)

		rec := &ddbv1.Record{}
		err = proto.Unmarshal(b, rec)
		require.NoError(t, err)

		require.True(t, proto.Equal(write, rec))
		require.Equal(t, int(recordLen), n)

		off += int64(n)
	}
}

func TestStoreSync(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "store_sync_test")
	require.NoError(t, err)

	store, err := newStore(f)
	require.NoError(t, err)

	n, _, err := store.Append(expectedWrites[0])
	require.NoError(t, err)

	_, beforeSize, err := openFile(f.Name())
	require.NoError(t, err)
	require.Equal(t, beforeSize, int64(0))

	err = store.Sync()
	require.NoError(t, err)

	_, afterSize, err := openFile(f.Name())
	require.NoError(t, err)
	require.Equal(t, afterSize, int64(n))
}

func TestStoreClose(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "store_close_test")
	require.NoError(t, err)

	store, err := newStore(f)
	require.NoError(t, err)

	_, _, err = store.Append(expectedWrites[0])
	require.NoError(t, err)

	f, beforeSize, err := openFile(f.Name())
	require.NoError(t, err)

	err = store.Close()
	require.NoError(t, err)

	_, afterSize, err := openFile(f.Name())
	require.NoError(t, err)
	require.True(t, afterSize > beforeSize)
}

func openFile(name string) (file *os.File, size int64, err error) {
	f, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, 0, err
	}

	fi, err := f.Stat()
	if err != nil {
		return nil, 0, err
	}

	return f, fi.Size(), nil
}
