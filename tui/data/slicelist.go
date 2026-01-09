package data

type SliceList[T any] struct {
	items     []T
	observers []func()
}

func NewSliceList[T any](items []T) *SliceList[T] {
	return &SliceList[T]{
		items: items,
	}
}

func (list *SliceList[T]) Size() int {
	return len(list.items)
}

func (list *SliceList[T]) At(index int) T {
	return list.items[index]
}

func (list *SliceList[T]) Observe(observer func()) {
	list.observers = append(list.observers, observer)
}

func (list *SliceList[T]) SetSlice(items []T) {
	list.items = items
	list.notifyObservers()
}

func (list *SliceList[T]) notifyObservers() {
	for _, observer := range list.observers {
		observer()
	}
}
