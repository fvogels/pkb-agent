package data

type SliceList[T any] struct {
	items []T
}

func NewSliceList[T any](items []T) *SliceList[T] {
	return &SliceList[T]{
		items: items,
	}
}

func (slice *SliceList[T]) Size() int {
	return len(slice.items)
}

func (slice *SliceList[T]) At(index int) T {
	return slice.items[index]
}
