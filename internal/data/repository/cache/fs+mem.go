package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/guackamolly/insta-archiver/internal/data/storage"
)

const defaultCacheFile = "cache.json"

// Default cache policy is 1h.
const defaultCachePolicy = 1 * time.Hour

// A [CacheRepository] that uses the file system to persist cache values
// between different program sessions, as well as the program memory.
type FileSystemMemoryCacheRepository[V any] struct {
	fs          *storage.FileSystemStorage
	mem         *storage.MemoryStorage[string, CacheEntry[V]]
	cacheFile   string
	cachePolicy time.Duration
}

func (r FileSystemMemoryCacheRepository[V]) Load() (map[string]CacheEntry[V], error) {
	// lookup root directories
	rfs, err := r.fs.Lookup("")

	if err != nil {
		return nil, err
	}

	m := map[string]CacheEntry[V]{}

	for _, rf := range rfs {
		if !rf.IsDir {
			continue
		}

		// lookup for cache file inside each directory
		k := rf.Name()
		cfs, err := r.fs.LookupFile(fmt.Sprintf("%s/%s", k, r.cacheFile))

		if err != nil {
			continue
		}

		v, err := r.tryLoadCache(k, cfs)

		if err == nil {
			m[k] = v
		}
	}

	return m, nil
}

func (r FileSystemMemoryCacheRepository[V]) Update(id string, value V) (CacheEntry[V], error) {
	ce := cache(value, r.cachePolicy)

	// first update mem storage
	mce, err := r.mem.Store(id, ce)

	if err == nil {
		ce = mce
	}

	// then update fs storage
	_, err = r.trySaveCache(fmt.Sprintf("%s/%s", id, r.cacheFile), ce)

	return ce, err
}

func (r FileSystemMemoryCacheRepository[V]) Lookup(id string) (CacheEntry[V], error) {
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

func (r FileSystemMemoryCacheRepository[V]) Evict(id string) error {
	// Evict from memory only, since Load and Lookup call this method if the cache is outdated
	// ensuring that clients never access old values.
	return r.mem.Delete(id)
}

func (r *FileSystemMemoryCacheRepository[V]) tryLoadCache(id string, file storage.File) (CacheEntry[V], error) {
	var v CacheEntry[V]

	bs, err := os.ReadFile(file.Path)

	if err != nil {
		return v, err
	}

	err = json.Unmarshal(bs, &v)

	if err != nil {
		return v, err
	}

	if v.IsOutdated() {
		return v, errors.New("cache entry is outdated")
	}

	vm, err := r.mem.Store(id, v)

	if err == nil {
		v = vm
	}

	return v, err
}

func (r *FileSystemMemoryCacheRepository[V]) trySaveCache(id string, value CacheEntry[V]) (storage.File, error) {
	var f storage.File

	j, err := json.Marshal(value)

	if err != nil {
		return f, err
	}

	return r.fs.StoreRaw(id, j)
}
