package list

func Append[T any](list List[T], t T) List[T] {
	result := make([]T, list.Size()+1)

	ForEach(list, func(i int, x T) {
		result[i] = x
	})

	result[list.Size()] = t

	return FromSlice(result)
}

func ToSlice[T any](list List[T]) []T {
	result := make([]T, list.Size())

	ForEach(list, func(i int, x T) {
		result[i] = x
	})

	return result
}

func ForEach[T any](list List[T], f func(int, T)) {
	for i := 0; i != list.Size(); i++ {
		f(i, list.At(i))
	}
}
