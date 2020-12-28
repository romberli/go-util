package linux

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/romberli/go-util/constant"
)

func KillProcess(pid int, sleep time.Duration) {
	time.Sleep(sleep * time.Second)
	_, _ = ExecuteCommand(fmt.Sprintf("kill %d", pid))
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
	sleep = 2

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

	t.Log("==========GetPIDFromPidFile started.==========")
	pid, err = GetPIDFromPidFile(pidFile)
	asst.Nil(err, "GetPIDFromPidFile failed.\n%v", err)
	t.Log("==========GetPIDFromPidFile completed.==========")

	t.Log("==========HandleSignalsWithPIDFileAndLog started.==========")
	go KillProcess(pid, sleep)
	err = HandleSignalsWithPidFile(pidFile)
	asst.Nil(err, "HandleSignalsWithPIDFileAndLog failed.\n%v", err)
	t.Log("==========HandleSignalsWithPIDFileAndLog completed.==========")
}
