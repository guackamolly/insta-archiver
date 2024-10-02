package storage

import (
	"errors"
	"sync"
)

// A [Storage] implementation that uses the process memory as the provider.
// The input key should implement the [comparable] interface.
type MemoryStorage[I comparable, O any] struct {
	cache map[I]O
	lock  *sync.RWMutex
}

func NewMemoryStorage[I comparable, O any]() *MemoryStorage[I, O] {
	return &MemoryStorage[I, O]{
		cache: map[I]O{},
		lock:  &sync.RWMutex{},
	}
}

func (s *MemoryStorage[I, O]) Lookup(input I) (O, error) {
	v, ok := s.cache[input]

	if !ok {
		return v, errors.New("value is not cached")
	}

	return v, nil
}

func (s *MemoryStorage[I, O]) Store(input I, value O) (O, error) {
	s.lock.Lock()
	s.cache[input] = value
	s.lock.Unlock()

	return value, nil
}

func (s *MemoryStorage[I, O]) Delete(input I) error {
	s.lock.Lock()
	delete(s.cache, input)
	s.lock.Unlock()

	return nil
}
