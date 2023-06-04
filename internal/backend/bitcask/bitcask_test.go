package bitcask

import (
	"testing"

	ddbv1 "github.com/danielfsousa/ddb/gen/ddb/v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func TestCommitLog(t *testing.T) {
	tests := map[string]func(t *testing.T, log *Bitcask){
		"append and read a record succeeds": testAppendGet,
		"init with existing segments":       testInitExisting,
	}
	for scenario, fn := range tests {
		t.Run(scenario, func(t *testing.T) {
			dir := t.TempDir()

			config := Config{}
			config.Segment.MaxStoreBytes = 1024
			log, err := NewBitcaskBackend(dir, config)
			require.NoError(t, err)

			fn(t, log)
		})
	}
}

func testAppendGet(t *testing.T, log *Bitcask) {
	want := &ddbv1.Record{
		Timestamp: 12345,
		Key:       "foo",
		Value:     []byte("hello world"),
	}

	err := log.Set(want)
	require.NoError(t, err)

	got, exists, err := log.Get(want.Key)
	require.NoError(t, err)
	require.True(t, exists)
	require.True(t, proto.Equal(want, got))

	got, exists, err = log.Get("does_not_exist")
	require.NoError(t, err)
	require.False(t, exists)
	require.Nil(t, got)
}

func testInitExisting(t *testing.T, log *Bitcask) {
	recs := []*ddbv1.Record{
		{
			Timestamp: 12345,
			Key:       "foo",
			Value:     []byte("hello world"),
		},
		{
			Timestamp: 54321,
			Key:       "bar",
			Value:     []byte("hi world"),
		},
		{
			Timestamp: 99999,
			Key:       "Keeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeey",
			Value:     []byte("Vaaaaaaaaaaaaaaaaaaaaaaaaaaaaaalue"),
		},
	}

	for _, want := range recs {
		err := log.Set(want)
		require.NoError(t, err)
	}

	log.Close()
	log, err := NewBitcaskBackend(log.Dir, log.Config)
	require.NoError(t, err)

	for _, want := range recs {
		got, exists, err := log.Get(want.Key)
		require.NoError(t, err)
		require.True(t, exists)
		require.True(t, proto.Equal(want, got))
	}
}
