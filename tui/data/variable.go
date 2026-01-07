package data

type Variable[T any] struct {
	value T
}

func NewVariable[T any](value T) *Variable[T] {
	result := Variable[T]{value: value}

	return &result
}

func (v *Variable[T]) Set(value T) {
	v.value = value
}

func (v *Variable[T]) Get() T {
	return v.value
}
