package rabbitmq

import (
	"github.com/romberli/log"
)

const (
	testAddr  = "192.168.137.11:5672"
	testUser  = "guest"
	testPass  = "guest"
	testVhost = "/test_vhost"

	testExchangeName    = "test_exchange"
	testExchangeType    = "topic"
	testQueueName       = "test_queue"
	testKey             = "test_key"
	testMessage         = `{"dbs": {"id": 1, "db_name": "test_db", "cluster_id": 1}}`
	testMessageTemplate = `{"dbs": {"id": %d, "db_name": "test_db", "cluster_id": 1}}`
	testExpiration      = 1000 * 60 * 60 * 5 // 5 minutes
	testPublishCount    = 5

	testConsumerName  = "test_consumer"
	testPrefetchCount = 3
	testGlobal        = true
	testExclusive     = true
	testMultiple      = true
	testRequeue       = true
)

var testConn *Conn

func init() {
	testConn = testCreateConn(testAddr, testUser, testPass)
}

// testCreateConn returns a new *Conn with given address, user and password
func testCreateConn(addr, user, pass string) *Conn {
	var err error

	testConn, err = NewConnWithDefault(addr, user, pass)
	if err != nil {
		log.Errorf("creating new Connection failed. %s", err)
	}

	return testConn
}
