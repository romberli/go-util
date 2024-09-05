package linux

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFile(t *testing.T) {
	var (
		err                error
		dirExists          string
		dirNotExists       string
		dirName            string
		empDirName         string
		fileName           string
		pathExists         bool
		isDir              bool
		isEmpty            bool
		fileInfoList       []os.FileInfo
		expectedFileNum    int
		expectedElementNum int
		pathDirMap         map[string]string
		fileNameDest       string
		expectedFileName   string
	)

	asst := assert.New(t)

	dirExists = "/Users/romber"
	dirNotExists = "/Users/romber/xxxxdasfdsafas"
	dirName = "/Users/romber/test"
	fileName = "/Users/romber/test/1.txt"
	expectedFileNum = 4
	empDirName = "/Users/romber/test/subdir2"
	expectedElementNum = 5
	expectedFileName = "/Users/romber/test/subdir2/1.txt"

	// test PathExists()
	t.Log("==========test PathExists() started.==========")
	pathExists, err = PathExists(dirExists)
	asst.Nil(err, "test PathExists() failed")
	asst.True(pathExists, "test PathExists() failed")

	pathExists, err = PathExists(dirNotExists)
	asst.Nil(err, "test PathExists() failed")
	asst.False(pathExists, "test PathExists() failed")
	t.Log("==========test PathExists() completed.==========\n")

	// test IsDir()
	t.Log("==========test IsDir() started.==========")
	isDir, err = IsDir(dirName)
	asst.Nil(err, "test IsDir() failed")
	asst.True(isDir, "test IsDir() failed")

	isDir, err = IsDir(fileName)
	asst.Nil(err, "test IsDir() failed")
	asst.False(isDir, "test IsDir() failed")
	t.Log("==========test IsDir() completed.==========\n")

	// test ReadDir()
	t.Log("==========test ReadDir() started.==========")
	fileInfoList, err = ReadDir(dirName)
	asst.Nil(err, "test ReadDir() failed")
	asst.Equal(len(fileInfoList), expectedFileNum, "test ReadDir() failed")

	fileInfoList, err = ReadDir(fileName)
	asst.NotNil(err, "test ReadDir() failed")
	t.Log("==========test ReadDir() completed.==========\n")

	// test IsEmptyDir()
	t.Log("==========test IsEmptyDir() started.==========")
	isEmpty, err = IsEmptyDir(dirName)
	asst.Nil(err, "test IsEmptyDir() failed")
	asst.False(isEmpty, "test IsEmptyDir() failed")

	isEmpty, err = IsEmptyDir(empDirName)
	asst.Nil(err, "test IsEmptyDir() failed")
	asst.True(isEmpty, "test IsEmptyDir() failed")

	isEmpty, err = IsEmptyDir(fileName)
	asst.NotNil(err, "test IsEmptyDir() failed")
	t.Log("==========test IsEmptyDir() completed.==========\n")

	// test GetPathDirMapLocal()
	t.Log("==========test GetPathDirMapLocal() started.==========")
	pathDirMap, err = GetPathDirMapLocal(dirName, dirName)
	asst.Nil(err, "test GetPathDirMapLocal() failed")
	asst.Equal(len(pathDirMap), expectedElementNum, "test GetPathDirMapLocal() failed")
	t.Log("==========test GetPathDirMapLocal() completed.==========\n")

	// test GetFileNameDest()
	t.Log("==========test GetFileNameDest() started.==========")
	fileNameDest = GetFileNameDest(fileName, empDirName)
	asst.Nil(err, "test GetFileNameDest() failed")
	asst.Equal(fileNameDest, expectedFileName, "test GetFileNameDest() failed")
	t.Log("==========test GetFileNameDest() completed.==========\n")
}
