package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSsh(t *testing.T) {
	var (
		err            error
		cmd            string
		hostIp         string
		result         int
		stdOut         string
		stdErr         string
		sshConn        *MySshConn
		sftpConn       *MySftpConn
		fileNameSource string
		fileNameDest   string
	)

	cmd = "date"
	hostIp = "192.168.137.11"

	assert := assert.New(t)

	// create ssh connection
	t.Log("==========create ssh connection started.==========")
	if sshConn, err = NewMySshConn(hostIp); err != nil {
		assert.Nil(err, "create new ssh connection failed")
		panic(t)
	}
	t.Log("==========create ssh connection completed.==========\n")

	// create sftp connection
	t.Log("==========create sftp connection started.==========")
	if sftpConn, err = NewMySftpConn(hostIp); err != nil {
		assert.Nil(err, "create new sftp connection failed")
	}
	t.Log("==========create sftp connection completed.==========\n")

	// test execute remote shell command
	t.Log("==========execute remote shell command started.==========")
	if result, stdOut, stdErr, err = sshConn.ExecuteCommand(cmd); err != nil {
		assert.Nil(err, "execute command failed.\ncommand: %s", cmd)
	} else {
		t.Logf("return code: %d\n\t\t\t\tstdOut: %s\n\t\t\t\tstdErr: %s", result, stdOut, stdErr)
		assert.Zero(result, "return code is NOT ZERO\nreturn code: %d", result)
		assert.NotEmpty(stdOut, 0, "command output should NOT empty\nstdOut: %s", stdOut)
		assert.Empty(stdErr, 0, "command output MUST empty\nstdErr: %s", stdErr)
	}
	t.Log("==========execute remote shell command completed.==========\n")

	// test copy single file from remote
	t.Log("==========copy single file from remote started.==========")
	fileNameSource = "/tmp/test.txt"
	fileNameDest = "/Users/romber/text.txt"
	if err = sftpConn.CopySingleFileFromRemote(fileNameSource, fileNameDest); err != nil {
		assert.Nil(err, "copy single file from remote failed")
	}
	t.Log("==========copy single file from remote completed.==========\n")
}
