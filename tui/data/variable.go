package data

type Variable[T any] struct {
	value     T
	observers []func()
}

func NewVariable[T any](value T) Variable[T] {
	result := Variable[T]{value: value}

	return result
}

func (v *Variable[T]) Set(value T) {
	v.value = value

	for _, observer := range v.observers {
		observer()
	}
}

func (v *Variable[T]) Get() T {
	return v.value
}

func (v *Variable[T]) Observe(observer func()) {
	v.observers = append(v.observers, observer)
}

func (v *Variable[T]) Version() uint {
	// TODO
	return 0
}
