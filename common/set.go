package common

import "github.com/romberli/go-util/types"

// SubtractSlice returns all the elements in first but not in second
func SubtractSlice[T types.Primitive](first, second []T) []T {
	var result []T

	for _, v := range first {
		if !ElementInSlice(second, v) {
			result = append(result, v)
		}
	}

	return result
}
