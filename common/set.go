package common

import "github.com/romberli/go-util/types"

// UnionSlice returns the union of two slices
func UnionSlice[T types.Primitive](first, second []T) []T {
	var result []T

	for _, v := range first {
		result = append(result, v)
	}

	for _, v := range second {
		if !ElementInSlice(result, v) {
			result = append(result, v)
		}
	}

	return result
}

// IntersectSlice returns the intersection of two slices
func IntersectSlice[T types.Primitive](first, second []T) []T {
	var result []T

	for _, v := range first {
		if ElementInSlice(second, v) {
			result = append(result, v)
		}
	}

	return result
}

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
