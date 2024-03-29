package bitcask

import (
	"testing"

	"github.com/danielfsousa/ddb/internal/backend"
	"github.com/stretchr/testify/require"
)

var expectedIndexKeys = []string{
	"hello world 1",
	"hello world 2",
	"hello world 3",
}
var expectedIndexPos = []backend.RecordMetadata{
	{Pos: 0, DeletedAt: nil},
	{Pos: 25, DeletedAt: nil},
	{Pos: 50, DeletedAt: nil},
}

func TestIndexSetGetDelete(t *testing.T) {
	idx := newIndex()
	testSetGet(t, idx)
	testDeleteGet(t, idx)
}

func testSetGet(t *testing.T, idx *index) {
	t.Helper()
	for i, key := range expectedIndexKeys {
		idx.Set(key, expectedIndexPos[i])
		item, exists := idx.Get(key)
		require.Equal(t, expectedIndexPos[i], item)
		require.True(t, exists)
	}
}

func testDeleteGet(t *testing.T, idx *index) {
	t.Helper()

	item, exists := idx.Get(expectedIndexKeys[0])
	require.Equal(t, expectedIndexPos[0], item)
	require.True(t, exists)

	idx.Delete(expectedIndexKeys[0])
	item, exists = idx.Get(expectedIndexKeys[0])
	require.Zero(t, item)
	require.False(t, exists)
}
