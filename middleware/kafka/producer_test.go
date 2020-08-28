package kafka

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/romberli/go-util/linux"
	"github.com/romberli/log"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

const (
	DefaultProduceSeconds = 5 * time.Second
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

	assert := assert.New(t)

	kafkaVersion := "2.2.0"
	brokerList := []string{"10.0.0.63:9092", "10.0.0.84:9092", "10.0.0.92:9092"}
	topicName := "test001"
	//clusterNameHeader := p.BuildProducerMessageHeader("clusterName", "main01")
	clusterNameHeader := p.BuildProducerMessageHeader("clusterName", "main01")
	hostIP, err = linux.GetDefaultIP()
	if err != nil {
		assert.Nil(err, fmt.Sprintf("get host ip failed. message: %s", err.Error()))
	}
	addrHeader := p.BuildProducerMessageHeader("addr", fmt.Sprintf("%s:%d", hostIP, 3306))
	headers = append(headers, clusterNameHeader, addrHeader)
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

			if i%2 == 0 {
				err = p.Produce(topicName, ts)
				assert.Nil(err, "produce string message failed. topic: %s, message: %s", topicName, ts)
			} else {
				message = p.BuildProducerMessage(topicName, strconv.Itoa(i), ts, headers)
				err = p.Produce(topicName, message)
				assert.Nil(err, "produce producer message failed. topic: %s, message: %v", topicName, message)
			}

			err = ctx.Err()
			assert.Nil(err, "context error is not nil. topic: %s, errMessage: %v", topicName, err)
		}
	}()

	time.Sleep(DefaultProduceSeconds)

	cancel()

	err = ctx.Err()
	assert.EqualError(err, "context canceled", "context error is not nil. topic: %s", topicName)
}
