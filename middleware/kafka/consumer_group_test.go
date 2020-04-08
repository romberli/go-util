package kafka

import (
	"context"
	"github.com/Shopify/sarama"
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

	kafkaVersion := "2.2.0"
	brokerList := []string{"10.0.0.63:9092", "10.0.0.84:9092", "10.0.0.92:9092"}
	groupName := "group001"
	topicName := "test001"
	ctx, cancel := context.WithCancel(context.Background())
	handler := DefaultConsumerGroupHandler{}

	cg, err = NewConsumerGroup(kafkaVersion, brokerList, groupName, sarama.OffsetNewest)
	assert.Nil(err, "create consumer group failed. group: %s, topic: %s", groupName, topicName)

	go func() {
		err = cg.Consume(ctx, topicName, handler)
		assert.Nil(err, "consume failed. group: %s, topic: %s", groupName, topicName)

		err = ctx.Err()
		assert.EqualError(err, "context canceled", "context error is not nil. group: %s, topic: %s, message: %s", groupName, topicName, err.Error())
	}()

	time.Sleep(DefaultConsumeSeconds)

	cancel()

	err = ctx.Err()
	assert.EqualError(err, "context canceled", "context error is not nil. group: %s, topic: %s, message: %s", groupName, topicName, err.Error())
}
