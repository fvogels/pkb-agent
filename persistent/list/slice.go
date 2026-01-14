package list

func Slice[T any](list List[T], start int, end int) List[T] {
	return &slice[T]{
		original: list,
		start:    start,
		end:      end,
	}
}

type slice[T any] struct {
	original List[T]
	start    int
	end      int
}

func (list *slice[T]) Size() int {
	return list.end - list.start
}

func (list *slice[T]) At(index int) T {
	if index < 0 || index >= list.Size() {
		panic("out of bounds")
	}

	return list.original.At(index - list.start)
}

func DropLast[T any](list List[T]) List[T] {
	return Slice(list, 0, list.Size()-1)
}
