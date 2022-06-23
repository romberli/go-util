package rabbitmq

import (
	"testing"

	"github.com/romberli/go-util/common"
	"github.com/romberli/log"
	"github.com/stretchr/testify/assert"
)

var testConsumer *Consumer

func init() {
	testConn = testCreateConn(testAddr, testUser, testPass)
	testConsumer = testCreateConsumer(testConn)
}

func testCreateConsumer(conn *Conn) *Consumer {
	var err error

	testConsumer, err = NewConsumerWithConn(conn)
	if err != nil {
		log.Errorf("creating new Consumer failed. %s", err)
	}

	return testConsumer
}

func TestConsumer_All(t *testing.T) {
	TestConsumer_Consume(t)
}

func TestConsumer_Consume(t *testing.T) {
	asst := assert.New(t)

	deliveryChan, err := testConsumer.Consume(testQueueName, testConsumerName, testExclusive)
	asst.Nil(err, common.CombineMessageWithError("test Consume() failed", err))
	for d := range deliveryChan {
		log.Infof(" [x] %s", d.Body)
	}
}

func TestConsumer_Ack(t *testing.T) {
	asst := assert.New(t)

	deliveryChan, err := testConsumer.Consume(testQueueName, testConsumerName, testExclusive)
	asst.Nil(err, common.CombineMessageWithError("test Ack() failed", err))
	for d := range deliveryChan {
		log.Infof(" [x] %s", d.Body)
		err = testConsumer.Ack(d.DeliveryTag, testMultiple)
		asst.Nil(err, common.CombineMessageWithError("test Ack() failed", err))
	}
}
