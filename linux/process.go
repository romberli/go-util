package linux

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/romberli/log"

	"github.com/romberli/go-util/constant"
)

// IsRunningWithPID returns if given pid is running
func IsRunningWithPID(pid int) bool {
	if pid > 0 {
		err := syscall.Kill(pid, syscall.Signal(constant.ZeroInt))
		if err != nil {
			return false
		}

		return true
	}

	return false
}

// SavePID saves pid to pid file with given file mode
func SavePID(pid int, pidFile string, fileMode os.FileMode) error {
	return ioutil.WriteFile(pidFile, []byte(fmt.Sprintf("%d", pid)), fileMode)
}

// IsRunningWithPIDFile returns if process of which pid was saved in given pid file is running
func IsRunningWithPIDFile(pidFile string) (bool, error) {
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

	return IsRunningWithPID(pid), nil
}

// GetPID reads pid file and returns pid
func GetPIDFromPIDFile(pidFile string) (int, error) {
	pidBytes, err := ioutil.ReadFile(pidFile)
	if err != nil {
		return constant.ZeroInt, err
	}
	pid, err := strconv.Atoi(string(pidBytes))
	if err != nil {
		return constant.ZeroInt, err
	}

	return pid, nil
}

// HandleSignals handles operating system signals
func HandleSignalsWithPIDFile(pidFile string) error {
	var err error

	signals := make(chan os.Signal)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGHUP, syscall.SIGKILL, syscall.SIGTERM)

	for {
		sig := <-signals
		switch sig {
		case syscall.SIGINT, syscall.SIGHUP, syscall.SIGKILL, syscall.SIGTERM:
			log.Info(fmt.Sprintf("got operating system signal %d, this process will exit soon.", sig))

			err = os.Remove(pidFile)
			if err != nil {
				return err
			}

			os.Exit(constant.DefaultNormalExitCode)
		default:
			return errors.New(fmt.Sprintf("got wrong signal %d, only accept %d, %d, %d, %d",
				sig, syscall.SIGINT, syscall.SIGHUP, syscall.SIGKILL, syscall.SIGTERM))
		}
	}
}
