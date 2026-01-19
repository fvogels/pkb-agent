package list

import "strings"

func String[T any](xs List[T], f func(T) string) string {
	var builder strings.Builder

	builder.WriteString("[")

	for i := range xs.Size() {
		if i != 0 {
			builder.WriteString(", ")
		}

		x := xs.At(i)
		builder.WriteString(f(x))
	}

	builder.WriteString("]")

	return builder.String()
}
