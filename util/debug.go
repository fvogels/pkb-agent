package util

import "fmt"

func AddressOf[T any](value *T) string {
	return fmt.Sprintf("%p", value)
}
