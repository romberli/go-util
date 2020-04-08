package kafka

import (
	"context"
	"github.com/romber2001/log"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const (
	DefaultProduceSeconds = 5 * time.Second
)

func TestProduce(t *testing.T) {
	var (
		err error
		p   *AsyncProducer
		ts  string
	)

	assert := assert.New(t)

	kafkaVersion := "2.2.0"
	brokerList := []string{"10.0.0.63:9092", "10.0.0.84:9092", "10.0.0.92:9092"}
	topicName := "test001"
	ctx, cancel := context.WithCancel(context.Background())

	p, err = NewAsyncProducer(kafkaVersion, brokerList)
	assert.Nil(err, "create producer failed.")

	defer func() {
		err = p.Close()
		if err != nil {
			log.Errorf("close producer failed. topic: %s, message: %s", topicName, err.Error())
		}
	}()

	go func() {
		for i := 0; i < 10; i++ {
			ts = time.Now().String()
			err = p.Produce(topicName, ts)
			assert.Nil(err, "produce failed. topic: %s, message: %s", topicName, ts)

			err = ctx.Err()
			assert.Nil(err, "context error is not nil. topic: %s, errMessage: %v", topicName, err)
		}
	}()

	time.Sleep(DefaultProduceSeconds)

	cancel()

	err = ctx.Err()
	assert.EqualError(err, "context canceled", "context error is not nil. topic: %s", topicName)
}
