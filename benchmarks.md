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
