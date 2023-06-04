//nolint:errcheck
package bitcask

import (
	"fmt"
	"strings"
	"testing"
	"time"

	ddbv1 "github.com/danielfsousa/ddb/gen/ddb/v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func BenchmarkBitcaskSet100NoSync(b *testing.B)   { benchmark(b, 100, 0) }
func BenchmarkBitcaskSet100NoBatch(b *testing.B)  { benchmark(b, 100, 1) }
func BenchmarkBitcaskSet100Batch10(b *testing.B)  { benchmark(b, 100, 10) }
func BenchmarkBitcaskSet100Batch100(b *testing.B) { benchmark(b, 100, 100) }

func BenchmarkBitcaskSet1000NoSync(b *testing.B)   { benchmark(b, 1000, 0) }
func BenchmarkBitcaskSet1000NoBatch(b *testing.B)  { benchmark(b, 1000, 1) }
func BenchmarkBitcaskSet1000Batch10(b *testing.B)  { benchmark(b, 1000, 10) }
func BenchmarkBitcaskSet1000Batch100(b *testing.B) { benchmark(b, 1000, 100) }

func benchmark(b *testing.B, dataSize, batchSize int) {
	b.Helper()
	b.ReportAllocs()
	tempdir := b.TempDir()

	config := Config{}
	config.Segment.MaxIndexBytes = 1e+8 // 100MB
	config.Segment.MaxStoreBytes = 1e+9 // 1GB

	db, err := NewBitcaskBackend(tempdir, config)
	require.NoError(b, err)
	defer db.Close()

	record := &ddbv1.Record{
		Timestamp: time.Now().Unix(),
		Key:       "key",
		Value:     []byte(strings.Repeat("a", dataSize)),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		record.Key = fmt.Sprint(i)
		b.SetBytes(int64(proto.Size(record)))
		err := db.Set(record)
		require.NoError(b, err)
		if batchSize > 0 && i%batchSize == 0 {
			db.Sync()
		}
	}
	b.StopTimer()
}
