package config

import (
	"fmt"

	"github.com/pingcap/errors"
	"github.com/romberli/go-util/constant"
)

const defaultCallerSkip = 1

type ErrMessage struct {
	Header  string
	ErrCode int
	Raw     string
	errors.StackTrace
}

// NewErrMessage returns a *ErrMessage without stack trace
func NewErrMessage(header string, errCode int, raw string) *ErrMessage {
	return newErrMessage(header, errCode, raw, nil)
}

// NewErrMessageWithStack returns a *ErrMessage with stack trace
func NewErrMessageWithStack(header string, errCode int, raw string) *ErrMessage {
	return newErrMessage(header, errCode, raw, errors.NewStack(defaultCallerSkip).StackTrace())
}

// NewErrMessage returns a new *ErrMessage
func newErrMessage(header string, errCode int, raw string, stackTrace errors.StackTrace) *ErrMessage {
	return &ErrMessage{
		Header:     header,
		ErrCode:    errCode,
		Raw:        raw,
		StackTrace: stackTrace,
	}
}

// Code returns combined Header and ErrCode string
func (e *ErrMessage) Code() string {
	return fmt.Sprintf("%s-%d", e.Header, e.ErrCode)
}

// Error is an implementation fo Error interface
func (e *ErrMessage) Error() string {
	return fmt.Sprintf("%s: %s", e.Code(), e.Raw)
}

// String is an alias of Error()
func (e *ErrMessage) String() string {
	return e.Error()
}

// Renew returns a new *ErrMessage and specify with given input
func (e *ErrMessage) Renew(ins ...interface{}) *ErrMessage {
	c := e.Clone()
	c.Specify(ins...)

	return c
}

// Clone returns a new *ErrMessage with same member variables
func (e *ErrMessage) Clone() *ErrMessage {
	return newErrMessage(e.Header, e.ErrCode, e.Raw, e.StackTrace)
}

// Specify specifies placeholders with given data
func (e *ErrMessage) Specify(ins ...interface{}) {
	e.Raw = fmt.Sprintf(e.Raw, ins...)
}

// ErrorOrNil returns an error interface if both Header and ErrCode are not zero value, otherwise, returns nil.
// This function is useful at the end of accumulation to make sure that the value
// returned represents the existence of errors
func (e *ErrMessage) ErrorOrNil() error {
	if e == nil {
		return nil
	}

	if e.Header == constant.EmptyString || e.ErrCode == constant.ZeroInt {
		return nil
	}

	return e
}
