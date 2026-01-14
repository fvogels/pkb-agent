package list

type MaterialList[T any] struct {
	items []T
}

func New[T any]() List[T] {
	return &MaterialList[T]{
		items: nil,
	}
}

func FromSlice[T any](items []T) List[T] {
	return &MaterialList[T]{
		items: items,
	}
}

func (list MaterialList[T]) Size() int {
	return len(list.items)
}

func (list MaterialList[T]) At(index int) T {
	return list.items[index]
}
