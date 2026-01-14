package list

func MapList[T any, U any](list List[T], f func(T) U) List[U] {
	ys := make([]U, list.Size())

	ForEach(list, func(index int, x T) {
		ys[index] = f(x)
	})

	return FromSlice(ys)
}

func MapWithIndex[T any, U any](list List[T], f func(int, T) U) List[U] {
	ys := make([]U, list.Size())

	ForEach(list, func(index int, x T) {
		ys[index] = f(index, x)
	})

	return FromSlice(ys)
}
