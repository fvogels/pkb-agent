package list

func Concatenate[T any](xs List[T], ys List[T]) List[T] {
	result := listConcat[T]{
		xs: xs,
		ys: ys,
	}

	return &result
}

type listConcat[T any] struct {
	xs List[T]
	ys List[T]
}

func (list *listConcat[T]) Size() int {
	return list.xs.Size() + list.ys.Size()
}

func (list *listConcat[T]) At(index int) T {
	if index < 0 || index >= list.Size() {
		panic("out of bounds")
	}

	if index < list.xs.Size() {
		return list.xs.At(index)
	} else {
		return list.ys.At(index - list.xs.Size())
	}
}
