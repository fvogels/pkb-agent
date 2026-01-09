package data

func MapValue[T any, U any](value Value[T], transformer func(T) U) Value[U] {
	result := mappedValue[T, U]{
		originalValue: value,
		transformer:   transformer,
	}

	return &result
}

type mappedValue[T any, U any] struct {
	originalValue Value[T]
	transformer   func(T) U
}

func (value *mappedValue[T, U]) Get() U {
	return value.transformer(value.originalValue.Get())
}

func (value *mappedValue[T, U]) Observe(observer func()) {
	value.originalValue.Observe(observer)
}
