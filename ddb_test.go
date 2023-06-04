package ddb

import (
	"testing"

	ddbv1 "github.com/danielfsousa/ddb/gen/ddb/v1"
	"github.com/stretchr/testify/require"
)

func TestCommitLog(t *testing.T) {
	tests := map[string]func(t *testing.T, log *Ddb){
		"write, read and delete a record succeeds": testReadWriteDelete,
		"init with existing segments":              testInitExisting,
	}
	for scenario, fn := range tests {
		t.Run(scenario, func(t *testing.T) {
			dir := t.TempDir()
			ddb, err := newDdb(dir)
			require.NoError(t, err)
			fn(t, ddb)
		})
	}
}

func newDdb(dir string) (*Ddb, error) {
	return Open(dir, WithMaxSegmentDataSize(1024))
}

func testReadWriteDelete(t *testing.T, ddb *Ddb) {
	want := &ddbv1.Record{
		Key:   "foo",
		Value: []byte("hello world"),
	}

	require.False(t, ddb.Has(want.Key))

	err := ddb.Set(want.Key, want.Value)
	require.NoError(t, err)

	require.True(t, ddb.Has(want.Key))

	got, err := ddb.Get(want.Key)
	require.NoError(t, err)
	require.Equal(t, want.Value, got)

	err = ddb.Delete(want.Key)
	require.NoError(t, err)

	require.False(t, ddb.Has(want.Key))

	got, err = ddb.Get(want.Key)
	require.Error(t, err, ErrKeyNotFound)
	require.Nil(t, got)
}

func testInitExisting(t *testing.T, ddb *Ddb) {
	recs := []*ddbv1.Record{
		{
			Key:   "foo",
			Value: []byte("hello world"),
		},
		{
			Key:   "bar",
			Value: []byte("hi world"),
		},
		{
			Key:   "Keeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeey",
			Value: []byte("Vaaaaaaaaaaaaaaaaaaaaaaaaaaaaaalue"),
		},
	}

	for _, want := range recs {
		err := ddb.Set(want.Key, want.Value)
		require.NoError(t, err)
	}

	ddb.Close()
	ddb, err := newDdb(ddb.dir)
	require.NoError(t, err)

	for _, want := range recs {
		got, err := ddb.Get(want.Key)
		require.NoError(t, err)
		require.Equal(t, want.Value, got)
	}
}
