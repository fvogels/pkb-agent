package data

func MemoizeValue[T any](value Value[T]) Value[T] {
	result := memoizedValue[T]{
		original: value,
		cached:   value.Get(),
	}

	value.Observe(func() { result.refresh() })

	return &result
}

type memoizedValue[T any] struct {
	original Value[T]
	cached   T
}

func (value *memoizedValue[T]) Get() T {
	return value.cached
}

func (value *memoizedValue[T]) refresh() {
	value.cached = value.original.Get()
}

func (value *memoizedValue[T]) Observe(observer func()) {
	value.original.Observe(observer)
}

func (value *memoizedValue[T]) Version() uint {
	// TODO
	return 0
}
