package cache

import (
	"errors"
	"time"

	"github.com/guackamolly/insta-archiver/internal/data/storage"
)

// A [CacheRepository] that uses the memory to persist values.
type MemoryCacheRepository[V any] struct {
	mem         *storage.MemoryStorage[string, CacheEntry[V]]
	cachePolicy time.Duration
}

func (r MemoryCacheRepository[V]) Load() (map[string]CacheEntry[V], error) {
	return map[string]CacheEntry[V]{}, nil
}

func (r MemoryCacheRepository[V]) Update(id string, value V) (CacheEntry[V], error) {
	ce := cache(value, r.cachePolicy)

	return r.mem.Store(id, ce)
}

func (r MemoryCacheRepository[V]) Lookup(id string) (CacheEntry[V], error) {
	ce, err := r.mem.Lookup(id)

	if err != nil {
		return ce, err
	}

	// if cache is outdated, evict value
	if ce.IsOutdated() {
		r.Evict(id)

		return ce, errors.New("cache entry is outdated")
	}

	return ce, nil
}

func (r MemoryCacheRepository[V]) Evict(id string) error {
	return r.mem.Delete(id)
}
