package data

type mappedValue[T any, R any] struct {
	originalValue          Value[T]
	transformer            func(T) R
	dirty                  bool
	cachedTransformedValue R
}

func MapValue[T any, R any](value Value[T], transformer func(T) R) Value[R] {
	result := mappedValue[T, R]{
		originalValue:          value,
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
		value.cachedTransformedValue = value.transformer(value.originalValue.Get())
		value.dirty = false
	}

	return value.cachedTransformedValue
}

func (value *mappedValue[T, R]) Observe(observer func()) {
	value.originalValue.Observe(observer)
}

func (value *mappedValue[T, R]) Version() uint {
	// TODO
	return 0
}

type mappedValue2[T1, T2, R any] struct {
	originalValue1   Value[T1]
	originalValue2   Value[T2]
	transformer      func(T1, T2) R
	transformedValue R
	dirty            bool
}

func MapValue2[T1, T2, R any](value1 Value[T1], value2 Value[T2], transformer func(T1, T2) R) Value[R] {
	result := mappedValue2[T1, T2, R]{
		originalValue1:   value1,
		originalValue2:   value2,
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
		value.transformedValue = value.transformer(value.originalValue1.Get(), value.originalValue2.Get())
		value.dirty = false
	}
	return value.transformedValue
}

func (value *mappedValue2[T1, T2, R]) Observe(observer func()) {
	value.originalValue1.Observe(observer)
	value.originalValue2.Observe(observer)
}

func (value *mappedValue2[T1, T2, R]) Version() uint {
	// TODO
	return 0
}

type mappedValue3[T1, T2, T3, R any] struct {
	originalValue1   Value[T1]
	originalValue2   Value[T2]
	originalValue3   Value[T3]
	transformer      func(T1, T2, T3) R
	transformedValue R
	dirty            bool
}

func MapValue3[T1, T2, T3, R any](value1 Value[T1], value2 Value[T2], value3 Value[T3], transformer func(T1, T2, T3) R) Value[R] {
	result := mappedValue3[T1, T2, T3, R]{
		originalValue1:   value1,
		originalValue2:   value2,
		originalValue3:   value3,
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
		value.transformedValue = value.transformer(value.originalValue1.Get(), value.originalValue2.Get(), value.originalValue3.Get())
		value.dirty = false
	}

	return value.transformedValue
}

func (value *mappedValue3[T1, T2, T3, R]) Observe(observer func()) {
	value.originalValue1.Observe(observer)
	value.originalValue2.Observe(observer)
	value.originalValue3.Observe(observer)
}

func (value *mappedValue3[T1, T2, T3, R]) Version() uint {
	// TODO
	return 0
}
