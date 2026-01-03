package schema

import (
	"fmt"
	"reflect"
)

func BindMapEntry[K comparable, T any](unknown any, key K, target *T, errs *[]error) {
	table, ok := unknown.(map[K]any)
	if !ok {
		*errs = append(*errs, ErrNotAMap)
		return
	}

	value, ok := table[key]
	if !ok {
		*errs = append(*errs, fmt.Errorf("%w: %v", ErrMissingKey, key))
		return
	}

	castValue, ok := value.(T)
	if !ok {
		*errs = append(*errs, fmt.Errorf("%w: failed to cast map value (actual: %s, target: %s)", ErrWrongType, reflect.TypeOf(value).String(), reflect.TypeFor[T]().String()))
	}

	*target = castValue
}

func BindSlice[T any](unknown any, target *[]T, errs *[]error) {
	slice, ok := unknown.([]any)
	if !ok {
		*errs = append(*errs, ErrNotASlice)
		return
	}

	*target = make([]T, len(slice))

	for index, x := range slice {
		cast, ok := x.(T)
		if !ok {
			err := fmt.Errorf("%w, index %d, type %s", ErrWrongType, index, reflect.TypeOf(x).String())
			*errs = append(*errs, err)
		} else {
			(*target)[index] = cast
		}
	}
}
