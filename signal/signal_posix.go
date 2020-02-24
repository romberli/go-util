package signal

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	log "github.com/romber2001/go-util/mylog"
	"go.uber.org/zap"
)

// SetupSignalHandler setup signal handler for TiDB Server
func SetupSignalHandler(shutdownFunc func(bool)) {
	usrDefSignalChan := make(chan os.Signal, 1)

	signal.Notify(usrDefSignalChan, syscall.SIGUSR1)
	go func() {
		buf := make([]byte, 1<<16)
		for {
			sig := <-usrDefSignalChan
			if sig == syscall.SIGUSR1 {
				stackLen := runtime.Stack(buf, true)
				log.Info(
					fmt.Sprintf("\n=== Got signal [%s] to dump goroutine stack. ===\n"+
						"%s\n=== Finished dumping goroutine stack. ===\n", sig, buf[:stackLen]))
			}
		}
	}()

	closeSignalChan := make(chan os.Signal, 1)
	signal.Notify(closeSignalChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		sig := <-closeSignalChan
		log.Info("got signal to exit", zap.Stringer("signal", sig))
		shutdownFunc(sig == syscall.SIGQUIT)
	}()
}
