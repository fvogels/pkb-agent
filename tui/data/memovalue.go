package data

type memoizedValue[T any] struct {
	original      Value[T]
	cached        T
	cachedVersion uint
}

func MemoizeValue[T any](value Value[T]) Value[T] {
	result := memoizedValue[T]{
		original:      value,
		cached:        value.Get(),
		cachedVersion: value.Version(),
	}

	return &result
}

func (value *memoizedValue[T]) Get() T {
	if value.cachedVersion != value.original.Version() {
		value.refresh()
	}

	return value.cached
}

func (value *memoizedValue[T]) refresh() {
	value.cached = value.original.Get()
	value.cachedVersion = value.original.Version()
}

func (value *memoizedValue[T]) Observe(observer func()) {
	value.original.Observe(observer)
}

func (value *memoizedValue[T]) Version() uint {
	// TODO
	return 0
}
