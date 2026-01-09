package data

func MapList[T any, U any](list List[T], f func(T) U) List[U] {
	return &mappedList[T, U]{
		original: list,
		f:        f,
	}
}

type mappedList[T any, U any] struct {
	original List[T]
	f        func(T) U
}

func (list *mappedList[T, U]) Size() int {
	return list.original.Size()
}

func (list *mappedList[T, U]) At(index int) U {
	return list.f(list.original.At(index))
}

func (list *mappedList[T, U]) Observe(func()) {}
