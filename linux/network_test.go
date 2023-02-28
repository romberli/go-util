package linux

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testValidIP1  = "127.0.0.1"
	testValidIP2  = "192.168.137.11"
	testInvalidIP = "127.0.0.256"
	testAddr1     = "127.0.0.1:3306"
	testAddr2     = "127.0.0.1:3307"
	testAddr3     = "127.0.0.1:3305"
)

func TestNetwork_All(t *testing.T) {
	TestNetwork_GetDefaultIP(t)
	TestNetwork_IsValidIP(t)
	TestNetwork_CompareIP(t)
	TestNetwork_CompareAddr(t)
	TestNetwork_GetMinAddr(t)
	TestNetwork_SortAddrs(t)
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

func TestNetwork_CompareAddr(t *testing.T) {
	asst := assert.New(t)

	result, err := CompareAddr(testAddr1, testAddr2)
	asst.Nil(err, "test CompareAddr() failed")
	asst.Equal(-1, result, "test CompareAddr() failed")
}

func TestNetwork_GetMinAddr(t *testing.T) {
	asst := assert.New(t)

	result, err := GetMinAddr([]string{testAddr1, testAddr2})
	asst.Nil(err, "test GetMinAddr() failed")
	asst.Equal(testAddr1, result, "test GetMinAddr() failed")
}

func TestNetwork_SortAddrs(t *testing.T) {
	asst := assert.New(t)

	addrs := []string{testAddr1, testAddr2, testAddr3}
	err := SortAddrs(addrs)
	asst.Nil(err, "test SortAddrs() failed")
	asst.Equal([]string{testAddr3, testAddr1, testAddr2}, addrs, "test SortAddrs() failed")
}
