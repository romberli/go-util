package common

import (
	"github.com/romberli/go-util/constant"
)

const (
	CompareResultEqual CompareResult = 0
	CompareResultGT    CompareResult = 1
	CompareResultLT    CompareResult = -1
)

type CompareResult int

type Comparable interface {
	Compare(other Comparable) (CompareResult, error)
}

func QuickSort(comparableList []Comparable) error {
	return QuickSortAsc(comparableList)
}

func QuickSortAsc(comparableList []Comparable) error {
	if len(comparableList) == constant.ZeroInt {
		return nil
	}

	return quickSort(comparableList, constant.ZeroInt, len(comparableList)-constant.OneInt, CompareResultLT)
}

func QuickSortDesc(comparableList []Comparable) error {
	if len(comparableList) == constant.ZeroInt {
		return nil
	}

	return quickSort(comparableList, constant.ZeroInt, len(comparableList)-constant.OneInt, CompareResultGT)
}

func quickSort(comparableList []Comparable, left, right int, cr CompareResult) error {
	if left >= right {
		return nil
	}

	var err error
	pivot := left
	index := pivot + 1
	i := index
	for i <= right {
		result, err := comparableList[i].Compare(comparableList[pivot])
		if err != nil {
			return err
		}
		if result == cr {
			comparableList[i], comparableList[index] = comparableList[index], comparableList[i]
			index++
		}
		i++
	}
	comparableList[pivot], comparableList[index-constant.OneInt] = comparableList[index-constant.OneInt], comparableList[pivot]
	pivot = index - 1

	err = quickSort(comparableList, left, pivot-constant.OneInt, cr)
	if err != nil {
		return err
	}
	err = quickSort(comparableList, pivot+constant.OneInt, right, cr)
	if err != nil {
		return err
	}

	return nil
}
