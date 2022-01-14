package linux

import (
	"bytes"
	"os/exec"

	"github.com/pingcap/errors"
)

// ExecuteCommand is an alias of ExecuteCommandAndWait
func ExecuteCommand(command string) (output string, err error) {
	return ExecuteCommandAndWait(command)
}

// ExecuteCommandAndWait executes shell command and wait for it to complete
func ExecuteCommandAndWait(command string) (output string, err error) {
	var stdoutBuffer bytes.Buffer

	cmd := exec.Command("sh", "-c", command)
	cmd.Stdout = &stdoutBuffer
	cmd.Stderr = &stdoutBuffer

	err = cmd.Run()

	return stdoutBuffer.String(), errors.Trace(err)
}

// ExecuteCommandNoWait executes shell command and does not wait for it to complete
func c(command string) (output string, err error) {
	var stdoutBuffer bytes.Buffer

	cmd := exec.Command("sh", "-c", command)
	cmd.Stdout = &stdoutBuffer
	cmd.Stderr = &stdoutBuffer

	err = cmd.Start()

	return stdoutBuffer.String(), errors.Trace(err)
}
