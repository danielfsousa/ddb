//nolint:errcheck
package bitcask

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	ddbv1 "github.com/danielfsousa/ddb/gen/ddb/v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func BenchmarkStoreAppend100NoSync(b *testing.B)   { benchmarkStore(b, 100, 0) }
func BenchmarkStoreAppend100NoBatch(b *testing.B)  { benchmarkStore(b, 100, 1) }
func BenchmarkStoreAppend100Batch10(b *testing.B)  { benchmarkStore(b, 100, 10) }
func BenchmarkStoreAppend100Batch100(b *testing.B) { benchmarkStore(b, 100, 100) }

func BenchmarkStoreAppend1000NoSync(b *testing.B)   { benchmarkStore(b, 1000, 0) }
func BenchmarkStoreAppend1000NoBatch(b *testing.B)  { benchmarkStore(b, 1000, 1) }
func BenchmarkStoreAppend1000Batch10(b *testing.B)  { benchmarkStore(b, 1000, 10) }
func BenchmarkStoreAppend1000Batch100(b *testing.B) { benchmarkStore(b, 1000, 100) }

func benchmarkStore(b *testing.B, dataSize, batchSize int) {
	b.Helper()
	b.ReportAllocs()
	tempdir := b.TempDir()
	fname := filepath.Join(tempdir, "commitlog_bench")
	f, err := os.OpenFile(fname, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	require.NoError(b, err)
	defer f.Close()

	record := &ddbv1.Record{
		Timestamp: time.Now().Unix(),
		Key:       "key",
		Value:     []byte(strings.Repeat("a", dataSize)),
	}
	b.SetBytes(int64(proto.Size(record)))

	store, err := newStore(f)
	require.NoError(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		store.Append(record)
		if batchSize > 0 && i%batchSize == 0 {
			store.Sync()
		}
	}
	b.StopTimer()
}
