package listview

type emptyList[T any] struct{}

func (list *emptyList[T]) At(index int) T {
	panic("cannot index empty list")
}

func (list *emptyList[T]) Length() int {
	return 0
}
