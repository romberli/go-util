package kafka

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/romberli/log"
	"github.com/stretchr/testify/assert"

	"github.com/romberli/go-util/constant"
	"github.com/romberli/go-util/linux"
)

const (
	DefaultProduceTime = 10 * time.Second
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

	log.SetDisableEscape(true)
	log.SetDisableDoubleQuotes(true)

	kafkaVersion := "3.8.0"
	brokerList := []string{"192.168.137.11:9092"}
	topicName := "test001"

	p, err = NewAsyncProducer(kafkaVersion, brokerList)
	asst.Nil(err, "create producer failed.")

	topics, err := p.Client.Topics()
	asst.Nil(err, "get topics failed.")
	for _, topic := range topics {
		t.Logf("topic: %s", topic)
	}

	clusterNameHeader := p.BuildProducerMessageHeader("clusterName", "main01")
	hostIP, err = linux.GetDefaultIP()
	if err != nil {
		asst.Nil(err, fmt.Sprintf("get host ip failed. message: %s", err.Error()))
	}
	hostIP = "192.168.137.11"
	addrHeader := p.BuildProducerMessageHeader("addr", fmt.Sprintf("%s:%d", hostIP, 3306))
	headers = append(headers, clusterNameHeader, addrHeader)
	ctx, cancel := context.WithCancel(context.Background())

	defer func() {
		err = p.Close()
		if err != nil {
			log.Errorf("close producer failed. topic: %s, message: %s", topicName, err.Error())
		}
	}()

	var wg sync.WaitGroup
	go func() {
		wg.Add(constant.OneInt)
		for i := 0; i < 10; i++ {
			ts = time.Now().Format(constant.DefaultTimeLayout)
			msg := fmt.Sprintf(`{"id": "%d", "message": "hello, world!", "timestamp": %s}`, i, ts)

			if i%2 == 0 {
				err = p.Produce(topicName, msg)
				asst.Nil(err, "produce string message failed. topic: %s, message: %s", topicName, ts)
			} else {
				message = p.BuildProducerMessage(topicName, strconv.Itoa(i), msg, headers)
				err = p.Produce(topicName, message)
				asst.Nil(err, "produce producer message failed. topic: %s, message: %v", topicName, message)
			}

			err = ctx.Err()
			asst.Nil(err, "context error is not nil. topic: %s, errMessage: %v", topicName, err)
		}

		wg.Done()
	}()

	wg.Wait()
	time.Sleep(DefaultProduceTime)

	cancel()

	err = ctx.Err()
	asst.EqualError(err, "context canceled", "context error is not nil. topic: %s", topicName)
}
