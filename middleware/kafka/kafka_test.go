package kafka

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const (
	DefaultConsumeSeconds = 5 * time.Second
)

func TestConsume(t *testing.T) {
	var (
		err error
		cg  *ConsumerGroup
	)

	assert := assert.New(t)

	kafkaVersion := "V1.0.0.0"
	brokerList := []string{"192.168.137.11:9092"}
	groupName := "group001"
	topicName := "topic001"
	ctx, cancel := context.WithCancel(context.Background())
	handler := DefaultConsumerGroupHandler{}

	cg, err = NewConsumerGroup(ctx, kafkaVersion, brokerList, groupName)
	assert.Nil(err, "create consumer group failed.")

	go func() {
		err = cg.Consume(topicName, handler)
		assert.Nil(err, "consume failed. group: %s, topic: %s")

		err = cg.Ctx.Err()
		assert.Nil(err, "context error is not nil. group: %s, topic: %s")
	}()

	time.Sleep(DefaultConsumeSeconds)

	cancel()

	err = cg.Ctx.Err()
	assert.Nil(err, "context error is not nil. group: %s, topic: %s")
}
