package data

type variable[T any] struct {
	value T
}

func NewVariable[T any](value T) Data[T] {
	result := variable[T]{value: value}

	return &result
}

func (v *variable[T]) Get() T {
	return v.value
}
