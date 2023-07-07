package kafka

import (
	"context"
	"testing"
	"time"

	"github.com/Shopify/sarama"
	"github.com/stretchr/testify/assert"
)

const (
	DefaultConsumeTime = 5 * time.Second
)

func TestConsume(t *testing.T) {
	var (
		err error
		cg  *ConsumerGroup
	)

	asst := assert.New(t)

	kafkaVersion := "2.2.0"
	brokerList := []string{"10.0.0.63:9092", "10.0.0.84:9092", "10.0.0.92:9092"}
	groupName := "group001"
	topicName := "test001"
	ctx, cancel := context.WithCancel(context.Background())
	handler := DefaultConsumerGroupHandler{}

	cg, err = NewConsumerGroup(kafkaVersion, brokerList, groupName, sarama.OffsetNewest)
	asst.Nil(err, "create consumer group failed. group: %s, topic: %s", groupName, topicName)

	go func() {
		err = cg.Consume(ctx, topicName, handler)
		asst.Nil(err, "consume failed. group: %s, topic: %s", groupName, topicName)

		err = ctx.Err()
		asst.EqualError(err, "context canceled", "context error is not nil. group: %s, topic: %s, message: %s", groupName, topicName, err.Error())
	}()

	time.Sleep(DefaultConsumeTime)

	cancel()

	err = ctx.Err()
	asst.EqualError(err, "context canceled", "context error is not nil. group: %s, topic: %s, message: %s", groupName, topicName, err.Error())
}
