package rabbitmq

import (
	"testing"
	"time"

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
	TestConsumer_ExchangeDeclare(t)
	TestConsumer_QueueDeclare(t)
	TestConsumer_QueueBind(t)
	TestConsumer_Qos(t)
	TestConsumer_Consume(t)
	TestConsumer_Cancel(t)
	TestConsumer_Ack(t)
	TestConsumer_Nack(t)
}

func TestConsumer_ExchangeDeclare(t *testing.T) {
	asst := assert.New(t)

	err := testConsumer.ExchangeDeclare(testExchangeName, testExchangeType)
	asst.Nil(err, common.CombineMessageWithError("test ExchangeDeclare() failed", err))
}

func TestConsumer_QueueDeclare(t *testing.T) {
	asst := assert.New(t)

	queue, err := testConsumer.QueueDeclare(testQueueName)
	asst.Nil(err, common.CombineMessageWithError("test QueueDeclare() failed", err))
	asst.Equal(testQueueName, queue.Name, "test QueueDeclare() failed")
}

func TestConsumer_QueueBind(t *testing.T) {
	asst := assert.New(t)

	err := testConsumer.QueueBind(testQueueName, testExchangeName, testKey)
	asst.Nil(err, common.CombineMessageWithError("test QueueBind() failed", err))
}

func TestConsumer_Qos(t *testing.T) {
	asst := assert.New(t)

	err := testConsumer.Qos(testPrefetchCount, testGlobal)
	asst.Nil(err, common.CombineMessageWithError("test Qos() failed", err))
}

func TestConsumer_Consume(t *testing.T) {
	asst := assert.New(t)

	deliveryChan, err := testConsumer.Consume(testQueueName, testExclusive)
	asst.Nil(err, common.CombineMessageWithError("test Consume() failed", err))

	for {
		select {
		case d := <-deliveryChan:
			log.Infof("%s", d.Body)
		default:
			log.Infof("no message to consume, will exit now")
			err = testConsumer.Cancel()
			asst.Nil(err, common.CombineMessageWithError("test Consume() failed", err))
			return
		}
	}
}

func TestConsumer_Cancel(t *testing.T) {
	asst := assert.New(t)

	_, err := testConsumer.Consume(testQueueName, testExclusive)
	asst.Nil(err, common.CombineMessageWithError("test Ack() failed", err))
	err = testConsumer.Cancel()
	asst.Nil(err, common.CombineMessageWithError("test Cancel() failed", err))
}

func TestConsumer_Ack(t *testing.T) {
	asst := assert.New(t)

	deliveryChan, err := testConsumer.Consume(testQueueName, testExclusive)
	asst.Nil(err, common.CombineMessageWithError("test Ack() failed", err))
	for {
		select {
		case d := <-deliveryChan:
			log.Infof("%s", d.Body)
			err = testConsumer.Ack(d.DeliveryTag, testMultiple)
			asst.Nil(err, common.CombineMessageWithError("test Ack() failed", err))
		default:
			log.Infof("no message to consume, will exit now")
			err = testConsumer.Cancel()
			asst.Nil(err, common.CombineMessageWithError("test Ack() failed", err))
			return
		}
	}
}

func TestConsumer_Nack(t *testing.T) {
	asst := assert.New(t)

	for {
		deliveryChan, err := testConsumer.Consume(testQueueName, testExclusive)
		// asst.Nil(err, common.CombineMessageWithError("test Nack() failed", err))
		if err != nil {
			if testConsumer.IsExclusiveUseError(testQueueName, err) {
				log.Infof("queue %s is exclusive used, will sleep 3 seconds and try again", testQueueName)
				time.Sleep(time.Second * 3)
				continue
			}
		}
		for {
			select {
			case d, ok := <-deliveryChan:
				if !ok {
					log.Infof("delivery chan has been closed")
					time.Sleep(time.Second * 3)
					continue
				}
				log.Infof("%s", d.Body)
				// err = testConsumer.Nack(d.DeliveryTag, testMultiple, testRequeue)
				asst.Nil(err, common.CombineMessageWithError("test Nack() failed", err))
				time.Sleep(time.Second * 3)
			default:
				log.Infof("no message to consume...")
				time.Sleep(time.Second * 3)
				continue
				// log.Infof("no message to consume, will exit now")
				// err = testConsumer.Cancel()
				// asst.Nil(err, common.CombineMessageWithError("test Nack() failed", err))
				// return
			}
		}
	}
}

func TestConsumer_Nack1(t *testing.T) {
	asst := assert.New(t)

	for {
		deliveryChan, err := testConsumer.Consume(testQueueName, testExclusive)
		// asst.Nil(err, common.CombineMessageWithError("test Nack() failed", err))
		if err != nil {
			if testConsumer.IsExclusiveUseError(testQueueName, err) {
				log.Infof("queue %s is exclusive used, will sleep 3 seconds and try again", testQueueName)
				time.Sleep(time.Second * 3)
				continue
			}

			if testConsumer.IsChannelOrConnectionClosedError(err) {
				log.Infof("channel or connection is closed, will open connection and channel again")

				testConn = testCreateConn(testAddr, testUser, testPass)
				testConsumer = testCreateConsumer(testConn)
				time.Sleep(time.Second * 3)
				continue
			}
		}

		for {
			select {
			case d, ok := <-deliveryChan:
				if !ok {
					log.Infof("delivery chan has been closed")
					time.Sleep(time.Second * 3)
					continue
				}
				log.Infof("%s", d.Body)
				// err = testConsumer.Nack(d.DeliveryTag, testMultiple, testRequeue)
				asst.Nil(err, common.CombineMessageWithError("test Nack() failed", err))
				time.Sleep(time.Second * 3)
			default:
				log.Infof("no message to consume...")
				time.Sleep(time.Second * 3)
				continue
				// log.Infof("no message to consume, will exit now")
				// err = testConsumer.Cancel()
				// asst.Nil(err, common.CombineMessageWithError("test Nack() failed", err))
				// return
			}
		}
	}
}
