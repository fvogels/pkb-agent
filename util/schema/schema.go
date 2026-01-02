package schema

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
