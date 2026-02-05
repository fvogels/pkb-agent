package list

func Concatenate[T any](xss ...List[T]) List[T] {
	if len(xss) == 0 {
		return FromItems[T]()
	}

	result := xss[0]

	for index := 1; index < len(xss); index++ {
		result = &listConcat[T]{
			xs: result,
			ys: xss[index],
		}
	}

	return result
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
