package domain

// An use case describes an action in the domain.
type Usecase[I any, O any] interface {
	Invoke(I) (O, error)
}
