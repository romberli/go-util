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
	testProducer *Producer
)

func init() {
	testProducer = testCreateProducer()
}

func testCreateProducer() *Producer {
	var err error

	producer, err := NewProducer(testAddr, testUser, testPass, testVhost, testTag, testExchangeName, testQueueName, testKey)
	if err != nil {
		log.Errorf("creating new producer failed. %s", err)
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

	err := testProducer.ExchangeDeclare(testExchangeType)
	asst.Nil(err, common.CombineMessageWithError("test ExchangeDeclare() failed", err))
}

func TestProducer_QueueDeclare(t *testing.T) {
	asst := assert.New(t)

	err := testProducer.QueueDeclare()
	asst.Nil(err, common.CombineMessageWithError("test QueueDeclare() failed", err))
	asst.Equal(testQueueName, testProducer.Queue.Name, "test QueueDeclare() failed")
}

func TestProducer_QueueBind(t *testing.T) {
	asst := assert.New(t)

	err := testProducer.QueueBind()
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
		err := testProducer.Publish(testProducer.BuildMessageWithExpiration(constant.DefaultJSONContentType, message, testExpiration))
		asst.Nil(err, common.CombineMessageWithError("test Publish() failed", err))
	}
}

func TestProducer_PublishWithContext(t *testing.T) {
	asst := assert.New(t)

	for i := constant.ZeroInt; i < testPublishCount; i++ {
		message := fmt.Sprintf(testMessageTemplate, i)
		err := testProducer.PublishWithContext(context.Background(), testProducer.BuildMessageWithExpiration(constant.DefaultJSONContentType, message, testExpiration))
		asst.Nil(err, common.CombineMessageWithError("test Publish() failed", err))
	}
}
