package common

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/romberli/go-util/constant"
)

type nestStruct struct {
	ID    int
	slice []string
}

type trimStruct struct {
	ID   int
	Name string `middleware:"name"`
	B    bool
	s    string
	NSA  *nestStruct
	NSB  *nestStruct
}

func (ts *trimStruct) GetNSA() *nestStruct {
	return ts.NSA
}

type expectStructA struct {
	ID   int
	Name string
	NSA  *nestStruct
}

type EnvInfo struct {
	ID      int    `middleware:"id"`
	EnvName string `middleware:"env_name"`
	DelFlag string `middleware:"del_flag"`
}

type testStructA struct {
	ID            int
	TestInterface TestInterface
}

type TestInterface interface {
	GetNSA() *nestStruct
}

func TestCommon(t *testing.T) {
	var (
		err            error
		s              string
		exists         bool
		sliceInterface []interface{}
		mapInterface   map[interface{}]interface{}
		ns1            *nestStruct
		ns2            *nestStruct
		ns3            *nestStruct
		ts             *trimStruct
		expectTs       *trimStruct
	)

	asst := assert.New(t)

	originStr := "ABCDEFGabcdefg1234567+-"
	reversedStr := "-+7654321gfedcbaGFEDCBA"

	str1 := "a"
	str2 := "xxx"
	int1 := 1234
	int2 := 3
	int3 := 6
	float1 := 1.234

	sliceInt := []int{1, 2, 3, 5, 7}
	sliceStr := []string{"a", "b", "c"}

	mapStrInt := map[string]int{"a": 1, "b": 2, "c": 3}
	mapStrStr := map[string]string{"a": "xxx", "b": "yyy", "c": "zzz"}

	ts = &trimStruct{
		1,
		"    a    b   ",
		false,
		"             s    ",
		nil,
		nil,
	}

	tsList := []*trimStruct{ts}

	// test ReverseString()
	t.Log("==========test ReverseString() started==========")
	s = ReverseString(originStr)
	asst.Equal(reversedStr, s, "test ReverseString() failed")
	t.Log("==========test ReverseString() completed==========")

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

	// test ConvertInterfaceToString()
	t.Log("==========test ConvertInterfaceToString() started==========")
	s = ConvertInterfaceToString(ts)
	t.Logf("convert string: %s", s)
	jsonBytes, err := json.Marshal(ts)
	asst.Nil(err, "json.Marshal failed")
	t.Logf("marshal string: %s", jsonBytes)
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
	exists = ElementInSlice(sliceStr, str1)
	asst.True(exists, "test ElementInSlice() failed")

	exists = ElementInSlice(sliceInt, int1)
	asst.False(exists, "test ElementInSlice() failed")

	exists = ElementInSlice(sliceStr, str2)
	asst.False(exists, "test ElementInSlice() failed")
	t.Log("==========test ElementInSlice() completed==========")

	t.Log("==========test ElementInSortedSlice() started==========")
	exists = ElementInSortedSlice(sliceInt, int1)
	asst.False(exists, "test ElementInSortedSlice() failed")
	exists = ElementInSortedSlice(sliceInt, int2)
	asst.True(exists, "test ElementInSortedSlice() failed")
	exists = ElementInSortedSlice(sliceInt, int3)
	asst.False(exists, "test ElementInSortedSlice() failed")
	t.Log("==========test ElementInSortedSlice() completed==========")

	t.Log("==========test ElementInSliceInterface() started==========")
	exists, err = ElementInSliceInterface(tsList, ts)
	asst.Nil(err, "test ElementInSliceInterface() failed")
	asst.True(exists, "test ElementInSliceInterface() failed")
	t.Log("==========test ElementInSliceInterface() completed==========")

	t.Log("==========test KeyInMap() started==========")
	exists, err = KeyInMap(mapStrInt, str1)
	asst.Nil(err, "test KeyInMap() failed")
	asst.True(exists, "test KeyInMap() failed")

	exists, err = KeyInMap(mapStrStr, str2)
	asst.Nil(err, "test KeyInMap() failed")
	asst.False(exists, "test KeyInMap() failed")
	t.Log("==========test KeyInMap() completed==========")

	t.Log("==========test ValueInMap() started==========")
	exists, err = ValueInMap(mapStrInt, str1)
	asst.Nil(err, "test ValueInMap() failed")
	asst.False(exists, "test ValueInMap() failed")

	exists, err = ValueInMap(mapStrStr, str1)
	asst.Nil(err, "test ValueInMap() failed")
	asst.False(exists, "test ValueInMap() failed")

	exists, err = ValueInMap(mapStrStr, str2)
	asst.Nil(err, "test ValueInMap() failed")
	asst.True(exists, "test ValueInMap() failed")
	t.Log("==========test ValueInMap() completed==========")

	t.Log("==========test TrimSpaceOfStructString() started==========")
	t.Logf("old ts: %v", *ts)
	err = TrimSpaceOfStructString(ts)
	asst.Nil(err, "test TrimSpaceOfStructString() failed")
	t.Logf("new ts: %v", *ts)
	t.Log("==========test TrimSpaceOfStructString() completed==========")

	ns1 = &nestStruct{
		ID:    100,
		slice: []string{"a", "b", "c"},
	}
	ns2 = &nestStruct{
		ID:    200,
		slice: []string{"a", "b"},
	}
	ns3 = &nestStruct{
		ID:    200,
		slice: []string{"a", "b"},
	}
	ts = &trimStruct{
		ID:   1,
		Name: "aaa",
		B:    true,
		s:    "small",
		NSA:  ns1,
		NSB:  ns1,
	}
	expectTs = &trimStruct{
		ID:  1,
		B:   false,
		s:   "big",
		NSA: ns2,
	}

	t.Log("==========test GetValueOfStruct() started==========")
	name, err := GetValueOfStruct(ts, "Name")
	asst.Nil(err, "test GetValueOfStruct() failed")
	asst.Equal("aaa", name, "test GetValueOfStruct() failed")
	t.Log("==========test GetValueOfStruct() completed==========")

	t.Log("==========test SetValueOfStruct() started==========")
	// set bool
	err = SetValueOfStruct(ts, "B", false)
	asst.Nil(err, "test SetValueOfStruct() failed")
	asst.Equal(expectTs.B, ts.B, "test SetValueOfStruct() failed")
	// set not exists field
	err = SetValueOfStruct(ts, "a", "big")
	asst.NotNil(err, "test SetValueOfStruct() failed")
	// set unexported field
	err = SetValueOfStruct(ts, "s", "big")
	asst.NotNil(err, "test SetValueOfStruct() failed")
	// set struct
	err = SetValueOfStruct(ts, "NSA", ns3)
	asst.Nil(err, "test SetValueOfStruct() failed")
	asst.Equal(expectTs.NSA.ID, ts.NSA.ID, "test SetValueOfStruct() failed")
	asst.Equal(expectTs.NSA.slice, ts.NSA.slice, "test SetValueOfStruct() failed")
	t.Logf("ts: %v", ts)
	t.Log("==========test SetValueOfStruct() completed==========")

	t.Log("==========test SetValueOfStructByTag() started==========")
	err = SetValueOfStructByTag(ts, "name", "newName", constant.DefaultMiddlewareTag)
	asst.Nil(err, "test SetValueOfStructByTag() failed")
	asst.Equal("newName", ts.Name, "test SetValueOfStructByTag() failed")
	err = SetValueOfStructByTag(ts, "name", "aaa", constant.DefaultMiddlewareTag)
	asst.Nil(err, "test SetValueOfStructByTag() failed")
	t.Log("==========test SetValueOfStructByTag() completed==========")

	t.Log("==========test SetValuesWithMapByTag() started==========")
	err = SetValuesWithMapByTag(ts, map[string]interface{}{"name": "newName"}, constant.DefaultMiddlewareTag)
	asst.Nil(err, "test SetValueOfStructByTag() failed")
	err = SetValuesWithMapByTag(ts, map[string]interface{}{"name": "aaa"}, constant.DefaultMiddlewareTag)
	asst.Nil(err, "test SetValueOfStructByTag() failed")
	t.Log("==========test SetValuesWithMapByTag() completed==========")

	t.Log("==========test CopyStructWithFields() started==========")
	es := &expectStructA{
		ID:   1,
		Name: "bbb",
		NSA:  ns2,
	}
	ns2.ID = 300
	ts.NSA = ns2
	jets, err := json.Marshal(es)
	asst.Nil(err, "test CopyStructWithFields() failed")
	nts, err := CopyStructWithFields(ts, []string{"ID", "Name", "NSA"}...)
	asst.Nil(err, "test CopyStructWithFields() failed")
	err = SetValueOfStruct(nts, "Name", "bbb")
	asst.Nil(err, "test CopyStructWithFields() failed")
	jnts, err := json.Marshal(nts)
	asst.Nil(err, "test CopyStructWithFields() failed")
	// asst.Equal(string(jets), string(jnts), "test CopyStructWithFields() failed")
	t.Log(jets, jnts)
	ti := &testStructA{
		ID:            1,
		TestInterface: ts,
	}

	ns4, err := CopyStructWithFields(ti, "TestInterface")
	asst.Nil(err, "test CopyStructWithFields() failed")
	jsonBytes, err = json.Marshal(ns4)
	t.Logf("copied field: %s", string(jsonBytes))
	asst.Nil(err, "test CopyStructWithFields() failed")

	t.Log("==========test CopyStructWithFields() completed==========")

	t.Log("==========test MarshalStructWithFields() started==========")
	es = &expectStructA{
		ID:   1,
		Name: "aaa",
		NSA:  ns2,
	}
	asst.Nil(err, "test MarshalStructWithFields() failed")
	jets, err = json.Marshal(es)
	asst.Nil(err, "test MarshalStructWithFields() failed")
	jnts, err = MarshalStructWithFields(ts, []string{"ID", "Name", "NSA"}...)
	asst.Nil(err, "test MarshalStructWithFields() failed")
	// asst.Equal(string(jets), string(jnts), "test MarshalStructWithFields() failed")
	t.Log(jets, jnts)
	t.Log("==========test MarshalStructWithFields() completed==========")

	t.Log("==========test MarshalStructWithoutFields() started==========")
	jnts, err = MarshalStructWithoutFields(ts, []string{"B", "NSB"}...)
	asst.Nil(err, "test MarshalStructWithoutFields() failed")
	// asst.Equal(string(jets), string(jnts), "test MarshalStructWithoutFields() failed")
	t.Log(jets, jnts)
	t.Log("==========test MarshalStructWithoutFields() completed==========")

	t.Log("==========test MarshalStructWithTag() started==========")
	jnts, err = MarshalStructWithTag(ts, "middleware")
	asst.Nil(err, "test MarshalStructWithTag() failed")
	asst.Equal("{\"Name\":\"aaa\"}", string(jnts), "test MarshalStructWithoutFields() failed")
	t.Log("==========test MarshalStructWithTag() completed==========")

	t.Log("==========test NewMapWithStructTag() started==========")
	envInfo := &EnvInfo{}
	oldMap := map[string]interface{}{"id": 1, "env_name": "test"}
	expectMap := map[string]interface{}{"ID": 1, "EnvName": "test"}
	newMap, err := NewMapWithStructTag(oldMap, envInfo, "middleware")
	asst.Nil(err, "test NewMapWithStructTag() failed")
	asst.True(reflect.DeepEqual(newMap, expectMap), "test NewMapWithStructTag() failed")
	t.Log("==========test NewMapWithStructTag() completed==========")

	t.Log("==========test UnmarshalToMapWithStructTagFromString() started==========")
	data, _ := json.Marshal(oldMap)
	result, err := UnmarshalToMapWithStructTagFromString(string(data), envInfo, "middleware")
	asst.Nil(err, "test UnmarshalToMapWithStructTagFromString() failed")
	asst.True(reflect.DeepEqual(result, expectMap), "test UnmarshalToMapWithStructTagFromString() failed")
	t.Log("==========test UnmarshalToMapWithStructTagFromString() completed==========")
}

func TestSetValuesWithMapAndRandom(t *testing.T) {
	asst := assert.New(t)

	ei := &EnvInfo{}
	fields := map[string]interface{}{"EnvName": "test"}

	err := SetValuesWithMapAndRandom(ei, fields)
	asst.Nil(err, "test SetValuesWithMapAndRandom() failed")
}

func TestSetValuesWithMapAndRandomByTag(t *testing.T) {
	asst := assert.New(t)

	ei := &EnvInfo{}
	fields := map[string]interface{}{"env_name": "test"}
	err := SetValuesWithMapAndRandomByTag(ei, fields, "middleware")
	asst.Nil(err, "test SetValuesWithMapAndRandomByTag() failed")
}

func TestSetValueOfStruct(t *testing.T) {
	asst := assert.New(t)

	type T int
	type StructA struct {
		T T
	}

	s := &StructA{}

	err := SetValueOfStruct(s, "T", 1)
	asst.Nil(err, "test SetValueOfStruct() failed")
}

func TestIsIsRandomValue(t *testing.T) {
	asst := assert.New(t)

	type T int

	a := T(1)

	isRandom, err := IsRandomValue(a)
	asst.Nil(err, "test IsRandomValue() failed")
	asst.False(isRandom, "test IsRandomValue() failed")

	b := T(constant.DefaultRandomInt)

	isRandom, err = IsRandomValue(b)
	asst.Nil(err, "test IsRandomValue() failed")
	asst.True(isRandom, "test IsRandomValue() failed")
}
