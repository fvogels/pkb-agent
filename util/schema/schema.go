package schema

import "fmt"

func BindMapEntry[K comparable, T any](unknown any, key K, target *T, errs *[]error) {
	table, ok := unknown.(map[K]any)
	if !ok {
		*errs = append(*errs, ErrNotAMap)
		return
	}

	value, ok := table[key]
	if !ok {
		*errs = append(*errs, ErrMissingKey)
		return
	}

	castValue, ok := value.(T)
	if !ok {
		*errs = append(*errs, ErrWrongType)
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
			err := fmt.Errorf("%w, index %d", ErrWrongType, index)
			*errs = append(*errs, err)
		} else {
			(*target)[index] = cast
		}
	}
}
