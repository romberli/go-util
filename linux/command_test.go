package linux

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommand(t *testing.T) {
	var (
		err    error
		cmd    string
		stdout string
	)

	asst := assert.New(t)

	cmd = "ls -l /tmp;ls -l /tmp/1234"

	// test command
	t.Log("==========test command started.==========")
	stdout, err = ExecuteCommand(cmd)
	t.Log(fmt.Sprintf("stdout: %s", stdout))
	asst.Nil(err, "test command failed.\ncmd: %s\n%v", cmd, err)
	t.Log("==========test command completed.==========\n")
}
