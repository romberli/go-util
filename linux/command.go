package linux

import (
	"bytes"
	"os/exec"
	"strings"

	"github.com/pingcap/errors"

	"github.com/romberli/go-util/constant"
)

// ExecuteCommand is an alias of ExecuteCommandAndWait
func ExecuteCommand(command string) (output string, err error) {
	return ExecuteCommandAndWait(command)
}

// ExecuteCommandAndWait executes shell command and wait for it to complete
func ExecuteCommandAndWait(command string) (output string, err error) {
	var stdoutBuffer bytes.Buffer

	c := strings.Split(command, constant.SpaceString)
	cmd := exec.Command(c[constant.ZeroInt], c[constant.OneInt:]...)
	cmd.Stdout = &stdoutBuffer
	cmd.Stderr = &stdoutBuffer

	err = cmd.Run()

	return stdoutBuffer.String(), errors.Trace(err)
}

// ExecuteCommandNoWait executes shell command and does not wait for it to complete
func ExecuteCommandNoWait(command string) (output string, err error) {
	var stdoutBuffer bytes.Buffer

	cmd := exec.Command("sh", "-c", command)
	cmd.Stdout = &stdoutBuffer
	cmd.Stderr = &stdoutBuffer

	err = cmd.Start()

	return stdoutBuffer.String(), errors.Trace(err)
}
