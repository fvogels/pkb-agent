package data

type List[T any] interface {
	Size() int
	At(index int) T
	Observe(func())
}
