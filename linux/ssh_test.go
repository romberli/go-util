package linux

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSSH(t *testing.T) {
	var (
		err            error
		cmd            string
		hostIP         string
		portNum        int
		userName       string
		userPass       string
		result         int
		stdOut         string
		sshConn        *MySSHConn
		fileNameSource string
		fileNameDest   string
	)

	cmd = "date"
	hostIP = "192.168.137.11"
	portNum = 22
	userName = "root"
	userPass = "shit"

	asst := assert.New(t)

	// create ssh connection
	t.Log("==========create ssh connection started.==========")
	sshConn, err = NewMySSHConn(hostIP, portNum, userName, userPass)
	if err != nil {
		asst.Nil(err, "create new ssh connection failed")
		panic(t)
	}
	t.Log("==========create ssh connection completed.==========\n")

	// test execute remote shell command
	t.Log("==========execute remote shell command started.==========")
	result, stdOut, err = sshConn.ExecuteCommand(cmd)
	asst.Nil(err, "execute command failed.\ncommand: %s", cmd)
	t.Logf("return code: %d\n\t\t\t\tstdOut: %s\n\t\t\t\t", result, stdOut)
	asst.Zero(result, "return code is NOT ZERO\nreturn code: %d", result)
	asst.NotEmpty(stdOut, 0, "command output should NOT empty\nstdOut: %s", stdOut)
	t.Log("==========execute remote shell command completed.==========\n")

	// test copy single file from remote
	t.Log("==========copy single file from remote started.==========")
	fileNameSource = "/root/common.go"
	fileNameDest = "/Users/romber/common.go"
	err = sshConn.CopyFromRemote(fileNameSource, fileNameDest)
	asst.Nil(err, "copy single file from remote failed")
	t.Log("==========copy single file from remote completed.==========\n")

	// test copy single file to remote
	t.Log("==========copy single file from remote started.==========")
	fileNameSource = "/Users/romber/common2.go"
	fileNameDest = "/root"
	err = sshConn.CopyToRemote(fileNameSource, fileNameDest)
	asst.Nil(err, "copy single file from remote failed")
	t.Log("==========copy single file from remote completed.==========\n")

	// test copy a directory from remote
	t.Log("==========copy a directory from remote started.==========")
	fileNameSource = "/tmp/test"
	fileNameDest = "/Users/romber/test2"
	// err = os.RemoveAll(fileNameDest)
	asst.Nil(err, "remove directory on local host failed")
	err = sshConn.CopyFromRemote(fileNameSource, fileNameDest)
	asst.Nil(err, "copy directory from remote failed")
	t.Log("==========copy a directory from remote completed.==========\n")

	// test copy a directory to remote
	t.Log("==========copy a directory to remote started.==========")
	fileNameSource = "/Users/romber/test"
	fileNameDest = "/tmp/test2"
	err = sshConn.RemoveAll(fileNameDest)
	asst.Nil(err, "remove directory on remote host failed")
	err = sshConn.CopyToRemote(fileNameSource, fileNameDest)
	asst.Nil(err, "copy directory to remote failed")
	fileNameSource = "/Users/romber/test/1.txt"
	fileNameDest = "/tmp/test2"
	t.Log("==========copy a directory to remote completed.==========\n")

}
