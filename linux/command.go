package linux

import (
	"bytes"
	"os/exec"
	"strings"

	"github.com/pingcap/errors"

	"github.com/romberli/go-util/common"
	"github.com/romberli/go-util/constant"
)

const (
	shCommand   = "sh"
	shPath      = "/bin/sh"
	bashCommand = "/bin/bash"
	dashCArg    = "-c"
)

type CommandOption interface {
	apply(*exec.Cmd)
}

type optionFunc func(*exec.Cmd)

func (f optionFunc) apply(cmd *exec.Cmd) {
	f(cmd)
}

func WorkDirOption(workDir string) CommandOption {
	return optionFunc(func(cmd *exec.Cmd) {
		cmd.Dir = workDir
	})
}

func UseSHCOption() CommandOption {
	return optionFunc(func(cmd *exec.Cmd) {
		cmd.Path = shPath
		cmd.Args = []string{shCommand, dashCArg, common.ConvertSliceToString(cmd.Args, constant.SpaceString)}
	})
}

// ExecuteCommand is an alias of ExecuteCommandAndWait
// 1. first arg must be the working directory
func ExecuteCommand(command string, options ...CommandOption) (output string, err error) {
	return ExecuteCommandAndWait(command, options...)
}

// ExecuteCommandAndWait executes shell command and wait for it to complete
func ExecuteCommandAndWait(command string, options ...CommandOption) (output string, err error) {
	cmd, err := getCommand(command, options...)
	if err != nil {
		return constant.EmptyString, err
	}

	err = cmd.Run()

	return cmd.Stdout.(*bytes.Buffer).String(), errors.Trace(err)
}

// ExecuteCommandNoWait executes shell command and does not wait for it to complete
func ExecuteCommandNoWait(command string, options ...CommandOption) (output string, err error) {
	cmd, err := getCommand(command, options...)
	if err != nil {
		return constant.EmptyString, err
	}

	err = cmd.Start()

	return cmd.Stdout.(*bytes.Buffer).String(), errors.Trace(err)
}

func getCommand(command string, options ...CommandOption) (*exec.Cmd, error) {
	var (
		stdoutBuffer bytes.Buffer
	)

	commandList := strings.Split(command, constant.SpaceString)
	if len(commandList) == constant.ZeroInt {
		return nil, errors.New("command is empty")
	}

	cmdName := commandList[constant.ZeroInt]
	if len(commandList) > constant.TwoInt && (cmdName == shCommand || cmdName == bashCommand) && commandList[constant.OneInt] == dashCArg {
		commandList = strings.SplitN(command, constant.SpaceString, constant.ThreeInt)
	}

	cmd := exec.Command(commandList[constant.ZeroInt], commandList[constant.OneInt:]...)
	cmd.Stdout = &stdoutBuffer
	cmd.Stderr = &stdoutBuffer

	for _, option := range options {
		option.apply(cmd)
	}

	return cmd, nil
}
