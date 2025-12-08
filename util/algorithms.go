package util

func Map[T any, R any](xs []T, transformer func(t T) R) []R {
	result := make([]R, len(xs))

	for index, x := range xs {
		result[index] = transformer(x)
	}

	return result
}

func Filter[T any](xs []T, predicate func(t T) bool) []T {
	result := []T{}

	for _, x := range xs {
		if predicate(x) {
			result = append(result, x)
		}
	}

	return result
}

func FindIndex[T any](xs []T, predicate func(t T) bool) int {
	for index, x := range xs {
		if predicate(x) {
			return index
		}
	}

	return -1
}

func Compose[T, U, R any](f func(T) U, g func(U) R) func(T) R {
	return func(t T) R {
		return g(f(t))
	}
}

func Curry2[T1, T2, R any](f func(T1, T2) R, x T1) func(T2) R {
	return func(y T2) R {
		return f(x, y)
	}
}

func IsLowercaseLetter(char byte) bool {
	return 'a' <= char && char <= 'z'
}

func IsDigit(char byte) bool {
	return '0' <= char && char <= '9'
}

func LowercaseWords(str string) []string {
	str = strings.ToLower(str)
	words := []string{}
	i := 0

	for i < len(str) {
		j := i
		for j != len(str) && IsLowercaseLetter(str[j]) {
			j++
		}

		word := str[i:j]
		words = append(words, word)
		i = j + 1
	}

	return words
}
