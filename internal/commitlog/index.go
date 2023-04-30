package commitlog

import (
	cmap "github.com/orcaman/concurrent-map/v2"
)

// index is an in-memory index of the commit log.
type index struct {
	items cmap.ConcurrentMap[string, *item]
}

// item is an index item.
type item struct {
	offset    uint64
	timestamp int64
	size      int64
}

// newIndex creates a new index.
func newIndex() *index {
	return &index{
		items: cmap.New[*item](),
	}
}

// Get returns an item from the index.
func (i *index) Get(key string) (*item, bool) {
	return i.items.Get(key)
}

// Set stores an item in the index.
func (i *index) Set(key string, value *item) {
	i.items.Set(key, value)
}

// Delete removes an item from the index.
func (i *index) Delete(key string) {
	i.items.Remove(key)
}
