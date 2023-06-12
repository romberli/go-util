package common

import (
	"time"

	"github.com/pingcap/errors"

	"github.com/romberli/go-util/constant"
)

const (
	DefaultAttempts int64 = 10
	DefaultDelay          = 500 * time.Millisecond
	DefaultTimeout        = 10 * time.Second
)

// RetryOption is options for Retry()
type RetryOption struct {
	Attempts int64
	Delay    time.Duration
	Timeout  time.Duration
}

// NewRetryOption returns RetryOption
func NewRetryOption(attempts int64, delay, timeout time.Duration) RetryOption {
	return RetryOption{
		Attempts: attempts,
		Delay:    delay,
		Timeout:  timeout,
	}
}

// Retry retries the func until it returns no error or reaches attempts limit or
// timed out, either one is earlier
func Retry(doFunc func() error, attempts int64, delay, timeout time.Duration) error {
	retryOption := NewRetryOption(attempts, delay, timeout)
	return RetryWithRetryOption(doFunc, retryOption)
}

// RetryWithRetryOption retries the func until it returns no error or reaches attempts limit or
// timed out, either one is earlier
func RetryWithRetryOption(doFunc func() error, opts ...RetryOption) error {
	var (
		err          error
		retryOption  RetryOption
		attemptCount int64
	)

	if len(opts) > constant.ZeroInt {
		retryOption = opts[constant.ZeroInt]
	} else {
		retryOption = NewRetryOption(DefaultAttempts, DefaultDelay, DefaultTimeout)
	}

	if retryOption.Timeout < constant.ZeroInt {
		return errors.Errorf("timeout must NOT be less than 0, %d is not valid.", retryOption.Timeout)
	}
	if retryOption.Delay < constant.ZeroInt {
		return errors.Errorf("delay must NOT be less than 0, %d is not valid.", retryOption.Delay)
	}

	timeoutChan := time.After(retryOption.Timeout)

	// call the function
	for attemptCount = constant.ZeroInt; attemptCount <= retryOption.Attempts; attemptCount++ {
		err = doFunc()
		if err == nil {
			return nil
		}
		// if attempts or timeout equal to 0, then not to retry
		if retryOption.Attempts == constant.ZeroInt || retryOption.Timeout == constant.ZeroInt {
			return errors.Trace(err)
		}

		// check for timeout
		select {
		case <-timeoutChan:
			return errors.Trace(err)
		default:
			time.Sleep(retryOption.Delay)
		}
	}

	return errors.Trace(err)
}
