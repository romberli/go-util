package common

import (
	"testing"

	"github.com/pingcap/errors"
	"github.com/romberli/log"
	"github.com/stretchr/testify/assert"
)

func testFunc() error {
	return errors.New("test error")
}

func TestRetry(t *testing.T) {
	asst := assert.New(t)

	option := NewRetryOptionWithDefault()
	log.SetDisableEscape(true)

	err := Retry(testFunc, option)
	asst.NotNil(err, "test Retry() failed")
}
