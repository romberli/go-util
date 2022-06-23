package rabbitmq

import (
	"strconv"
	"testing"

	"github.com/romberli/go-util/common"
	"github.com/romberli/go-util/constant"
	"github.com/romberli/log"
	"github.com/stretchr/testify/assert"
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

	msg := testProducer.BuildMessage(testMessage, constant.DefaultJSONContentType)
	asst.Equal(testMessage, string(msg.Body), "test BuildMessage() failed")
}

func TestProducer_BuildMessageWithExpiration(t *testing.T) {
	asst := assert.New(t)

	msg := testProducer.BuildMessageWithExpiration(testMessage, constant.DefaultJSONContentType, testExpiration)
	asst.Equal(testMessage, string(msg.Body), "test BuildMessageWithExpiration() failed")
	asst.Equal(strconv.Itoa(testExpiration), msg.Expiration, "test BuildMessageWithExpiration() failed")
}

func TestProducer_Publish(t *testing.T) {
	asst := assert.New(t)

	err := testProducer.Publish(testExchangeName, testKey, testProducer.BuildMessageWithExpiration(testMessage, constant.DefaultJSONContentType, testExpiration))
	asst.Nil(err, common.CombineMessageWithError("test Publish() failed", err))
}
