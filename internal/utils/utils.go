package utils

import "fmt"

type Iter[T any] func(yield func(T) bool)

func ToString[T any](v T) string {
	return fmt.Sprintf("%+v", v)
}
