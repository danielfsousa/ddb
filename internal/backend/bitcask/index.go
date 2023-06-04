package bitcask

import (
	"github.com/danielfsousa/ddb/internal/backend"
	cmap "github.com/orcaman/concurrent-map/v2"
)

// index is an in-memory index of the commit log.
type index struct {
	items cmap.ConcurrentMap[string, backend.RecordMetadata]
}

// newIndex creates a new index.
func newIndex() *index {
	return &index{
		items: cmap.New[backend.RecordMetadata](),
	}
}

// Keys returns all keys in the index.
func (i *index) Keys() []string {
	return i.items.Keys()
}

// Get returns an item from the index.
func (i *index) Get(key string) (entry backend.RecordMetadata, exists bool) {
	return i.items.Get(key)
}

// Set stores an item in the index.
func (i *index) Set(key string, entry backend.RecordMetadata) {
	i.items.Set(key, entry)
}

// Delete removes an item from the index.
func (i *index) Delete(key string) {
	i.items.Remove(key)
}

// Clear removes all items from the index.
func (i *index) Clear() {
	i.items.Clear()
}
