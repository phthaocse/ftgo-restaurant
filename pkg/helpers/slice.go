package helpers

type ReadOnlySlice[T any] struct {
	slice []T
}

func NewReadOnlySlice[T any](slice []T) *ReadOnlySlice[T] {
	return &ReadOnlySlice[T]{
		slice: slice,
	}
}

func (r *ReadOnlySlice[T]) Get() []T {
	return r.slice
}
