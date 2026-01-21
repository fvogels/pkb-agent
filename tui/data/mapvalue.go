package data

type mappedValue[T any, R any] struct {
	value                  Value[T]
	transformer            func(T) R
	dirty                  bool
	cachedTransformedValue R
}

func MapValue[T any, R any](value Value[T], transformer func(T) R) Value[R] {
	result := mappedValue[T, R]{
		value:                  value,
		transformer:            transformer,
		cachedTransformedValue: transformer(value.Get()),
		dirty:                  false,
	}

	value.Observe(func() {
		result.dirty = true
	})

	return &result
}

func (value *mappedValue[T, R]) Get() R {
	if value.dirty {
		value.cachedTransformedValue = value.transformer(value.value.Get())
		value.dirty = false
	}

	return value.cachedTransformedValue
}

func (value *mappedValue[T, R]) Observe(observer func()) {
	value.value.Observe(observer)
}

func (value *mappedValue[T, R]) Version() uint {
	// TODO
	return 0
}

type mappedValue2[T1, T2, R any] struct {
	value1           Value[T1]
	value2           Value[T2]
	transformer      func(T1, T2) R
	transformedValue R
	dirty            bool
}

func MapValue2[T1, T2, R any](value1 Value[T1], value2 Value[T2], transformer func(T1, T2) R) Value[R] {
	result := mappedValue2[T1, T2, R]{
		value1:           value1,
		value2:           value2,
		transformer:      transformer,
		transformedValue: transformer(value1.Get(), value2.Get()),
		dirty:            false,
	}

	observer := func() {
		result.dirty = true
	}

	value1.Observe(observer)
	value2.Observe(observer)

	return &result
}

func (value *mappedValue2[T1, T2, R]) Get() R {
	if value.dirty {
		value.transformedValue = value.transformer(value.value1.Get(), value.value2.Get())
		value.dirty = false
	}
	return value.transformedValue
}

func (value *mappedValue2[T1, T2, R]) Observe(observer func()) {
	value.value1.Observe(observer)
	value.value2.Observe(observer)
}

func (value *mappedValue2[T1, T2, R]) Version() uint {
	// TODO
	return 0
}

type mappedValue3[T1, T2, T3, R any] struct {
	value1           Value[T1]
	value2           Value[T2]
	value3           Value[T3]
	transformer      func(T1, T2, T3) R
	transformedValue R
	dirty            bool
}

func MapValue3[T1, T2, T3, R any](value1 Value[T1], value2 Value[T2], value3 Value[T3], transformer func(T1, T2, T3) R) Value[R] {
	result := mappedValue3[T1, T2, T3, R]{
		value1:           value1,
		value2:           value2,
		value3:           value3,
		transformer:      transformer,
		transformedValue: transformer(value1.Get(), value2.Get(), value3.Get()),
		dirty:            false,
	}

	observer := func() {
		result.dirty = true
	}

	value1.Observe(observer)
	value2.Observe(observer)
	value3.Observe(observer)

	return &result
}

func (value *mappedValue3[T1, T2, T3, R]) Get() R {
	if value.dirty {
		value.transformedValue = value.transformer(value.value1.Get(), value.value2.Get(), value.value3.Get())
		value.dirty = false
	}

	return value.transformedValue
}

func (value *mappedValue3[T1, T2, T3, R]) Observe(observer func()) {
	value.value1.Observe(observer)
	value.value2.Observe(observer)
	value.value3.Observe(observer)
}

func (value *mappedValue3[T1, T2, T3, R]) Version() uint {
	// TODO
	return 0
}
