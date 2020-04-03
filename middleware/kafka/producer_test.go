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

	p, err = NewAsyncProducer(ctx, kafkaVersion, brokerList)
	assert.Nil(err, "create producer failed.")

	defer func() {
		err = p.Producer.Close()
		if err != nil {
			log.Errorf("close producer failed. topic: %s, message: %s", topicName, err.Error())
		}
	}()

	go func() {
		for i := 0; i < 10; i++ {
			ts = time.Now().String()
			err = p.Produce(topicName, ts)
			assert.Nil(err, "produce failed. topic: %s, message: %s", topicName, ts)

			err = p.Ctx.Err()
			assert.EqualError(err, "context canceled", "context error is not nil. topic: %s", topicName)
		}
	}()

	time.Sleep(DefaultProduceSeconds)

	cancel()

	err = p.Ctx.Err()
	assert.EqualError(err, "context canceled", "context error is not nil. topic: %s", topicName)
}
