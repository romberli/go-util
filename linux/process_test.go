package linux

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/romberli/go-util/constant"
)

func killProcess(pid int, sleep time.Duration) {
	command := fmt.Sprintf("kill %d", pid)
	time.Sleep(sleep * time.Second)
	_, _ = ExecuteCommand(command)
}

func TestProcess(t *testing.T) {
	var (
		err       error
		pid       int
		pidFile   string
		isRunning bool
		sleep     time.Duration
	)

	asst := assert.New(t)

	pid = os.Getpid()
	pidFile = "go-util.pid"
	sleep = 1

	t.Log("==========SavePid started.==========")
	err = SavePid(pid, pidFile, constant.DefaultFileMode)
	asst.Nil(err, "SavePid failed.\n%v", err)
	t.Log("==========SavePid completed.==========")

	t.Log("==========IsRunningWithPid started.==========")
	isRunning = IsRunningWithPid(pid)
	asst.True(isRunning, "IsRunningWithPid failed.")
	t.Log("==========IsRunningWithPid completed.==========")

	t.Log("==========IsRunningWithPidFile started.==========")
	isRunning, err = IsRunningWithPidFile(pidFile)
	asst.Nil(err, "IsRunningWithPidFile failed.\n%v", err)
	asst.True(isRunning, "IsRunningWithPidFile failed.")
	t.Log("==========IsRunningWithPidFile completed.==========")

	t.Log("==========GetPidFromPidFile started.==========")
	pid, err = GetPidFromPidFile(pidFile)
	asst.Nil(err, "GetPidFromPidFile failed.\n%v", err)
	t.Log("==========GetPidFromPidFile completed.==========")

	t.Log("==========HandleSignalsWithPidFile started.==========")
	go killProcess(pid, sleep)
	// pidFile = "/tmp123"
	time.Sleep(sleep * time.Second)
	HandleSignalsWithPidFile(pidFile)
	asst.Nil(err, "HandleSignalsWithPidFile failed.")
	t.Log("==========HandleSignalsWithPidFile completed.==========")
}
