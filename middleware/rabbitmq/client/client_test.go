package client

import (
	"github.com/romberli/log"
)

const (
	testAddr  = "192.168.137.11:5672"
	testUser  = "guest"
	testPass  = "guest"
	testVhost = "/"
	testTag   = "test_consumer"
)

var testConn *Conn

func init() {
	testConn = testCreateConn(testAddr, testUser, testPass)
}

// testCreateConn returns a new *Conn with given address, user and password
func testCreateConn(addr, user, pass string) *Conn {
	var err error

	testConn, err = NewConn(addr, user, pass, testVhost)
	if err != nil {
		log.Errorf("creating new Connection failed. %s", err)
	}

	return testConn
}
