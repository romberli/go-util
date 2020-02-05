package log

import (
	log "github.com/sirupsen/logrus"
)

const (
	MaxLogFileSizeDefault     = 100 * 1024 * 1024 // 100MB
	RotateTimeIntervalDefault = 1                 // 1 day
)

type MyLogger struct {
	log.Logger
	MaxLogFileSize     int
	RotateTimeInterval int
}

func NewMyLogger(maxLogFileSize int, rotateTimeInterval int) (myLogger *MyLogger, err error) {
	if maxLogFileSize < 0 {
		maxLogFileSize = MaxLogFileSizeDefault
	}

	if rotateTimeInterval < 0 {
		rotateTimeInterval = RotateTimeIntervalDefault
	}

}
