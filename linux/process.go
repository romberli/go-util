package linux

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/romberli/log"

	"github.com/romberli/go-util/constant"

	"github.com/shirou/gopsutil/v3/process"
)

// IsRunningWithPid returns if given pid is running
func IsRunningWithPid(pid int) (bool, error) {
	_, err := process.NewProcess(int32(pid))
	if err != nil {
		if err == process.ErrorProcessNotRunning {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// SavePid saves pid to pid file with given file mode
func SavePid(pid int, pidFile string, fileMode os.FileMode) error {
	return ioutil.WriteFile(pidFile, []byte(fmt.Sprintf("%d", pid)), fileMode)
}

// IsRunningWithPidFile returns if process of which pid was saved in given pid file is running
func IsRunningWithPidFile(pidFile string) (bool, error) {
	// check if pid file exists
	exists, err := PathExists(pidFile)
	if err != nil {
		return false, err
	}
	if !exists {
		return false, errors.New(fmt.Sprintf("pid file does not exists. pid file: %s", pidFile))
	}

	// read pid from pid file
	pidBytes, err := ioutil.ReadFile(pidFile)
	if err != nil {
		return false, err
	}
	pid, err := strconv.Atoi(string(pidBytes))
	if err != nil {
		return false, err
	}

	return IsRunningWithPid(pid)
}

// GetPidFromPidFile reads pid file and returns pid
func GetPidFromPidFile(pidFile string) (int, error) {
	var (
		err      error
		pidBytes []byte
		pidStr   string
		pid      int
	)
	pidBytes, err = ioutil.ReadFile(pidFile)
	if err != nil {
		return constant.ZeroInt, err
	}
	pidStr = strings.TrimSpace(string(pidBytes))
	pid, err = strconv.Atoi(pidStr)
	if err != nil {
		return constant.ZeroInt, err
	}

	return pid, nil
}

func RemovePidFile(pidFile string) error {
	var (
		err           error
		pidFileExists bool
	)

	pidFileExists, err = PathExists(pidFile)
	if err != nil {
		return err
	}
	if !pidFileExists {
		return errors.New(fmt.Sprintf("pid file does not exists, please have a check. pid file: %s", pidFile))
	}

	return os.Remove(pidFile)
}

// SendSignal sends signal to given pid,
// it will also remove pid file if pid file path is specified as opts
func SendSignal(pid int, sig syscall.Signal, opts ...string) (err error) {
	p, err := os.FindProcess(pid)
	if err != nil {
		return err
	}

	err = p.Signal(sig)
	if err != nil {
		return err
	}

	// remove pid file
	if len(opts) > constant.ZeroInt {
		pidFile := opts[constant.ZeroInt]
		return RemovePidFile(pidFile)
	}

	return nil
}

// KillServer kills server with given pid, it acts like shell command "kill -9",
// it will also remove pid file if pid file path is specified as opts
func KillServer(pid int, opts ...string) (err error) {
	return SendSignal(pid, syscall.SIGKILL, opts...)
}

// ShutdownServer kills server with given pid, it acts like shell command "kill -15",
// it will also remove pid file if pid file path is specified as opts
func ShutdownServer(pid int, opts ...string) (err error) {
	return SendSignal(pid, syscall.SIGTERM, opts...)
}

// HandleSignalsWithPidFile handles operating system signals
func HandleSignalsWithPidFile(pidFile string) {
	signals := make(chan os.Signal)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGHUP, syscall.SIGKILL, syscall.SIGTERM)

	for {
		sig := <-signals
		switch sig {
		case syscall.SIGINT, syscall.SIGHUP, syscall.SIGKILL, syscall.SIGTERM:
			log.Info(fmt.Sprintf("got operating system signal %d, this process will exit soon.", sig))

			err := os.Remove(pidFile)
			if err != nil {
				log.Error(fmt.Sprintf("got wrong when removing pid file. pid file: %s", pidFile))
				os.Exit(constant.DefaultAbnormalExitCode)
			}

			os.Exit(constant.DefaultNormalExitCode)
		default:
			log.Error(fmt.Sprintf("got wrong signal %d, only accept %d, %d, %d, %d",
				sig, syscall.SIGINT, syscall.SIGHUP, syscall.SIGKILL, syscall.SIGTERM))
		}
	}
}
