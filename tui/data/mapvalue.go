package data

func MapValue[T any, U any](value Value[T], transformer func(T) U) Value[U] {
	result := mappedValue[T, U]{
		originalValue:    value,
		transformedValue: transformer(value.Get()),
	}

	value.Observe(func() {
		result.transformedValue = transformer(value.Get())
	})

	return &result
}

type mappedValue[T any, U any] struct {
	originalValue    Value[T]
	transformedValue U
}

func (value *mappedValue[T, U]) Get() U {
	return value.transformedValue
}

func (value *mappedValue[T, U]) Observe(observer func()) {
	value.originalValue.Observe(observer)
}

func MapValue2[T1, T2, R any](value1 Value[T1], value2 Value[T2], transformer func(T1, T2) R) Value[R] {
	result := mappedValue2[T1, T2, R]{
		originalValue1:   value1,
		originalValue2:   value2,
		transformedValue: transformer(value1.Get(), value2.Get()),
	}

	observer := func() {
		result.transformedValue = transformer(value1.Get(), value2.Get())
	}

	value1.Observe(observer)
	value2.Observe(observer)

	return &result
}

type mappedValue2[T1, T2, R any] struct {
	originalValue1   Value[T1]
	originalValue2   Value[T2]
	transformedValue R
}

func (value *mappedValue2[T1, T2, R]) Get() R {
	return value.transformedValue
}

func (value *mappedValue2[T1, T2, R]) Observe(observer func()) {
	value.originalValue1.Observe(observer)
	value.originalValue2.Observe(observer)
}

func MapValue3[T1, T2, T3, R any](value1 Value[T1], value2 Value[T2], value3 Value[T3], transformer func(T1, T2, T3) R) Value[R] {
	result := mappedValue3[T1, T2, T3, R]{
		originalValue1:   value1,
		originalValue2:   value2,
		originalValue3:   value3,
		transformedValue: transformer(value1.Get(), value2.Get(), value3.Get()),
	}

	observer := func() {
		result.transformedValue = transformer(value1.Get(), value2.Get(), value3.Get())
	}

	value1.Observe(observer)
	value2.Observe(observer)
	value3.Observe(observer)

	return &result
}

type mappedValue3[T1, T2, T3, R any] struct {
	originalValue1   Value[T1]
	originalValue2   Value[T2]
	originalValue3   Value[T3]
	transformedValue R
}

func (value *mappedValue3[T1, T2, T3, R]) Get() R {
	return value.transformedValue
}

func (value *mappedValue3[T1, T2, T3, R]) Observe(observer func()) {
	value.originalValue1.Observe(observer)
	value.originalValue2.Observe(observer)
	value.originalValue3.Observe(observer)
}
