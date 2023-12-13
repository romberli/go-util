package rabbitmq

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/romberli/log"
	"github.com/stretchr/testify/assert"

	"github.com/romberli/go-util/common"
	"github.com/romberli/go-util/constant"
)

var testProducer *Producer

func init() {
	testConn = testCreateConn(testAddr, testUser, testPass)
	testProducer = testCreateProducer(testConn)
}

func testCreateProducer(conn *Conn) *Producer {
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
		err := testProducer.Publish(testExchangeName, testKey, testProducer.BuildMessageWithExpiration(message, constant.DefaultJSONContentType, testExpiration))
		asst.Nil(err, common.CombineMessageWithError("test Publish() failed", err))
	}
}

func TestProducer_PublishWithContext(t *testing.T) {
	asst := assert.New(t)

	for i := constant.ZeroInt; i < testPublishCount; i++ {
		message := fmt.Sprintf(testMessageTemplate, i)
		err := testProducer.PublishWithContext(context.Background(), testExchangeName, testKey, testProducer.BuildMessageWithExpiration(message, constant.DefaultJSONContentType, testExpiration))
		asst.Nil(err, common.CombineMessageWithError("test Publish() failed", err))
	}
}
