<!-- Store -->

goos: darwin
goarch: arm64
pkg: github.com/danielfsousa/ddb/internal/commitlog
BenchmarkAppend100NoSync-8      	 4230967	       265.1 ns/op	 433.73 MB/s	     144 B/op	       2 allocs/op
BenchmarkAppend100NoBatch-8     	     319	   3695215 ns/op	   0.03 MB/s	     144 B/op	       2 allocs/op
BenchmarkAppend100Batch10-8     	    3147	    382251 ns/op	   0.30 MB/s	     144 B/op	       2 allocs/op
BenchmarkAppend100Batch100-8    	   30296	     40691 ns/op	   2.83 MB/s	     144 B/op	       2 allocs/op
BenchmarkAppend1000NoSync-8     	  796854	      1433 ns/op	 708.94 MB/s	    1040 B/op	       2 allocs/op
BenchmarkAppend1000NoBatch-8    	     314	   3779296 ns/op	   0.27 MB/s	    1040 B/op	       2 allocs/op
BenchmarkAppend1000Batch10-8    	    2658	    411826 ns/op	   2.47 MB/s	    1040 B/op	       2 allocs/op
BenchmarkAppend1000Batch100-8   	   25466	     48132 ns/op	  21.11 MB/s	    1040 B/op	       2 allocs/op
PASS
ok  	github.com/danielfsousa/ddb/internal/commitlog	12.256s

==========

goos: darwin
goarch: arm64
pkg: github.com/danielfsousa/ddb/internal/commitlog
BenchmarkAppend100NoSync-8      	 4187924	       270.4 ns/op	 425.23 MB/s	     144 B/op	       2 allocs/op
BenchmarkAppend100NoBatch-8     	     321	   3643622 ns/op	   0.03 MB/s	     144 B/op	       2 allocs/op
BenchmarkAppend100Batch10-8     	    3252	    376027 ns/op	   0.31 MB/s	     144 B/op	       2 allocs/op
BenchmarkAppend100Batch100-8    	   29262	     42020 ns/op	   2.74 MB/s	     144 B/op	       2 allocs/op
BenchmarkAppend1000NoSync-8     	  807238	      1392 ns/op	 729.80 MB/s	    1040 B/op	       2 allocs/op
BenchmarkAppend1000NoBatch-8    	     319	   3759983 ns/op	   0.27 MB/s	    1040 B/op	       2 allocs/op
BenchmarkAppend1000Batch10-8    	    2953	    403528 ns/op	   2.52 MB/s	    1040 B/op	       2 allocs/op
BenchmarkAppend1000Batch100-8   	   25308	     46974 ns/op	  21.63 MB/s	    1040 B/op	       2 allocs/op
PASS
ok  	github.com/danielfsousa/ddb/internal/commitlog	12.337s

<!-- Bitcask -->

<!-- fsync disabled -->
goos: darwin
goarch: arm64
pkg: github.com/danielfsousa/ddb/internal/backend/bitcask
BenchmarkBitcaskSet100NoSync-8       	 1526660	       782.9 ns/op	 149.45 MB/s	     270 B/op	       4 allocs/op
BenchmarkBitcaskSet100NoBatch-8      	     318	   3658448 ns/op	   0.03 MB/s	     227 B/op	       3 allocs/op
BenchmarkBitcaskSet100Batch10-8      	    3138	    378517 ns/op	   0.30 MB/s	     264 B/op	       4 allocs/op
BenchmarkBitcaskSet100Batch100-8     	   27882	     40261 ns/op	   2.86 MB/s	     351 B/op	       4 allocs/op
BenchmarkBitcaskSet1000NoSync-8      	  539024	      2053 ns/op	 495.34 MB/s	    1207 B/op	       4 allocs/op
BenchmarkBitcaskSet1000NoBatch-8     	     315	   3872722 ns/op	   0.26 MB/s	    1125 B/op	       3 allocs/op
BenchmarkBitcaskSet1000Batch10-8     	    3060	    410733 ns/op	   2.47 MB/s	    1163 B/op	       4 allocs/op
BenchmarkBitcaskSet1000Batch100-8    	   22622	     55116 ns/op	  18.43 MB/s	    1174 B/op	       4 allocs/op
BenchmarkStoreAppend100NoSync-8      	 4028493	       261.0 ns/op	 433.03 MB/s	     144 B/op	       2 allocs/op
BenchmarkStoreAppend100NoBatch-8     	     328	   3687137 ns/op	   0.03 MB/s	     144 B/op	       2 allocs/op
BenchmarkStoreAppend100Batch10-8     	    3074	    385116 ns/op	   0.29 MB/s	     144 B/op	       2 allocs/op
BenchmarkStoreAppend100Batch100-8    	   26202	     44280 ns/op	   2.55 MB/s	     144 B/op	       2 allocs/op
BenchmarkStoreAppend1000NoSync-8     	  780655	      1414 ns/op	 717.13 MB/s	    1040 B/op	       2 allocs/op
BenchmarkStoreAppend1000NoBatch-8    	     314	   3794648 ns/op	   0.27 MB/s	    1040 B/op	       2 allocs/op
BenchmarkStoreAppend1000Batch10-8    	    3052	    425936 ns/op	   2.38 MB/s	    1040 B/op	       2 allocs/op
BenchmarkStoreAppend1000Batch100-8   	   24505	     48740 ns/op	  20.80 MB/s	    1040 B/op	       2 allocs/op
PASS
ok  	github.com/danielfsousa/ddb/internal/backend/bitcask	25.353s
