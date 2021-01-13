package config

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
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

	errMessage = newErrMessage(header, errCode, raw)

	t.Log("==========test Code() started.==========")
	expectCode := fmt.Sprintf("%s-%d", header, errCode)
	asst.Equal(expectCode, errMessage.Code(), "test Code() failed.")
	t.Log("==========test Code() completed.==========")

	t.Log("==========test Error() started.==========")
	expectString = fmt.Sprintf("%s-%d: %s", header, errCode, raw)
	asst.Equal(expectString, errMessage.Error(), "test Error() failed.")
	t.Log("==========test Error() completed.==========")

	t.Log("==========test Renew() started.==========")
	expectString = fmt.Sprintf("%s-%d: %s", header, errCode, fmt.Sprintf(raw, line))
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
