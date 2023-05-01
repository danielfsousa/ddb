package commitlog

import (
	"io"
	"testing"

	ddbv1 "github.com/danielfsousa/ddb/gen/ddb/v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func TestSegment(t *testing.T) {
	dir := t.TempDir()

	want := &ddbv1.Record{
		Timestamp: 0,
		Key:       "54",
		Value:     []byte("Hello world"),
	}
	wantSize := proto.Size(want)

	config := Config{}
	config.Segment.MaxStoreBytes = uint64((wantSize + storeHeaderSize) * 3)

	segment, err := newSegment(dir, 16, config)
	require.NoError(t, err)
	require.False(t, segment.IsMaxed())

	for i := uint64(0); i < 3; i++ {
		want.Timestamp = int64(i)
		err = segment.Append(want)
		require.NoError(t, err)

		got, exists, err := segment.Get(want.Key) //nolint
		require.NoError(t, err)
		require.True(t, exists)
		require.True(t, proto.Equal(want, got))
	}

	err = segment.Append(want)
	require.Equal(t, io.EOF, err)
	require.True(t, segment.IsMaxed())

	segment.Close()
	segment, err = newSegment(dir, 16, config)
	require.NoError(t, err)
	require.True(t, segment.IsMaxed())

	err = segment.Remove()
	require.NoError(t, err)

	segment.Close()
	segment, err = newSegment(dir, 16, config)
	require.NoError(t, err)
	require.False(t, segment.IsMaxed())
}
