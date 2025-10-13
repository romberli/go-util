package linux

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/romberli/log"
	"github.com/shirou/gopsutil/v4/process"
	"github.com/stretchr/testify/assert"

	"github.com/romberli/go-util/constant"
)

func StartSandbox(cmd string) {
	_, err := ExecuteCommand(cmd)
	if err != nil {
		fmt.Println(fmt.Sprintf("error: %s", err.Error()))
	}

}

func killProcess(pid int, sleep int) {
	p, err := process.NewProcess(int32(pid))
	if err != nil {
		fmt.Println(fmt.Sprintf("error: %s", err.Error()))
	} else {
		_ = p.Kill()
	}

	time.Sleep(time.Duration(sleep) * time.Second)

	err = p.Kill()
	if err != nil {
		log.Errorf("kill process failed: %s", err.Error())
	}
}

func TestProcess(t *testing.T) {
	var (
		err       error
		pid       int
		pidFile   string
		isRunning bool
		sleep     int
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
	isRunning, err = IsRunningWithPid(pid)
	asst.Nil(err, "IsRunningWithPid failed.")
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

	t.Log("==========KillServer started.==========")
	binPath := "/Users/romber/source_code/go/src/github.com/romberli/go-util/linux/test/bin/go-sandbox"
	count := 100
	pidFileSandbox := "/tmp/go-sandbox.pid"
	cmd := fmt.Sprintf("%s --count=%d --pid-file=%s", binPath, count, pidFileSandbox)
	go StartSandbox(cmd)
	time.Sleep(time.Duration(sleep) * time.Second)
	asst.Nil(err, "start go-sandbox failed.")
	pidSandbox, err := GetPidFromPidFile(pidFileSandbox)
	asst.Nil(err, "get pid of go-sandbox failed.")
	err = KillServer(pidSandbox, pidFileSandbox)
	asst.Nil(err, "KillServer failed.\n%v", err)
	t.Log("==========KillServer completed.==========")

	t.Log("==========ShutdownServer started.==========")
	go StartSandbox(cmd)
	time.Sleep(time.Duration(sleep) * time.Second)
	pidSandbox, err = GetPidFromPidFile(pidFileSandbox)
	asst.Nil(err, "get pid of go-sandbox failed.")
	err = ShutdownServer(pidSandbox, pidFileSandbox)
	asst.Nil(err, "ShutdownServer failed.\n%v", err)
	t.Log("==========ShutdownServer completed.==========")

	t.Log("==========HandleSignals started.==========")
	go killProcess(pid, sleep)
	time.Sleep(time.Duration(sleep) * time.Second)
	HandleSignals(pidFile)
	asst.Nil(err, "HandleSignals failed.")
	t.Log("==========HandleSignalsWithPidFile completed.==========")
}

func TestHandleSignals(t *testing.T) {
	pid := os.Getpid()
	pidFile := "go-util.pid"
	sleep := 10

	go killProcess(pid, sleep)
	HandleSignals(pidFile)
}
