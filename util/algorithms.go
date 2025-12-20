package util

import (
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/unicode/norm"
)

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

func All[T any](xs []T, predicate func(t T) bool) bool {
	for _, x := range xs {
		if !predicate(x) {
			return false
		}
	}

	return true
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

func RemoveAccents(s string) string {
	t := norm.NFD.String(s)
	t = runes.Remove(runes.In(unicode.Mn)).String(t)
	return t
}

func KeepOnlyLettersAndSpaces(s string) string {
	var builder strings.Builder
	lastWasSpace := true

	for _, r := range s {
		if unicode.IsLetter(r) {
			// Returns error, but documentation says it will be nil, so we ignore it
			builder.WriteRune(r)
			lastWasSpace = false
		} else if unicode.IsSpace(r) {
			if !lastWasSpace {
				// Returns error, but documentation says it will be nil, so we ignore it
				builder.WriteRune(r)
				lastWasSpace = true
			}
		} else {
			if !lastWasSpace {
				builder.WriteRune(' ')
				lastWasSpace = true
			}
		}
	}

	return builder.String()
}

func CollectTo[T any](receiver *[]T) func(T) bool {
	return func(item T) bool {
		*receiver = append(*receiver, item)
		return true
	}
}

func MaxInt(a int, b int) int {
	if a >= b {
		return a
	} else {
		return b
	}
}

func MinInt(a int, b int) int {
	if a <= b {
		return a
	} else {
		return b
	}
}

func SplitInLines(str string) []string {
	lines := []string{}
	lineGenerator := strings.Lines(str)

	lineGenerator(func(line string) bool {
		trimmedLine := strings.TrimRight(line, "\t\n\r ")
		lines = append(lines, trimmedLine)

		return true
	})

	return lines
}
