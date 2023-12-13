package consumer

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/romberli/log"
	"github.com/stretchr/testify/assert"

	"github.com/romberli/go-util/common"
	"github.com/romberli/go-util/constant"
	"github.com/romberli/go-util/middleware/rabbitmq/producer"
)

var (
	testConsumerPool *Pool
)

func init() {
	testConsumerPool = testCreateConsumerPool(testAddr, testUser, testPass, testVhost, testTag)
}

func testCreateConsumerPool(addr, user, pass, vhost, tag string) *Pool {
	var err error

	testConsumerPool, err = NewPoolWithDefault(addr, user, pass, vhost, tag)
	if err != nil {
		log.Errorf("creating new consumer pool failed. %s", err)
	}

	return testConsumerPool
}

func TestPool(t *testing.T) {
	var (
		err error
		pc  *PoolConsumer
	)

	asst := assert.New(t)

	// get connection from the pool
	pc, err = testConsumerPool.Get()
	asst.Nil(err, "get pool consumer from pool failed")

	// sleep to test pool maintaining mechanism
	t.Logf("sleep 10 seconds to test pool maintaining mechanism")
	time.Sleep(10 * time.Second)

	err = pc.Close()
	asst.Nil(err, "close pool consumer failed")

	// create producer
	p, err := producer.NewProducer(testAddr, testUser, testPass, testVhost, testTag)
	asst.Nil(err, common.CombineMessageWithError("create producer failed", err))
	defer func() {
		err = p.Disconnect()
		asst.Nil(err, common.CombineMessageWithError("close producer failed", err))
	}()
	// send message to the queue
	for i := constant.ZeroInt; i < testPublishCount; i++ {
		message := fmt.Sprintf(testMessageTemplate, i)
		err = p.PublishWithContext(context.Background(), testExchangeName, testKey,
			p.BuildMessageWithExpiration(constant.DefaultJSONContentType, message, testExpiration))
		asst.Nil(err, common.CombineMessageWithError("test Publish() failed", err))
	}

	// get connection from the pool
	pc, err = testConsumerPool.Get()
	asst.Nil(err, "get pool consumer from pool failed")

	deliveryChan, err := pc.Consume(testQueueName, testExclusive)
	asst.Nil(err, common.CombineMessageWithError("test Consume() failed", err))

	expireTime := time.Now().Add(testMaxWaitTime)

	for {
		select {
		case d := <-deliveryChan:
			t.Logf("%s", d.Body)
			err = pc.Ack(d.DeliveryTag, testMultiple)
			asst.Nil(err, common.CombineMessageWithError("test Ack() failed", err))
		default:
			if time.Now().Sub(expireTime) > constant.ZeroInt {
				t.Logf("no message to consume, will exit now")
				err = testConsumer.Cancel()
				asst.Nil(err, common.CombineMessageWithError("test Consume() failed", err))
				return
			}

			t.Logf("no message to consume, will sleep 1 seconds and try again")
			time.Sleep(1 * time.Second)
		}
	}
}
