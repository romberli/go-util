package common

import (
	"testing"
	"time"

	"github.com/pingcap/errors"
	"github.com/romberli/log"
	"github.com/stretchr/testify/assert"

	"github.com/romberli/go-util/constant"
)

const (
	testMaxRetryCount = 3
	testMaxWaitTime   = 1000 * time.Second
)

type testStruct struct {
	index int
}

func newTestStruct() *testStruct {
	return &testStruct{index: constant.ZeroInt}
}

func (t *testStruct) testFunc() error {
	if t.index < testMaxRetryCount {
		t.index++
		return errors.New("test error")
	}

	return nil
}

func TestRetry(t *testing.T) {
	asst := assert.New(t)

	log.SetDisableEscape(true)

	ts := newTestStruct()
	option := NewRetryOption(testMaxRetryCount, DefaultDelayTime, testMaxWaitTime, log.L())

	err := Retry(ts.testFunc, option)
	asst.Nil(err, "test Retry() failed")
}
