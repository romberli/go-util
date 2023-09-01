package common

import (
	"reflect"
	"runtime"
	"time"

	"github.com/pingcap/errors"
	"github.com/romberli/log"

	"github.com/romberli/go-util/constant"
)

const (
	MinMaxRetryCount = -1
	MaxMaxRetryCount = constant.MaxInt

	DefaultMaxRetryCount = 10
	DefaultDelayTime     = 500 * time.Millisecond
	DefaultMaxWaitTime   = 10 * time.Second
)

// RetryOption is options for Retry()
type RetryOption struct {
	MaxRetryCount int
	MaxWaitTime   time.Duration
	DelayTime     time.Duration
	Logger        *log.Logger
}

// NewRetryOption returns RetryOption
func NewRetryOption(maxRetryCount int, maxWaitTime, delayTime time.Duration, logger *log.Logger) *RetryOption {
	return &RetryOption{
		MaxRetryCount: maxRetryCount,
		MaxWaitTime:   maxWaitTime,
		DelayTime:     delayTime,
		Logger:        logger,
	}
}

// NewRetryOptionWithLogger returns RetryOption with given logger
func NewRetryOptionWithLogger(logger *log.Logger) *RetryOption {
	return &RetryOption{
		MaxRetryCount: DefaultMaxRetryCount,
		MaxWaitTime:   DefaultMaxWaitTime,
		DelayTime:     DefaultDelayTime,
		Logger:        logger,
	}
}

// NewRetryOptionWithDefault returns RetryOption with default values
func NewRetryOptionWithDefault() *RetryOption {
	return &RetryOption{
		MaxRetryCount: DefaultMaxRetryCount,
		MaxWaitTime:   DefaultMaxWaitTime,
		DelayTime:     DefaultDelayTime,
		Logger:        log.L(),
	}
}

// Validate validates RetryOption
func (ro *RetryOption) Validate() error {
	if ro.MaxRetryCount < MinMaxRetryCount || ro.MaxRetryCount > MaxMaxRetryCount {
		return errors.Errorf("max retry count must be between %d and %d, %d is not valid", MinMaxRetryCount, MaxMaxRetryCount, ro.MaxRetryCount)
	}

	return nil
}

// Retry retries the function until it returns no error or reaches max retry count or
// max wait time, either one is earlier, if option is nil,
// it will only call the function once, and no retry.
func Retry(doFunc func() error, option *RetryOption) error {
	if option == nil {
		return doFunc()
	}
	err := option.Validate()
	if err != nil {
		return err
	}

	timeoutChan := time.After(option.MaxWaitTime)

	var i int

	for {
		// run the function
		err = doFunc()
		if err != nil {
			if option.Logger != nil {
				funName := runtime.FuncForPC(reflect.ValueOf(doFunc).Pointer()).Name()
				option.Logger.Errorf("common.Retry(): execute function failed. function name: %s, error:\n%+v", funName, err)
			}
			// check retry count
			if option.MaxRetryCount >= constant.ZeroInt && i >= option.MaxRetryCount {
				return errors.Trace(err)
			}
			// check wait timeout
			select {
			case <-timeoutChan:
				return errors.Trace(err)
			default:
				time.Sleep(option.DelayTime)
			}

			i++
			continue
		}

		return nil
	}
}
