package linux

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/romberli/go-util/constant"
)

const (
	testHostIP   = "192.168.137.11"
	testHostName = "192-168-137-11"
	testPortNum  = 22
	testUserName = "romber"
	testUserPass = "romber"
	testUseSudo  = true

	testRootHomeDir    = "/root"
	testLocalPath      = "/Users/romber/test_local"
	testLocalFileName  = "test_local.txt"
	testRemotePath     = "/data"
	testRemoteSubPath  = "/data/test_remote"
	testRemoteFileName = "test_remote.txt"
	testTmpDir         = "/tmp"
	testContent        = "test_content"

	testDateCommand             = "date"
	echoCommand                 = `sh -c 'echo "%s" > %s'`
	echoAppendCommand           = `sh -c 'echo "%s" >> %s'`
	testCreateLocalDirCommand   = "mkdir -p " + testLocalPath
	testRemoveLocalDirCommand   = "rm -rf " + testLocalPath
	testCreateRemoteFileCommand = "touch /data/test_remote/test.txt"
)

var testSSHConn *SSHConn

func init() {
	// testInitSSHConn()
}

func testInitSSHConn() {
	var err error

	testSSHConn, err = NewSSHConn(testHostIP, testPortNum, testUserName, testUserPass, testUseSudo)
	if err != nil {
		panic(err)
	}
}

func TestSSHConn_All(t *testing.T) {
	TestSSHConn_ExecuteCommand(t)
	TestSSHConn_GetHostName(t)
	TestSSHConn_PathExists(t)
	TestSSHConn_IsDir(t)
	TestSSHConn_ListPath(t)
	TestSSHConn_ReadDir(t)
	TestSSHConn_MkdirAll(t)
	TestSSHConn_RemoveAll(t)
	TestSSHConn_IsEmptyDir(t)
	TestSSHConn_CopyFile(t)
	TestSSHConn_CopyFromRemote(t)
	TestSSHConn_CopyToRemote(t)
}

func TestSSHConn_ExecuteCommand(t *testing.T) {
	asst := assert.New(t)

	stdOut, err := testSSHConn.ExecuteCommand(testDateCommand)
	asst.Nil(err, "test ExecuteCommand() failed")
	asst.NotEmpty(stdOut, "test ExecuteCommand() failed")
	t.Logf("output: %s", stdOut)
}

func TestSSHConn_GetHostName(t *testing.T) {
	asst := assert.New(t)

	hostName, err := testSSHConn.GetHostName()
	asst.Nil(err, "test GetHostName() failed")
	asst.Equal(testHostName, hostName, "test GetHostName() failed")
}

func TestSSHConn_PathExists(t *testing.T) {
	asst := assert.New(t)

	exists, err := testSSHConn.PathExists(testRemotePath)
	asst.Nil(err, "test PathExists() failed")
	asst.True(exists, "test PathExists() failed")
}

func TestSSHConn_IsDir(t *testing.T) {
	asst := assert.New(t)

	isDir, err := testSSHConn.IsDir(testRemotePath)
	asst.Nil(err, "test IsDir() failed")
	asst.True(isDir, "test IsDir() failed")
}

func TestSSHConn_ListPath(t *testing.T) {
	asst := assert.New(t)

	output, err := testSSHConn.ListPath(testRootHomeDir)
	if testUseSudo {
		asst.Nil(err, "test ListPath() failed")
		asst.NotEmpty(output, "test ListPath() failed")
		t.Log("sub directories or files: ", output)
	} else {
		asst.NotNil(err, "test ListPath() failed")
		asst.Empty(output, "test ListPath() failed")
	}
}

func TestSSHConn_ReadDir(t *testing.T) {
	asst := assert.New(t)

	fileInfos, err := testSSHConn.ReadDir(testRemotePath)
	asst.Nil(err, "test ReadDir() failed")
	asst.NotNil(fileInfos, "test ReadDir() failed")
	for _, fileInfo := range fileInfos {
		t.Log("file name: ", fileInfo.Name())
	}
}

func TestSSHConn_MkdirAll(t *testing.T) {
	asst := assert.New(t)

	err := testSSHConn.MkdirAll(testRemoteSubPath)
	asst.Nil(err, "test MkdirAll() failed")
	exists, err := testSSHConn.PathExists(testRemoteSubPath)
	asst.Nil(err, "test MkdirAll() failed")
	asst.True(exists, "test MkdirAll() failed")
}

func TestSSHConn_RemoveAll(t *testing.T) {
	asst := assert.New(t)

	err := testSSHConn.MkdirAll(testRemoteSubPath)
	asst.Nil(err, "test RemoveAll() failed")
	output, err := testSSHConn.ExecuteCommand(testCreateRemoteFileCommand)
	asst.Nil(err, "test RemoveAll() failed")
	asst.Empty(output, "test RemoveAll() failed")
	err = testSSHConn.RemoveAll(testRemoteSubPath)
	asst.Nil(err, "test RemoveAll() failed")
	exists, err := testSSHConn.PathExists(testRemoteSubPath)
	asst.Nil(err, "test RemoveAll() failed")
	asst.False(exists, "test RemoveAll() failed")
}

func TestSSHConn_Cat(t *testing.T) {
	asst := assert.New(t)

	err := testSSHConn.MkdirAll(testRemoteSubPath)
	asst.Nil(err, "test Cat() failed")
	output, err := testSSHConn.ExecuteCommand(testCreateRemoteFileCommand)
	asst.Nil(err, "test Cat() failed")
	asst.Empty(output, "test Cat() failed")
	remoteFilePath := filepath.Join(testRemoteSubPath, testRemoteFileName)
	err = testSSHConn.ExecuteCommandWithoutOutput(fmt.Sprintf(echoCommand, testContent, remoteFilePath))
	asst.Nil(err, "test Cat() failed")
	output, err = testSSHConn.Cat(remoteFilePath)
	asst.Nil(err, "test Cat() failed")
	asst.Equal(testContent, output, "test Cat() failed")
	t.Log("cat output: ", output)
	err = testSSHConn.RemoveAll(testRemoteSubPath)
	asst.Nil(err, "test Cat() failed")
}

func TestSSHConn_IsEmptyDir(t *testing.T) {
	asst := assert.New(t)

	err := testSSHConn.MkdirAll(testRemoteSubPath)
	asst.Nil(err, "test RemoveAll() failed")
	isEmpty, err := testSSHConn.IsEmptyDir(testRemoteSubPath)
	asst.Nil(err, "test IsEmptyDir() failed")
	asst.True(isEmpty, "test IsEmptyDir() failed")
	output, err := testSSHConn.ExecuteCommand(testCreateRemoteFileCommand)
	asst.Nil(err, "test RemoveAll() failed")
	asst.Empty(output, "test RemoveAll() failed")
	isEmpty, err = testSSHConn.IsEmptyDir(testRemoteSubPath)
	asst.Nil(err, "test IsEmptyDir() failed")
	asst.False(isEmpty, "test IsEmptyDir() failed")
	err = testSSHConn.RemoveAll(testRemoteSubPath)
	asst.Nil(err, "test RemoveAll() failed")
	exists, err := testSSHConn.PathExists(testRemoteSubPath)
	asst.Nil(err, "test RemoveAll() failed")
	asst.False(exists, "test RemoveAll() failed")
}

func TestSSHConn_CopyFile(t *testing.T) {
	asst := assert.New(t)

	output, err := ExecuteCommand(testCreateLocalDirCommand)
	asst.Nil(err, "test CopyFile() failed")
	asst.Empty(output, "test CopyFile() failed")
	fileNameSource := filepath.Join(testLocalPath, testLocalFileName)
	_, err = os.Create(fileNameSource)
	asst.Nil(err, "test CopyFile() failed")
	fileSource, err := os.Open(fileNameSource)
	asst.Nil(err, "test CopyFile() failed")
	defer func() { _ = fileSource.Close() }()

	fileNameDest := filepath.Join(constant.DefaultTmpDir, testRemoteFileName)
	fileDest, err := testSSHConn.SFTPClient.Create(fileNameDest)
	asst.Nil(err, "test CopyFile() failed")
	defer func() { _ = fileDest.Close() }()

	err = testSSHConn.CopyFile(fileSource, fileDest, DefaultByteBufferSize)
	asst.Nil(err, "test CopyFile() failed")
	exists, err := testSSHConn.PathExists(fileNameDest)
	asst.Nil(err, "test CopyFile() failed")
	asst.True(exists, "test CopyFile() failed")

	output, err = ExecuteCommand(testRemoveLocalDirCommand)
	asst.Nil(err, "test CopyFile() failed")
	asst.Empty(output, "test CopyFile() failed")
	err = testSSHConn.RemoveAll(fileNameDest)
	asst.Nil(err, "test CopyFile() failed")
	exists, err = testSSHConn.PathExists(fileNameDest)
	asst.Nil(err, "test CopyFile() failed")
	asst.False(exists, "test CopyFile() failed")
}

func TestSSHConn_CopyFromRemote(t *testing.T) {
	asst := assert.New(t)

	// prepare remote
	fileNameSource := filepath.Join(testRemotePath, testRemoteFileName)
	err := testSSHConn.Touch(fileNameSource)
	asst.Nil(err, "test CopyFromRemote() failed")
	pathExists, err := testSSHConn.PathExists(fileNameSource)
	asst.Nil(err, "test CopyFromRemote() failed")
	asst.True(pathExists, "test CopyFromRemote() failed")
	err = testSSHConn.Chmod(fileNameSource, "0600")
	asst.Nil(err, "test CopyFromRemote() failed")
	// prepare local
	output, err := ExecuteCommand(testCreateLocalDirCommand)
	asst.Nil(err, "test CopyFromRemote() failed")
	asst.Empty(output, "test CopyFromRemote() failed")
	fileNameDest := filepath.Join(testLocalPath, testLocalFileName)
	// copy from remote
	err = testSSHConn.CopyFromRemote(fileNameSource, fileNameDest)
	asst.NotNil(err, "test CopyFromRemote() failed")
	err = testSSHConn.CopyFromRemote(fileNameSource, fileNameDest, testTmpDir)
	asst.Nil(err, "test CopyFromRemote() failed")
	exists, err := PathExists(fileNameDest)
	asst.Nil(err, "test CopyFromRemote() failed")
	asst.True(exists, "test CopyFromRemote() failed")
	// clean up
	err = os.Remove(fileNameDest)
	asst.Nil(err, "test CopyFromRemote() failed")
	exists, err = PathExists(fileNameDest)
	asst.Nil(err, "test CopyFromRemote() failed")
	asst.False(exists, "test CopyFromRemote() failed")
	err = testSSHConn.RemoveAll(fileNameSource)
	asst.Nil(err, "test CopyFromRemote() failed")
	exists, err = testSSHConn.PathExists(fileNameSource)
	asst.Nil(err, "test CopyFromRemote() failed")
	asst.False(exists, "test CopyFromRemote() failed")
}

func TestSSHConn_CopyToRemote(t *testing.T) {
	asst := assert.New(t)

	// prepare local
	output, err := ExecuteCommand(testCreateLocalDirCommand)
	asst.Nil(err, "test CopyToRemote() failed")
	asst.Empty(output, "test CopyToRemote() failed")
	fileNameSource := filepath.Join(testLocalPath, testLocalFileName)
	_, err = os.Create(fileNameSource)
	asst.Nil(err, "test CopyToRemote() failed")
	// prepare remote
	fileNameDest := filepath.Join(testRemotePath, testRemoteFileName)
	// copy to remote
	err = testSSHConn.CopyToRemote(fileNameSource, fileNameDest, testTmpDir)
	asst.Nil(err, "test CopyToRemote() failed")
	exists, err := testSSHConn.PathExists(fileNameDest)
	asst.Nil(err, "test CopyToRemote() failed")
	asst.True(exists, "test CopyToRemote() failed")
	// clean up
	err = os.Remove(fileNameSource)
	asst.Nil(err, "test CopyToRemote() failed")
	exists, err = PathExists(fileNameSource)
	asst.Nil(err, "test CopyToRemote() failed")
	asst.False(exists, "test CopyToRemote() failed")
	err = testSSHConn.RemoveAll(fileNameDest)
	asst.Nil(err, "test CopyToRemote() failed")
	exists, err = testSSHConn.PathExists(fileNameDest)
	asst.Nil(err, "test CopyToRemote() failed")
	asst.False(exists, "test CopyToRemote() failed")
}

func TestSSHConn_Common(t *testing.T) {
	asst := assert.New(t)

	testSSHConn.SetUseSudo(false)
	cmd := `sudo su - && whoami`
	output, err := testSSHConn.ExecuteCommand(cmd)
	asst.Nil(err, "test ExecuteCommand() failed")
	output, err = testSSHConn.ExecuteCommand("whoami")
	t.Log(output)
}
