package linux

import (
	"errors"
	"fmt"

	"github.com/romberli/log"

	"github.com/romberli/go-util/common"
	"io/ioutil"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

// IsRunningWithPID returns if given pid is running
func IsRunningWithPID(pid int) bool {
	if pid > 0 {
		err := syscall.Kill(pid, 0)
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
	exists, err := common.PathExistsLocal(pidFile)
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

// SavePid saves pid to pid file with given file mode
func SavePid(pid int, pidFile string, fileMode os.FileMode) error {
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
			log.Infof("got operating system signal %d, will")
			err = os.Remove(pidFile)
			if err != nil {
				log.Errorf("remove pid file failed. pid file: %s", pidFile)
				os.Exit(2)
			}
		}
	}
}

//forkDaemon,当checkPid为true时，检查是否有存活的，有则不执行
func forkDaemon() error {
	args := os.Args
	os.Setenv("__Daemon", "true")
	procAttr := &syscall.ProcAttr{
		Env:   os.Environ(),
		Files: []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()},
	}
	pid, err := syscall.ForkExec(os.Args[0], args, procAttr)
	if err != nil {
		return err
	}
	log.Printf("[%d] %s start daemon\n", pid, appName)
	savePid(pid)
	return nil
}
