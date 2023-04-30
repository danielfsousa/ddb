package commitlog

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var expectedIndexKeys = []string{
	"hello world 1",
	"hello world 2",
	"hello world 3",
}

var expectedIndexItems = []*item{
	{offset: 0, timestamp: 0, size: 25},
	{offset: 25, timestamp: 0, size: 25},
	{offset: 50, timestamp: 0, size: 25},
}

func TestIndexSetGetDelete(t *testing.T) {
	idx := newIndex()
	testSetGet(t, idx)
	testDeleteGet(t, idx)
}

func testSetGet(t *testing.T, idx *index) {
	t.Helper()
	for i, key := range expectedIndexKeys {
		idx.Set(key, expectedIndexItems[i])
		item, exists := idx.Get(key)
		require.Equal(t, expectedIndexItems[i], item)
		require.True(t, exists)
	}
}

func testDeleteGet(t *testing.T, idx *index) {
	t.Helper()

	item, exists := idx.Get(expectedIndexKeys[0])
	require.Equal(t, expectedIndexItems[0], item)
	require.True(t, exists)

	idx.Delete(expectedIndexKeys[0])
	item, exists = idx.Get(expectedIndexKeys[0])
	require.Nil(t, item)
	require.False(t, exists)
}
