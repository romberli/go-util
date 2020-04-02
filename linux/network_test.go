package linux

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidHost(t *testing.T) {
	assert := assert.New(t)

	host_ip, err := GetDefaultIP()
	t.Logf("host ip: %s", host_ip)
	assert.Nil(err, "test GetDefaultIP failed")

	isValid := IsValidHost(host_ip)
	assert.True(isValid, "test IsValidHost failed")
}
