package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
		exists         bool
		sliceInterface []interface{}
		mapInterface   map[interface{}]interface{}
	)

	assert := assert.New(t)

	str1 := "a"
	str2 := "xxx"

	sliceInt := []int{1, 2, 3}
	sliceStr := []string{"a", "b", "c"}

	mapStrInt := map[string]int{"a": 1, "b": 2, "c": 3}
	mapStrStr := map[string]string{"a": "xxx", "b": "yyy", "c": "zzz"}

	testRemote := false
	pathExists := "common.go"
	pathNotExists := "not_exists.go"
	hostIpRemote := "192.168.137.11"

	ts := &trimStruct{
		1,
		"    a    b   ",
		false,
		"             s    ",
	}

	// test ConvertInterfaceToSliceInterface()
	t.Log("==========test ConvertInterfaceToSliceInterface() started==========")
	sliceInterface, err = ConvertInterfaceToSliceInterface(sliceInt)
	assert.Nil(err, "test ConvertInterfaceToSliceInterface sliceInt failed")
	for _, v := range sliceInterface {
		switch v.(type) {
		case interface{}:
		default:
			assert.True(false, "test ConvertInterfaceToSliceInterface sliceInt failed")
		}
	}
	t.Logf("sliceInt convert to %v", sliceInterface)

	sliceInterface, err = ConvertInterfaceToSliceInterface(sliceStr)
	assert.Nil(err, "test ConvertInterfaceToSlice sliceStr failed")
	for _, v := range sliceInterface {
		switch v.(type) {
		case interface{}:
		default:
			assert.True(false, "test ConvertInterfaceToSliceInterface sliceStr failed")
		}
	}
	t.Logf("sliceStr convert to %v", sliceInterface)
	t.Log("==========test ConvertInterfaceToSliceInterface() completed==========")

	// test ConvertInterfaceToMapInterfaceInterface()
	t.Log("==========test ConvertInterfaceToMapInterfaceInterface() started==========")
	mapInterface, err = ConvertInterfaceToMapInterfaceInterface(mapStrInt)
	assert.Nil(err, "test ConvertInterfaceToMapInterfaceInterface mapStrInt failed")
	for k, v := range mapInterface {
		switch k.(type) {
		case interface{}:
		default:
			assert.True(false, "test ConvertInterfaceToMapInterfaceInterface mapStrInt failed")
		}

		switch v.(type) {
		case interface{}:
		default:
			assert.True(false, "test ConvertInterfaceToMapInterfaceInterface mapStrInt failed")
		}
	}
	t.Logf("mapStrInt convert to %v", mapInterface)

	mapInterface, err = ConvertInterfaceToMapInterfaceInterface(mapStrStr)
	assert.Nil(err, "test ConvertInterfaceToMapInterfaceInterface mapStrStr failed")
	for k, v := range mapInterface {
		switch k.(type) {
		case interface{}:
		default:
			assert.True(false, "test ConvertInterfaceToMapInterfaceInterface mapStrStr failed")
		}

		switch v.(type) {
		case interface{}:
		default:
			assert.True(false, "test ConvertInterfaceToMapInterfaceInterface mapStrStr failed")
		}
	}
	t.Logf("mapStrStr convert to %v", mapInterface)
	t.Log("==========test ConvertInterfaceToMapInterfaceInterface() completed==========")

	t.Log("==========test ElementInSlice() started==========")
	exists, err = ElementInSlice(str1, sliceStr)
	assert.Nil(err, "test ElementInSlice str1 sliceStr failed")
	assert.True(exists, "test ElementInSlice str1 sliceStr failed")

	exists, err = ElementInSlice(str1, sliceInt)
	assert.Nil(err, "test ElementInSlice str1 sliceInt failed")
	assert.False(exists, "test ElementInSlice str1 sliceInt failed")

	exists, err = ElementInSlice(str2, sliceStr)
	assert.Nil(err, "test ElementInSlice str1 failed")
	assert.False(exists, "test ElementInSlice str2 failed")
	t.Log("==========test ElementInSlice() completed==========")

	t.Log("==========test KeyInMap() started==========")
	exists, err = KeyInMap(str1, mapStrInt)
	assert.Nil(err, "test KeyInMap str1 mapStrInt failed")
	assert.True(exists, "test ElementInSlice str1 mapStrInt failed")

	exists, err = KeyInMap(str2, mapStrStr)
	assert.Nil(err, "test ElementInSlice str1 failed")
	assert.False(exists, "test ElementInSlice str2 failed")
	t.Log("==========test KeyInMap() completed==========")

	t.Log("==========test ValueInMap() started==========")
	exists, err = ValueInMap(str1, mapStrInt)
	assert.Nil(err, "test ValueInMap str1 mapStrInt failed")
	assert.False(exists, "test ValueInMap str1 mapStrInt failed")

	exists, err = ValueInMap(str1, mapStrStr)
	assert.Nil(err, "test ValueInMap str1 mapStrStr failed")
	assert.False(exists, "test ValueInMap str2 mapStrStr failed")

	exists, err = ValueInMap(str2, mapStrStr)
	assert.Nil(err, "test ValueInMap str2 mapStrStr failed")
	assert.True(exists, "test ValueInMap str2 mapStrStr failed")
	t.Log("==========test ValueInMap() completed==========")

	t.Log("==========test PathExistsLocal() started==========")
	exists, err = PathExistsLocal(pathExists)
	assert.Nil(err, "test PathExistsLocal pathExists failed")
	assert.True(exists, "test PathExistsLocal pathExists failed")

	exists, err = PathExistsLocal(pathNotExists)
	assert.Nil(err, "test PathExistsLocal pathNotExists failed")
	assert.False(exists, "test PathExistsLocal pathNotExists failed")
	t.Log("==========test PathExistsLocal() completed==========")

	if testRemote {
		t.Log("==========test PathExistsRemote() started==========")
		sftpConn, err := NewMySftpConn(hostIpRemote)
		assert.Nil(err, "cann't connect to remote host")

		exists, err = PathExistsRemote(pathExists, sftpConn.SftpClient)
		assert.Nil(err, "test PathExistsRemote pathExists failed")
		assert.True(exists, "test PathExistsRemote pathExists failed")

		exists, err = PathExistsRemote(pathNotExists, sftpConn.SftpClient)
		assert.Nil(err, "test PathExistsRemote pathNotExists failed")
		assert.False(exists, "test PathExistsRemote pathNotExists failed")
		t.Log("==========test PathExistsRemote() completed==========")
	}

	t.Log("==========test PathExists() started==========")
	exists, err = PathExists(nil)
	assert.NotNil(err, "test PathExists pathExists failed")
	assert.False(exists, "test PathExists pathExists failed")

	exists, err = PathExists(pathExists)
	assert.Nil(err, "test PathExists pathExists failed")
	assert.True(exists, "test PathExists pathExists failed")

	exists, err = PathExists(pathNotExists)
	assert.Nil(err, "test PathExists pathNotExists failed")
	assert.False(exists, "test PathExists pathNotExists failed")

	exists, err = PathExists(pathExists, nil)
	assert.NotNil(err, "test PathExists pathExists failed")
	assert.False(exists, "test PathExists pathExists failed")

	if testRemote {
		sftpConn, err := NewMySftpConn(hostIpRemote)
		assert.Nil(err, "cann't connect to remote host")

		exists, err = PathExists(pathExists, sftpConn.SftpClient)
		assert.Nil(err, "test PathExists pathExists failed")
		assert.True(exists, "test PathExists pathExists failed")

		exists, err = PathExists(pathNotExists, sftpConn.SftpClient)
		assert.Nil(err, "test PathExists pathNotExists failed")
		assert.False(exists, "test PathExists pathNotExists failed")
	}
	t.Log("==========test PathExists() completed==========")

	t.Log("==========test TrimSpaceOfStructString() started==========")
	t.Logf("old ts: %v", *ts)
	err = TrimSpaceOfStructString(ts)
	assert.Nil(err, "test TrimSpaceOfStructString failed")
	t.Logf("new ts: %v", *ts)
	t.Log("==========test TrimSpaceOfStructString() completed==========")

}
