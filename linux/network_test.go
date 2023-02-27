package linux

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testValidIP1  = "127.0.0.1"
	testValidIP2  = "192.168.137.11"
	testInvalidIP = "127.0.0.256"
)

func TestNetwork_All(t *testing.T) {
	TestNetwork_GetDefaultIP(t)

}

func TestNetwork_GetDefaultIP(t *testing.T) {
	asst := assert.New(t)

	hostIP, err := GetDefaultIP()
	t.Logf("host ip: %s", hostIP)
	asst.Nil(err, "test GetDefaultIP failed")
}

func TestNetwork_IsValidIP(t *testing.T) {
	asst := assert.New(t)

	asst.True(IsValidIP(testValidIP1), "test IsValidIP() failed")
	asst.False(IsValidIP(testInvalidIP), "test IsValidIP() failed")
}

func TestNetwork_CompareIP(t *testing.T) {
	asst := assert.New(t)

	result, err := CompareIP(testValidIP1, testValidIP2)
	asst.Nil(err, "test CompareIP() failed")
	asst.Equal(-1, result, "test CompareIP() failed")
}
