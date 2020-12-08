package common

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/romberli/go-util/constant"
)

type trimStruct struct {
	Id   int
	Name string
	B    bool
	s    string
}

func TestCommon(t *testing.T) {
	var (
		err            error
		s              string
		exists         bool
		sliceInterface []interface{}
		mapInterface   map[interface{}]interface{}
	)

	asst := assert.New(t)

	str1 := "a"
	str2 := "xxx"
	int1 := 1234
	float1 := 1.234

	sliceInt := []int{1, 2, 3}
	sliceStr := []string{"a", "b", "c"}

	mapStrInt := map[string]int{"a": 1, "b": 2, "c": 3}
	mapStrStr := map[string]string{"a": "xxx", "b": "yyy", "c": "zzz"}

	ts := &trimStruct{
		1,
		"    a    b   ",
		false,
		"             s    ",
	}

	// test ConvertInterfaceToString()
	t.Log("==========test ConvertInterfaceToString() started==========")
	s, err = ConvertNumberToString(str1)
	asst.Nil(err, "convert string to string failed")
	asst.Equal(str1, s, "convert string to string failed")

	intStr := "1234"
	s, err = ConvertNumberToString(int1)
	asst.Nil(err, "convert string to string failed")
	asst.Equal(intStr, s, "convert string to string failed")

	floatStr := "1.234"
	s, err = ConvertNumberToString(float1)
	asst.Nil(err, "convert string to string failed")
	asst.Equal(floatStr, s, "convert string to string failed")

	s, err = ConvertNumberToString(true)
	asst.Nil(err, "convert string to string failed")
	asst.Equal(constant.TrueString, s, "convert string to string failed")
	t.Log("==========test ConvertInterfaceToString() completed==========")

	// test ConvertInterfaceToSliceInterface()
	t.Log("==========test ConvertInterfaceToSliceInterface() started==========")
	sliceInterface, err = ConvertInterfaceToSliceInterface(sliceInt)
	asst.Nil(err, "test ConvertInterfaceToSliceInterface sliceInt failed")
	for _, v := range sliceInterface {
		switch v.(type) {
		case interface{}:
		default:
			asst.True(false, "test ConvertInterfaceToSliceInterface sliceInt failed")
		}
	}
	t.Logf("sliceInt convert to %v", sliceInterface)

	sliceInterface, err = ConvertInterfaceToSliceInterface(sliceStr)
	asst.Nil(err, "test ConvertInterfaceToSlice sliceStr failed")
	for _, v := range sliceInterface {
		switch v.(type) {
		case interface{}:
		default:
			asst.True(false, "test ConvertInterfaceToSliceInterface sliceStr failed")
		}
	}
	t.Logf("sliceStr convert to %v", sliceInterface)
	t.Log("==========test ConvertInterfaceToSliceInterface() completed==========")

	// test ConvertInterfaceToMapInterfaceInterface()
	t.Log("==========test ConvertInterfaceToMapInterfaceInterface() started==========")
	mapInterface, err = ConvertInterfaceToMapInterfaceInterface(mapStrInt)
	asst.Nil(err, "test ConvertInterfaceToMapInterfaceInterface mapStrInt failed")
	for k, v := range mapInterface {
		switch k.(type) {
		case interface{}:
		default:
			asst.True(false, "test ConvertInterfaceToMapInterfaceInterface mapStrInt failed")
		}

		switch v.(type) {
		case interface{}:
		default:
			asst.True(false, "test ConvertInterfaceToMapInterfaceInterface mapStrInt failed")
		}
	}
	t.Logf("mapStrInt convert to %v", mapInterface)

	mapInterface, err = ConvertInterfaceToMapInterfaceInterface(mapStrStr)
	asst.Nil(err, "test ConvertInterfaceToMapInterfaceInterface mapStrStr failed")
	for k, v := range mapInterface {
		switch k.(type) {
		case interface{}:
		default:
			asst.True(false, "test ConvertInterfaceToMapInterfaceInterface mapStrStr failed")
		}

		switch v.(type) {
		case interface{}:
		default:
			asst.True(false, "test ConvertInterfaceToMapInterfaceInterface mapStrStr failed")
		}
	}
	t.Logf("mapStrStr convert to %v", mapInterface)
	t.Log("==========test ConvertInterfaceToMapInterfaceInterface() completed==========")

	t.Log("==========test ElementInSlice() started==========")
	exists, err = ElementInSlice(str1, sliceStr)
	asst.Nil(err, "test ElementInSlice str1 sliceStr failed")
	asst.True(exists, "test ElementInSlice str1 sliceStr failed")

	exists, err = ElementInSlice(str1, sliceInt)
	asst.Nil(err, "test ElementInSlice str1 sliceInt failed")
	asst.False(exists, "test ElementInSlice str1 sliceInt failed")

	exists, err = ElementInSlice(str2, sliceStr)
	asst.Nil(err, "test ElementInSlice str1 failed")
	asst.False(exists, "test ElementInSlice str2 failed")
	t.Log("==========test ElementInSlice() completed==========")

	t.Log("==========test KeyInMap() started==========")
	exists, err = KeyInMap(str1, mapStrInt)
	asst.Nil(err, "test KeyInMap str1 mapStrInt failed")
	asst.True(exists, "test ElementInSlice str1 mapStrInt failed")

	exists, err = KeyInMap(str2, mapStrStr)
	asst.Nil(err, "test ElementInSlice str1 failed")
	asst.False(exists, "test ElementInSlice str2 failed")
	t.Log("==========test KeyInMap() completed==========")

	t.Log("==========test ValueInMap() started==========")
	exists, err = ValueInMap(str1, mapStrInt)
	asst.Nil(err, "test ValueInMap str1 mapStrInt failed")
	asst.False(exists, "test ValueInMap str1 mapStrInt failed")

	exists, err = ValueInMap(str1, mapStrStr)
	asst.Nil(err, "test ValueInMap str1 mapStrStr failed")
	asst.False(exists, "test ValueInMap str2 mapStrStr failed")

	exists, err = ValueInMap(str2, mapStrStr)
	asst.Nil(err, "test ValueInMap str2 mapStrStr failed")
	asst.True(exists, "test ValueInMap str2 mapStrStr failed")
	t.Log("==========test ValueInMap() completed==========")

	t.Log("==========test TrimSpaceOfStructString() started==========")
	t.Logf("old ts: %v", *ts)
	err = TrimSpaceOfStructString(ts)
	asst.Nil(err, "test TrimSpaceOfStructString failed")
	t.Logf("new ts: %v", *ts)
	t.Log("==========test TrimSpaceOfStructString() completed==========")

}
