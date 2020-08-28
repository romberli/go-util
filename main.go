package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/Shopify/sarama"
	"github.com/romberli/log"

	"github.com/romberli/go-util/linux"
	"github.com/romberli/go-util/middleware/kafka"
)

const (
	DefaultProduceSeconds = 5 * time.Second
	DefaultConsumeSeconds = 5 * time.Second
)

func main() {
	var (
		err     error
		ts      string
		hostIP  string
		headers []sarama.RecordHeader
		message *sarama.ProducerMessage
		p       *kafka.AsyncProducer
		cg      *kafka.ConsumerGroup
	)

	kafkaVersion := "2.2.0"
	brokerList := []string{"10.0.0.63:9092", "10.0.0.84:9092", "10.0.0.92:9092"}
	groupName := "group001"
	topicName := "test001"
	ctx, cancel := context.WithCancel(context.Background())
	clusterNameHeader := p.BuildProducerMessageHeader("clusterName", "main01")
	hostIP, err = linux.GetDefaultIP()
	if err != nil {
		log.Errorf("get host ip failed. message: %s", err.Error())
	}
	addrHeader := p.BuildProducerMessageHeader("addr", fmt.Sprintf("%s:%d", hostIP, 3306))
	headers = append(headers, clusterNameHeader, addrHeader)
	handler := kafka.DefaultConsumerGroupHandler{}

	p, err = kafka.NewAsyncProducer(kafkaVersion, brokerList)
	if err != nil {
		log.Errorf("create producer failed. topic: %s, errMessage: %s", topicName, err.Error())
	}

	go func() {
		for i := 0; i < 10; i++ {
			ts = time.Now().String()
			log.Infof("message: %s", ts)
			message = p.BuildProducerMessage(topicName, strconv.Itoa(i), ts, headers)
			err = p.Produce(topicName, message)
			if err != nil {
				log.Errorf("produce message failed. topic: %s, message: %s, errMessage: %s", topicName, ts, err.Error())
			}

			err = ctx.Err()
			if err != nil {
				log.Errorf("produce message failed. topic: %s, message: %s, errMessage: %s", topicName, ts, err.Error())
			}
		}
	}()

	time.Sleep(DefaultProduceSeconds * 1)

	cancel()

	err = ctx.Err()
	if err != nil {
		log.Errorf("context error not nil. group: %s, topic: %s, errMessage: %s", groupName, topicName, err.Error())
	}

	ctx, cancel = context.WithCancel(context.Background())
	cg, err = kafka.NewConsumerGroup(kafkaVersion, brokerList, groupName, sarama.OffsetNewest)
	if err != nil {
		log.Errorf("create consumer group failed. group: %s, topic: %s, errMessage: %s", groupName, topicName, err.Error())
	}

	go func() {
		err = cg.Consume(ctx, topicName, handler)
		if err != nil {
			log.Errorf("consume topic failed. group: %s, topic: %s, errMessage: %s", groupName, topicName, err.Error())
		}

		err = ctx.Err()
		if err != nil {
			log.Errorf("context error not nil. group: %s, topic: %s, errMessage: %s", groupName, topicName, err.Error())
		}
	}()

	time.Sleep(DefaultConsumeSeconds * 5)

	cancel()

	err = ctx.Err()
	if err != nil {
		log.Errorf("context error not nil. group: %s, topic: %s, errMessage: %s", groupName, topicName, err.Error())
	}
}
