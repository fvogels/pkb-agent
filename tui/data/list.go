package data

type List[T any] interface {
	Size() int
	At(index int) T
	Observe(func())
}

func CopyListToSlice[T any](list List[T]) []T {
	result := make([]T, list.Size())

	for i := range list.Size() {
		result[i] = list.At(i)
	}

	return result
}
