package config

import (
	"fmt"

	"github.com/romberli/go-util/constant"
)

type ErrMessage struct {
	Header  string
	ErrCode int
	Raw     string
}

// NewErrMessage is an exported alias of newErrMessage() function
func NewErrMessage(header string, errCode int, raw string) *ErrMessage {
	return newErrMessage(header, errCode, raw)
}

// NewErrMessage returns a new *ErrMessage
func newErrMessage(header string, errCode int, raw string) *ErrMessage {
	return &ErrMessage{
		Header:  header,
		ErrCode: errCode,
		Raw:     raw,
	}
}

// Code returns combined Header and ErrCode string
func (e *ErrMessage) Code() string {
	return fmt.Sprintf("%s-%d", e.Header, e.ErrCode)
}

// Error is implementation fo Error interface
func (e *ErrMessage) Error() string {
	return fmt.Sprintf("%s: %s", e.Code(), e.Raw)
}

// Renew returns a new *ErrMessage and specify with given input
func (e *ErrMessage) Renew(ins ...interface{}) *ErrMessage {
	c := e.Clone()
	c.Specify(ins...)

	return c
}

// Clone returns a new *ErrMessage with same member variables
func (e *ErrMessage) Clone() *ErrMessage {
	return newErrMessage(e.Header, e.ErrCode, e.Raw)
}

// Specify specifies place holders with given data
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
