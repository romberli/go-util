package mylog

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newRoutine(t *testing.T, wg *sync.WaitGroup) {
	defer wg.Done()
	t.Log("new routine test")
	Debug("this is new routine debug message")
	Info("this is new routine info message")
	Warn("this is new routine warn message")
}

func TestLog(t *testing.T) {
	var (
		err           error
		fileLogConfig *FileLogConfig
		logConfig     *LogConfig
	)

	wg := &sync.WaitGroup{}
	assert := assert.New(t)

	level := "info"
	format := "text"
	fileName := "/Users/romber/run.log"
	maxSize := 1
	maxDays := 1
	maxBackups := 2

	// init logger
	t.Log("==========init logger started==========")
	fileLogConfig, err = NewFileLogConfig(fileName, maxSize, maxDays, maxBackups)
	assert.Nil(err, "init file log config failed")

	logConfig = NewLogConfig(level, format, *fileLogConfig, false)
	err = InitZapLogger(logConfig)
	assert.Nil(err, "init logger failed")
	t.Log("==========init logger completed==========\n")

	// print log
	t.Log("==========print main log entry started==========")
	Debug("this is main debug message")
	Info("this is main info message")
	MyLogger.Warn("this is main warn message")
	// MyLogger.Error("this is main error message")
	// MyLogger.Fatal("this is main fatal message")
	t.Log("==========print main log entry completed==========")

	t.Log("==========print goroutine log entry started==========")

	wg.Add(1)
	go func() {
		defer wg.Done()
		t.Log("goroutine test")
		MyLogger.Debug("this is goroutine debug message")
		MyLogger.Info("this is goroutine info message")
		MyLogger.Warn("this is goroutine warn message")
	}()

	wg.Add(1)
	go newRoutine(t, wg)
	t.Log("==========print goroutine log entry completed==========")

	wg.Wait()
}
