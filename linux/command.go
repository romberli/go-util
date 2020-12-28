package linux

import (
	"bytes"
	"os/exec"
)

// ExecuteCommand executes shell command
func ExecuteCommand(command string) (output string, err error) {
	var stdoutBuffer bytes.Buffer

	cmd := exec.Command("/bin/sh", "-c", command)
	cmd.Stdout = &stdoutBuffer
	cmd.Stderr = &stdoutBuffer

	err = cmd.Run()

	return stdoutBuffer.String(), err
}
