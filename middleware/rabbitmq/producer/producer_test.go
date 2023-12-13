package producer

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/romberli/log"
	"github.com/stretchr/testify/assert"

	"github.com/romberli/go-util/common"
	"github.com/romberli/go-util/constant"
	"github.com/romberli/go-util/middleware/rabbitmq/client"
)

const (
	testAddr  = "192.168.137.11:5672"
	testUser  = "guest"
	testPass  = "guest"
	testVhost = "/"
	testTag   = "test_consumer"

	testExchangeName    = "test_exchange"
	testExchangeType    = "topic"
	testQueueName       = "test_queue"
	testKey             = "test_key"
	testMessage         = `{"dbs": {"id": 1, "db_name": "test_db", "cluster_id": 1}}`
	testMessageTemplate = `{"dbs": {"id": %d, "db_name": "test_db", "cluster_id": 1}}`
	testExpiration      = 1000 * 60 * 60 * 5 // 5 minutes
	testPublishCount    = 5

	testExclusive = true
	testMultiple  = true

	testMaxWaitTime = 10 * time.Second
)

var (
	testConn     *client.Conn
	testProducer *Producer
)

func init() {
	testConn = testCreateConn(testAddr, testUser, testPass)
	testProducer = testCreateProducer(testConn)
}

// testCreateConn returns a new *Conn with given address, user and password
func testCreateConn(addr, user, pass string) *client.Conn {
	var err error

	testConn, err = client.NewConn(addr, user, pass, testVhost, testTag)
	if err != nil {
		log.Errorf("creating new Connection failed. %s", err)
	}

	return testConn
}

func testCreateProducer(conn *client.Conn) *Producer {
	var err error

	producer, err := NewProducerWithConn(conn)
	if err != nil {
		log.Errorf("creating new Producer failed. %s", err)
	}

	return producer
}

func TestProducer_All(t *testing.T) {
	TestProducer_ExchangeDeclare(t)
	TestProducer_QueueDeclare(t)
	TestProducer_QueueBind(t)
	TestProducer_BuildMessage(t)
	TestProducer_BuildMessageWithExpiration(t)
	TestProducer_Publish(t)
	TestProducer_PublishWithContext(t)
}

func TestProducer_ExchangeDeclare(t *testing.T) {
	asst := assert.New(t)

	err := testProducer.ExchangeDeclare(testExchangeName, testExchangeType)
	asst.Nil(err, common.CombineMessageWithError("test ExchangeDeclare() failed", err))
}

func TestProducer_QueueDeclare(t *testing.T) {
	asst := assert.New(t)

	queue, err := testProducer.QueueDeclare(testQueueName)
	asst.Nil(err, common.CombineMessageWithError("test QueueDeclare() failed", err))
	asst.Equal(testQueueName, queue.Name, "test QueueDeclare() failed")
}

func TestProducer_QueueBind(t *testing.T) {
	asst := assert.New(t)

	err := testProducer.QueueBind(testQueueName, testExchangeName, testKey)
	asst.Nil(err, common.CombineMessageWithError("test QueueBind() failed", err))
}

func TestProducer_BuildMessage(t *testing.T) {
	asst := assert.New(t)

	msg := testProducer.BuildMessage(constant.DefaultJSONContentType, testMessage)
	asst.Equal(testMessage, string(msg.Body), "test BuildMessage() failed")
}

func TestProducer_BuildMessageWithExpiration(t *testing.T) {
	asst := assert.New(t)

	msg := testProducer.BuildMessageWithExpiration(constant.DefaultJSONContentType, testMessage, testExpiration)
	asst.Equal(testMessage, string(msg.Body), "test BuildMessageWithExpiration() failed")
	asst.Equal(strconv.Itoa(testExpiration), msg.Expiration, "test BuildMessageWithExpiration() failed")
}

func TestProducer_Publish(t *testing.T) {
	asst := assert.New(t)

	for i := constant.ZeroInt; i < testPublishCount; i++ {
		message := fmt.Sprintf(testMessageTemplate, i)
		err := testProducer.Publish(testExchangeName, testKey, testProducer.BuildMessageWithExpiration(constant.DefaultJSONContentType, message, testExpiration))
		asst.Nil(err, common.CombineMessageWithError("test Publish() failed", err))
	}
}

func TestProducer_PublishWithContext(t *testing.T) {
	asst := assert.New(t)

	for i := constant.ZeroInt; i < testPublishCount; i++ {
		message := fmt.Sprintf(testMessageTemplate, i)
		err := testProducer.PublishWithContext(context.Background(), testExchangeName, testKey, testProducer.BuildMessageWithExpiration(constant.DefaultJSONContentType, message, testExpiration))
		asst.Nil(err, common.CombineMessageWithError("test Publish() failed", err))
	}
}
