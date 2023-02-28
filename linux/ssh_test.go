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
		stdOut         string
		sshConn        *SSHConn
		fileNameSource string
		fileNameDest   string
	)

	asst := assert.New(t)

	cmd = "date"
	hostIP = "192.168.137.11"
	portNum = 22
	userName = "root"
	userPass = "root"

	// create ssh connection
	t.Log("==========create ssh connection started.==========")
	sshConn, err = NewSSHConn(hostIP, portNum, userName, userPass)
	if err != nil {
		asst.Nil(err, "create new ssh connection failed")
		panic(t)
	}
	t.Log("==========create ssh connection completed.==========\n")

	// test execute remote shell command
	t.Log("==========execute remote shell command started.==========")
	stdOut, err = sshConn.ExecuteCommand(cmd)
	asst.Nil(err, "execute command failed.\ncommand: %s", cmd)
	t.Logf("return code: stdOut: %s\n\t\t\t\t", stdOut)
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
