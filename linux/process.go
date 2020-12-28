package linux

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/romberli/go-util/constant"
	"github.com/romberli/log"
)

// IsRunningWithPID returns if given pid is running
func IsRunningWithPID(pid int) bool {
	if pid > 0 {
		err := syscall.Kill(pid, constant.ZeroInt)
		if err != nil {
			return false
		}

		return true
	}

	return false
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

// SavePID saves pid to pid file with given file mode
func SavePID(pid int, pidFile string, fileMode os.FileMode) error {
	return ioutil.WriteFile(pidFile, []byte(fmt.Sprintf("%d", pid)), fileMode)
}

// HandleSignals handles operating system signals
func HandleSignalsWithPIDFileAndLog(pidFile string) {
	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT)
	var err error
	for {
		sig := <-signals
		switch sig {
		case syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM:
			log.Infof("got operating system signal %d, will exit soon.")
			err = os.Remove(pidFile)
			if err != nil {
				log.Errorf("remove pid file failed. pid file: %s", pidFile)
				os.Exit(constant.DefaultAbnormalExitCode)
			}
		}
	}
}
