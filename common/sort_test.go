package common

import (
	"reflect"
	"testing"

	"github.com/pingcap/errors"
	"github.com/stretchr/testify/assert"

	"github.com/romberli/go-util/constant"
)

var (
	_ Comparable = (*testComparable)(nil)
)

type testComparable struct {
	Value int
}

func (tc *testComparable) Compare(other Comparable) (CompareResult, error) {
	oth, ok := other.(*testComparable)
	if !ok {
		return constant.ZeroInt, errors.Errorf("type assertion failed, expectType: *testComparable, actualType: %s",
			reflect.TypeOf(other).Kind().String())
	}

	if tc.Value == oth.Value {
		return CompareResultEqual, nil
	}

	if tc.Value > oth.Value {
		return CompareResultGT, nil
	}

	return CompareResultLT, nil
}

func TestComparable_All(t *testing.T) {
	TestQuickSortAsc(t)
	TestQuickSortDesc(t)
}

func TestQuickSortAsc(t *testing.T) {
	asst := assert.New(t)
	comparableList := []Comparable{
		&testComparable{Value: 2},
		&testComparable{Value: 1},
		&testComparable{Value: 3},
	}

	err := QuickSort(comparableList)
	asst.Nil(err, "test QuickSortAsc() failed")
	for _, comparer := range comparableList {
		t.Logf("asc: %d", comparer.(*testComparable).Value)
	}
}

func TestQuickSortDesc(t *testing.T) {
	asst := assert.New(t)
	comparableList := []Comparable{
		&testComparable{Value: 2},
		&testComparable{Value: 1},
		&testComparable{Value: 3},
	}

	err := QuickSortDesc(comparableList)
	asst.Nil(err, "test TestQuickSortDesc() failed")
	for _, comparer := range comparableList {
		t.Logf("desc: %d", comparer.(*testComparable).Value)
	}
}
