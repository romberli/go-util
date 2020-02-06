package mylog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func newRoutine(t *testing.T) {
	t.Log("goroutine test")
	logger.Debug("this is goroutine debug message")
	logger.Info("this is goroutine info message")
	logger.Warn("this is goroutine warn message")
}

func TestLog(t *testing.T) {
	var (
		err           error
		fileLogConfig *FileLogConfig
		logConfig     *LogConfig
	)

	level := "info"
	format := "text"
	fileName := "/Users/romber/run.log"
	maxSize := 1
	maxDays := 1
	maxBackups := 2

	assert := assert.New(t)

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
	logger.Debug("this is main debug message")
	logger.Info("this is main info message")
	logger.Warn("this is main warn message")
	// logger.Error("this is main error message")
	// logger.Fatal("this is main fatal message")
	t.Log("==========print main log entry completed==========")

	t.Log("==========print goroutine log entry started==========")
	go func() {
		t.Log("goroutine test")
		logger.Debug("this is goroutine debug message")
		logger.Info("this is goroutine info message")
		logger.Warn("this is goroutine warn message")
	}()

	go newRoutine(t)
	t.Log("==========print goroutine log entry completed==========")
}
