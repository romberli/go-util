package kafka

import (
	"context"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"testing"
	"time"
)

const (
	DefaultConsumeSeconds = 10
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
	startTime := time.Now().Unix()

	cg, err = NewConsumerGroup(ctx, kafkaVersion, brokerList, groupName)
	assert.Nil(err, "create consumer group failed.")

	wg := &sync.WaitGroup{}
	go func() {
		wg.Add(1)
		defer wg.Done()

		for {
			if time.Now().Unix() >= startTime+DefaultConsumeSeconds {
				cancel()
			}

			err = cg.Consume(topicName, handler)
			assert.Nil(err, "consume failed. group: %s, topic: %s")

			err = cg.Ctx.Err()
			assert.Nil(err, "context error is not nil. group: %s, topic: %s")
		}
	}()

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-ctx.Done():
		log.Println("terminating: context cancelled")
	case <-sigterm:
		log.Println("terminating: via signal")
	}
	cancel()
	wg.Wait()
}
