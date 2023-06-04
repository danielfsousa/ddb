package bitcask

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

type hintArgs struct {
	key string
	pos uint64
}

var expectedHints = []hintArgs{
	{"hello world 1", 0},
	{"hello world 2", 25},
	{"hello world 3", 50},
}

func TestHintWriteScanClose(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "hint_write_scan_close_test")
	require.NoError(t, err)

	hint, err := newHint(f)
	require.NoError(t, err)

	testWrite(t, hint)
	testScan(t, hint)
	testClose(t, hint)
}

func testWrite(t *testing.T, hint *hint) {
	t.Helper()
	for _, arg := range expectedHints {
		err := hint.Write(arg.key, arg.pos)
		require.NoError(t, err)
	}
	err := hint.Sync()
	require.NoError(t, err)
}

func testScan(t *testing.T, hint *hint) {
	t.Helper()
	scanner, err := hint.Scanner()
	require.NoError(t, err)
	for i := 0; scanner.Scan(); i++ {
		err := scanner.Err()
		require.NoError(t, err)

		key, pos := scanner.Next()
		require.Equal(t, expectedHints[i].key, key)
		require.Equal(t, expectedHints[i].pos, pos)
	}
}

func testClose(t *testing.T, hint *hint) {
	t.Helper()
	err := hint.Close()
	require.NoError(t, err)
}
