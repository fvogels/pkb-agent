package data

import "pkb-agent/tui"

type Value[T any] interface {
	Get() T
}

type List[T any] interface {
	Size() int
	At(index int) T
}

func Sublist[T any](list List[T], startIndex int, length int) List[T] {
	return &sublist[T]{
		original:   list,
		startIndex: startIndex,
		length:     length,
	}
}

type sublist[T any] struct {
	original   List[T]
	startIndex int
	length     int
}

func (list *sublist[T]) Size() int {
	return list.length
}

func (list *sublist[T]) At(index int) T {
	if tui.SafeMode && !list.isValidIndex(index) {
		panic("invalid index")
	}

	return list.original.At(index - list.startIndex)
}

func (list *sublist[T]) isValidIndex(index int) bool {
	return 0 <= index && index < list.length && list.startIndex+index < list.original.Size()
}
