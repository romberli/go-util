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

	t.Log("==========SavePID started.==========")
	err = SavePID(pid, pidFile, constant.DefaultFileMode)
	asst.Nil(err, "SavePID failed.\n%v", err)
	t.Log("==========SavePID completed.==========")

	t.Log("==========IsRunningWithPID started.==========")
	isRunning = IsRunningWithPID(pid)
	asst.True(isRunning, "IsRunningWithPID failed.")
	t.Log("==========IsRunningWithPID completed.==========")

	t.Log("==========IsRunningWithPIDFile started.==========")
	isRunning, err = IsRunningWithPIDFile(pidFile)
	asst.Nil(err, "IsRunningWithPIDFile failed.\n%v", err)
	asst.True(isRunning, "IsRunningWithPIDFile failed.")
	t.Log("==========IsRunningWithPIDFile completed.==========")

	t.Log("==========GetPIDFromPIDFile started.==========")
	pid, err = GetPIDFromPIDFile(pidFile)
	asst.Nil(err, "GetPIDFromPIDFile failed.\n%v", err)
	t.Log("==========GetPIDFromPIDFile completed.==========")

	t.Log("==========HandleSignalsWithPIDFileAndLog started.==========")
	go KillProcess(pid, sleep)
	err = HandleSignalsWithPIDFile(pidFile)
	asst.Nil(err, "HandleSignalsWithPIDFileAndLog failed.\n%v", err)
	t.Log("==========HandleSignalsWithPIDFileAndLog completed.==========")
}
