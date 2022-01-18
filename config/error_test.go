package config

import (
	"fmt"
	"testing"

	"github.com/pingcap/errors"
	"github.com/romberli/log"
	"github.com/stretchr/testify/assert"

	"github.com/romberli/go-multierror"
)

const (
	testHeader   = "testEM"
	testErrCode1 = 100001
	testErrCode2 = 100002
	testRaw1     = "error message 1. id: %d"
	testRaw2     = "error message 2. id: %d"
)

func ReturnNotNil1() error {
	var em *ErrMessage
	return em
}

func ReturnNotNil2() error {
	em := &ErrMessage{}
	return em
}

func ReturnNil1() error {
	var em *ErrMessage
	return em.ErrorOrNil()
}

func ReturnNil2() error {
	em := &ErrMessage{}
	return em.ErrorOrNil()
}

func TestError(t *testing.T) {
	var (
		err        error
		errMessage *ErrMessage
		header     string
		errCode    int
		raw        string
		line       int

		expectString string
	)

	asst := assert.New(t)

	header = "test"
	errCode = 100001
	raw = "something goes wrong, line: %d"
	line = 100

	errMessage = NewErrMessage(header, errCode, raw, nil)

	t.Log("==========test Code() started.==========")
	expectCode := fmt.Sprintf("%s-%d", header, errCode)
	asst.Equal(expectCode, errMessage.Code(), "test Code() failed.")
	t.Log("==========test Code() completed.==========")

	t.Log("==========test Error() started.==========")
	expectString = fmt.Sprintf("%s-%d: %s\n", header, errCode, raw)
	asst.Equal(expectString, errMessage.Error(), "test Error() failed.")
	t.Log("==========test Error() completed.==========")

	t.Log("==========test Renew() started.==========")
	expectString = fmt.Sprintf("%s-%d: %s\n", header, errCode, fmt.Sprintf(raw, line))
	asst.Equal(expectString, errMessage.Renew(line).Error(), "test Renew() failed.")
	t.Log("==========test Renew() completed.==========")

	t.Log("==========test ErrorOrNil() started.==========")
	err = ReturnNotNil1()
	if err == nil {
		asst.Fail("test ErrorOrNil() failed.")
	}
	err = ReturnNotNil2()
	if err == nil {
		asst.Fail("test ErrorOrNil() failed.")
	}
	err = ReturnNil1()
	if err != nil {
		asst.Fail("test ErrorOrNil() failed.")
	}
	err = ReturnNil2()
	if err != nil {
		asst.Fail("test ErrorOrNil() failed.")
	}
	t.Log("==========test ErrorOrNil() completed.==========")
}

func funcA() error {
	return errors.New("function error")
}

func funcB() error {
	err := funcA()
	// raw := "function error"
	// return NewErrMessage("FUNCB", 100001, raw, err)
	en := errors.New("function b error")
	merr := &multierror.Error{}
	merr = multierror.Append(merr, errors.Trace(err), en)

	return errors.Trace(merr)
}

func funcC() error {
	return funcB()
}

func TestError_Log(t *testing.T) {
	log.SetDisableEscape(true)
	log.SetDisableDoubleQuotes(true)

	err := funcC()
	log.Errorf("got error: %+v", err)
	// log.Errorf("got error: %s", err.Error())
	// log.Error("got error", zap.Error(err))
}

func TestNestedError(t *testing.T) {
	log.SetDisableEscape(true)
	log.SetDisableDoubleQuotes(true)

	err := funcC()
	em1 := NewErrMessage(testHeader, testErrCode1, testRaw1, err).Renew(1)
	em2 := NewErrMessage(testHeader, testErrCode2, testRaw2, em1).Renew(2)
	log.Errorf("got error: %+v", em2)
}
