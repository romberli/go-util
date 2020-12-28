package linux

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidHost(t *testing.T) {
	asst := assert.New(t)

	hostIP, err := GetDefaultIP()
	t.Logf("host ip: %s", hostIP)
	asst.Nil(err, "test GetDefaultIP failed")
}
