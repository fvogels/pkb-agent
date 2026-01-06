package data

type Data[T any] interface {
	Get() T
}
