package kafka

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/Shopify/sarama"
	"github.com/romberli/log"
	"github.com/stretchr/testify/assert"

	"github.com/romberli/go-util/linux"
)

const (
	DefaultProduceTime = 5 * time.Second
)

func TestProduce(t *testing.T) {
	var (
		err     error
		p       *AsyncProducer
		ts      string
		hostIP  string
		headers []sarama.RecordHeader
		message *sarama.ProducerMessage
	)

	asst := assert.New(t)

	kafkaVersion := "2.2.0"
	brokerList := []string{"10.0.0.63:9092", "10.0.0.84:9092", "10.0.0.92:9092"}
	topicName := "test001"
	// clusterNameHeader := p.BuildProducerMessageHeader("clusterName", "main01")
	clusterNameHeader := p.BuildProducerMessageHeader("clusterName", "main01")
	hostIP, err = linux.GetDefaultIP()
	if err != nil {
		asst.Nil(err, fmt.Sprintf("get host ip failed. message: %s", err.Error()))
	}
	addrHeader := p.BuildProducerMessageHeader("addr", fmt.Sprintf("%s:%d", hostIP, 3306))
	headers = append(headers, clusterNameHeader, addrHeader)
	ctx, cancel := context.WithCancel(context.Background())

	p, err = NewAsyncProducer(kafkaVersion, brokerList)
	asst.Nil(err, "create producer failed.")

	defer func() {
		err = p.Close()
		if err != nil {
			log.Errorf("close producer failed. topic: %s, message: %s", topicName, err.Error())
		}
	}()

	go func() {
		for i := 0; i < 10; i++ {
			ts = time.Now().String()

			if i%2 == 0 {
				err = p.Produce(topicName, ts)
				asst.Nil(err, "produce string message failed. topic: %s, message: %s", topicName, ts)
			} else {
				message = p.BuildProducerMessage(topicName, strconv.Itoa(i), ts, headers)
				err = p.Produce(topicName, message)
				asst.Nil(err, "produce producer message failed. topic: %s, message: %v", topicName, message)
			}

			err = ctx.Err()
			asst.Nil(err, "context error is not nil. topic: %s, errMessage: %v", topicName, err)
		}
	}()

	time.Sleep(DefaultProduceTime)

	cancel()

	err = ctx.Err()
	asst.EqualError(err, "context canceled", "context error is not nil. topic: %s", topicName)
}
