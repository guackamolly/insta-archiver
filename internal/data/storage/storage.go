package storage

// Abstracts interaction with a data storage.
// Typed to fully abstract the input and output types.
type Storage[I any, O any] interface {
	Lookup(I) (O, error)
	Store(I, O) (O, error)
	Delete(I) error
}
