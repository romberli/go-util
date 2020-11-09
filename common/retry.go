package common

import (
	"errors"
	"fmt"
	"time"
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

// RetryWithRetryOption retries the func until it returns no error or reaches attempts limit or
// timed out, either one is earlier
func Retry(doFunc func() error, attempts int64, delay, timeout time.Duration) error {
	retryOption := NewRetryOption(attempts, delay, timeout)
	return RetryWithRetryOption(doFunc, retryOption)
}

// RetryWithRetryOption retries the func until it returns no error or reaches attempts limit or
// timed out, either one is earlier
func RetryWithRetryOption(doFunc func() error, opts ...RetryOption) (err error) {
	var (
		retryOption  RetryOption
		attemptCount int64
	)

	if len(opts) > 0 {
		retryOption = opts[0]
	} else {
		retryOption = NewRetryOption(DefaultAttempts, DefaultDelay, DefaultTimeout)
	}

	if retryOption.Timeout < 0 {
		return errors.New(fmt.Sprintf("timeout must NOT be less than 0, %d is not valid.", retryOption.Timeout))
	}
	if retryOption.Delay < 0 {
		return errors.New(fmt.Sprintf("delay must NOT be less than 0, %d is not valid.", retryOption.Delay))
	}

	timeoutChan := time.After(retryOption.Timeout)

	// call the function
	for attemptCount = 0; attemptCount <= retryOption.Attempts; attemptCount++ {
		err := doFunc()
		if err == nil {
			return err
		}
		// if attempts or timeout equal to 0, then not to retry
		if retryOption.Attempts == 0 || retryOption.Timeout == 0 {
			return err
		}

		// check for timeout
		select {
		case <-timeoutChan:
			return err
		default:
			time.Sleep(retryOption.Delay)
		}
	}

	return err
}
