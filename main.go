package main

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/romber2001/go-util/middleware/kafka"
	"github.com/romber2001/log"
	"time"
)

const (
	DefaultProduceSeconds = 5 * time.Second
	DefaultConsumeSeconds = 5 * time.Second
)

func main() {
	var (
		err error
		ts  string
		p   *kafka.AsyncProducer
		cg  *kafka.ConsumerGroup
	)

	kafkaVersion := "2.2.0"
	brokerList := []string{"10.0.0.63:9092", "10.0.0.84:9092", "10.0.0.92:9092"}
	groupName := "group001"
	topicName := "test001"
	ctx, cancel := context.WithCancel(context.Background())
	handler := kafka.DefaultConsumerGroupHandler{}

	p, err = kafka.NewAsyncProducer(ctx, kafkaVersion, brokerList)
	if err != nil {
		log.Errorf("create producer failed. topic: %s, errMessage: %s", topicName, err.Error())
	}

	go func() {
		for i := 0; i < 10; i++ {
			ts = time.Now().String()
			//ts = strconv.Itoa(i)
			log.Infof("message: %s", ts)
			err = p.Produce(topicName, ts)
			if err != nil {
				log.Errorf("produce message failed. topic: %s, message: %s, errMessage: %s", topicName, ts, err.Error())
			}

			err = p.Ctx.Err()
			if err != nil {
				log.Errorf("produce message failed. topic: %s, message: %s, errMessage: %s", topicName, ts, err.Error())
			}

		}
	}()

	time.Sleep(DefaultProduceSeconds * 1)

	cancel()

	err = p.Ctx.Err()
	if err != nil {
		log.Errorf("context error not nil. group: %s, topic: %s, errMessage: %s", groupName, topicName, err.Error())
	}

	ctx, cancel = context.WithCancel(context.Background())
	cg, err = kafka.NewConsumerGroup(ctx, kafkaVersion, brokerList, groupName, sarama.OffsetNewest)
	if err != nil {
		log.Errorf("create consumer group failed. group: %s, topic: %s, errMessage: %s", groupName, topicName, err.Error())
	}

	go func() {
		err = cg.Consume(topicName, handler)
		if err != nil {
			log.Errorf("consume topic failed. group: %s, topic: %s, errMessage: %s", groupName, topicName, err.Error())
		}

		err = cg.Ctx.Err()
		if err != nil {
			log.Errorf("context error not nil. group: %s, topic: %s, errMessage: %s", groupName, topicName, err.Error())
		}
	}()

	time.Sleep(DefaultConsumeSeconds * 3)

	cancel()

	err = cg.Ctx.Err()
	if err != nil {
		log.Errorf("context error not nil. group: %s, topic: %s, errMessage: %s", groupName, topicName, err.Error())
	}
}
