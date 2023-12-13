package producer

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/romberli/log"
	"github.com/stretchr/testify/assert"

	"github.com/romberli/go-util/common"
	"github.com/romberli/go-util/constant"
	"github.com/romberli/go-util/middleware/rabbitmq/consumer"
)

var (
	testProducerPool *Pool
)

func init() {
	testProducerPool = testCreateProducerPool(testAddr, testUser, testPass, testVhost, testTag)
}

func testCreateProducerPool(addr, user, pass, vhost, tag string) *Pool {
	var err error

	testProducerPool, err = NewPoolWithDefault(addr, user, pass, vhost, tag)
	if err != nil {
		log.Errorf("creating new producer pool failed. %s", err)
	}

	return testProducerPool
}

func TestPool(t *testing.T) {
	var (
		err error
		pp  *PoolProducer
	)

	asst := assert.New(t)

	// get connection from the pool
	pp, err = testProducerPool.Get()
	asst.Nil(err, "get pool producer from pool failed")

	// sleep to test pool maintaining mechanism
	t.Logf("sleep 10 seconds to test pool maintaining mechanism")
	time.Sleep(10 * time.Second)

	err = pp.Close()
	asst.Nil(err, "close pool producer failed")

	// get connection from the pool
	pp, err = testProducerPool.Get()
	asst.Nil(err, "get pool producer from pool failed")
	defer func() {
		err = pp.Close()
		asst.Nil(err, "close pool producer failed")
	}()
	// send message to the queue
	for i := constant.ZeroInt; i < testPublishCount; i++ {
		message := fmt.Sprintf(testMessageTemplate, i)
		err = pp.PublishWithContext(context.Background(), testExchangeName, testKey,
			pp.BuildMessageWithExpiration(constant.DefaultJSONContentType, message, testExpiration))
		asst.Nil(err, common.CombineMessageWithError("test Publish() failed", err))
	}

	// create consumer
	c, err := consumer.NewConsumer(testAddr, testUser, testPass, testVhost, testTag)
	asst.Nil(err, common.CombineMessageWithError("create consumer failed", err))
	defer func() {
		err = c.Disconnect()
		asst.Nil(err, common.CombineMessageWithError("close consumer failed", err))
	}()
	// consume message from the queue
	deliveryChan, err := c.Consume(testQueueName, testExclusive)
	asst.Nil(err, common.CombineMessageWithError("test Consume() failed", err))

	expireTime := time.Now().Add(testMaxWaitTime)

	for {
		select {
		case d := <-deliveryChan:
			t.Logf("%s", d.Body)
		default:
			if time.Now().Sub(expireTime) > constant.ZeroInt {
				t.Logf("no message to consume, will exit now")
				err = c.Cancel()
				asst.Nil(err, common.CombineMessageWithError("test Consume() failed", err))
				return
			}

			t.Logf("no message to consume, will sleep 1 seconds and try again")
			time.Sleep(1 * time.Second)
		}
	}
}
