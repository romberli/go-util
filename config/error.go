package config

import (
	"fmt"
	"io"

	"github.com/pingcap/errors"
	"github.com/romberli/go-multierror"
	"github.com/romberli/go-util/constant"
)

const defaultCallerSkip = 1

type ErrMessage struct {
	Header  string
	ErrCode int
	Raw     string
	Err     error
	Stack   errors.StackTrace
}

// NewErrMessage returns a *ErrMessage without stack trace,
// note that if errs contains more than one error, only the first one will be used,
// if you want to use more than one error, please use github.com/romberli/go-multierror to wrap those errors to one
func NewErrMessage(header string, errCode int, raw string, errs ...error) *ErrMessage {
	var err error

	if len(errs) > constant.ZeroInt {
		err = errs[constant.ZeroInt]
	}

	if err != nil && errors.HasStack(err) {
		return newErrMessage(header, errCode, raw, err, errors.GetStackTracer(err).StackTrace())
	}

	return newErrMessage(header, errCode, raw, err, errors.NewStack(defaultCallerSkip).StackTrace())
}

// NewErrMessage returns a new *ErrMessage
func newErrMessage(header string, errCode int, raw string, err error, stack errors.StackTrace) *ErrMessage {
	return &ErrMessage{
		Header:  header,
		ErrCode: errCode,
		Raw:     raw,
		Err:     err,
		Stack:   stack,
	}
}

func (e *ErrMessage) StackTrace() errors.StackTrace {
	merr, ok := e.Err.(*multierror.Error)
	if ok {
		if merr != nil && merr.Len() > constant.ZeroInt {
			return errors.GetStackTracer(merr.WrappedErrors()[constant.ZeroInt]).StackTrace()
		}
	}

	return e.Stack
}

// Code returns combined Header and ErrCode string
func (e *ErrMessage) Code() string {
	return fmt.Sprintf("%s-%d", e.Header, e.ErrCode)
}

// Error is an implementation fo Error interface
func (e *ErrMessage) Error() string {
	message := fmt.Sprintf("%s: %s", e.Code(), e.Raw)

	if e.Err != nil {
		message += fmt.Sprintf(". error: %s", e.Err.Error())
	}

	return message
}

// String is an alias of Error()
func (e *ErrMessage) String() string {
	return e.Error()
}

// Format implements fmt.Formatter interface
func (e *ErrMessage) Format(s fmt.State, verb rune) {
	message := fmt.Sprintf("%s: %s\n", e.Code(), e.Raw)
	if e.Err != nil {
		message += e.Err.Error()
	}

	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = io.WriteString(s, fmt.Sprintf("%s: %s\n", e.Code(), e.Raw))

			if e.Err != nil {
				merr, IsMulti := e.Err.(errors.ErrorGroup)
				if IsMulti {
					_, _ = io.WriteString(s, fmt.Sprintf("%+v", merr))
					return
				}

				em, IsEM := e.Err.(*ErrMessage)
				if IsEM {
					em.Format(s, verb)
					return
				}

				_, _ = io.WriteString(s, fmt.Sprintf("%s\n", e.Err.Error()))
				if errors.HasStack(e.Err) {
					errors.GetStackTracer(e.Err).StackTrace().Format(s, verb)
				}

				return
			}

			e.Stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		_, _ = io.WriteString(s, message)
	case 'q':
		_, _ = fmt.Fprintf(s, "%q", message)
	}
}

// SetError sets Err with the given error
func (e *ErrMessage) SetError(err error) *ErrMessage {
	e.Err = err

	return e
}

// Renew returns a new *ErrMessage and specify with given input
func (e *ErrMessage) Renew(ins ...interface{}) *ErrMessage {
	c := e.Clone()
	c.Specify(ins...)

	return c
}

// Clone returns a new *ErrMessage with same member variables
func (e *ErrMessage) Clone() *ErrMessage {
	return newErrMessage(e.Header, e.ErrCode, e.Raw, e.Err, e.Stack)
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
