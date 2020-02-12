package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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

	// test ConvertSliceToInterface()
	t.Log("==========test ConvertSliceToInterface() started==========")
	sliceInterface, err = ConvertSliceToInterface(sliceInt)
	assert.Nil(err, "test ConvertSliceToInterface sliceInt failed")
	for _, v := range sliceInterface {
		switch v.(type) {
		case interface{}:
		default:
			assert.True(false, "test ConvertSliceToInterface sliceInt failed")
		}
	}
	t.Logf("sliceInt convert to %v", sliceInt)

	sliceInterface, err = ConvertSliceToInterface(sliceStr)
	assert.Nil(err, "test ConvertSliceToInterface sliceStr failed")
	for _, v := range sliceInterface {
		switch v.(type) {
		case interface{}:
		default:
			assert.True(false, "test ConvertSliceToInterface sliceStr failed")
		}
	}
	t.Logf("sliceStr convert to %v", mapStrStr)
	t.Log("==========test StringInSlice() completed==========")

	// test ConvertMapToInterface()
	t.Log("==========test ConvertMapToInterface() started==========")
	mapInterface, err = ConvertMapToInterface(mapStrInt)
	assert.Nil(err, "test ConvertMapToInterface mapStrInt failed")
	for k, v := range mapInterface {
		switch k.(type) {
		case interface{}:
		default:
			assert.True(false, "test ConvertMapToInterface mapStrInt failed")
		}

		switch v.(type) {
		case interface{}:
		default:
			assert.True(false, "test ConvertMapToInterface mapStrInt failed")
		}
	}
	t.Logf("mapStrInt convert to %v", mapInterface)

	mapInterface, err = ConvertMapToInterface(mapStrStr)
	assert.Nil(err, "test ConvertMapToInterface mapStrStr failed")
	for k, v := range mapInterface {
		switch k.(type) {
		case interface{}:
		default:
			assert.True(false, "test ConvertMapToInterface mapStrStr failed")
		}

		switch v.(type) {
		case interface{}:
		default:
			assert.True(false, "test ConvertMapToInterface mapStrStr failed")
		}
	}
	t.Logf("mapStrStr convert to %v", mapInterface)
	t.Log("==========test ConvertMapToInterface() completed==========")

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
}
