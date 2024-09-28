package cache

import (
	"time"

	"github.com/guackamolly/insta-archiver/internal/data/storage"
)

// Represents an entry in a cache store.
type CacheEntry[V any] struct {
	NextHit       time.Time     `json:"nextHit"`
	RefreshPolicy time.Duration `json:"refreshPolicy"`
	Value         V             `json:"value"`
}

// Caches a value, returning an entry to add in the cache store.
func cache[V any](value V, refreshPolicy time.Duration) CacheEntry[V] {
	t := time.Now()

	return CacheEntry[V]{
		NextHit:       t.Add(refreshPolicy),
		RefreshPolicy: refreshPolicy,
		Value:         value,
	}
}

// Checks if a cache value is outdated.
func (ce CacheEntry[V]) IsOutdated() bool {
	return time.Now().After(ce.NextHit)
}

// Interface similar to [storage.Storage], but focused on
// caching.
// Provides methods to load and update the cache store, and
// method to lookup a cache value.
//
// Implementations of this repository should be capable of loading a previous cache version
// from a data source.
type CacheRepository[I comparable, V any] interface {
	// (Re)Loads the cache store.
	Load() (map[I]CacheEntry[V], error)
	// Returns the cache entry associated with the identifier. Outdated cache entries are evicted.
	Lookup(I) (CacheEntry[V], error)
	// Caches a value by associating a new cache entry with the identifier.
	Update(I, V) (CacheEntry[V], error)
	// Evicts/deletes a cache entry.
	Evict(I) error
}

func NewMemoryCacheRepository[V any](
	memStorage *storage.MemoryStorage[string, CacheEntry[V]],
) MemoryCacheRepository[V] {
	return MemoryCacheRepository[V]{
		mem:         memStorage,
		cachePolicy: time.Hour * 24,
	}
}

func NewFileSystemMemoryCacheRepository[V any](
	fsStorage *storage.FileSystemStorage,
	memStorage *storage.MemoryStorage[string, CacheEntry[V]],
) FileSystemMemoryCacheRepository[V] {
	return FileSystemMemoryCacheRepository[V]{
		fs:          fsStorage,
		mem:         memStorage,
		cacheFile:   defaultCacheFile,
		cachePolicy: defaultCachePolicy,
	}
}
