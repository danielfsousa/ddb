package bitcask

type Config struct {
	Segment struct {
		MaxStoreBytes uint64
		MaxIndexBytes uint64
	}
}
