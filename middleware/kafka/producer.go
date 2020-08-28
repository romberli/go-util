package kafka

import (
	"errors"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/romberli/log"
	"reflect"
	"time"
)

type AsyncProducer struct {
	KafkaVersion sarama.KafkaVersion
	BrokerList   []string
	Config       *sarama.Config
	Client       sarama.Client
	Producer     sarama.AsyncProducer
}

func NewAsyncProducer(kafkaVersion string, brokerList []string) (p *AsyncProducer, err error) {
	// Init config, specify appropriate version
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Flush.Messages = 1
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	config.Version, err = sarama.ParseKafkaVersion(kafkaVersion)
	if err != nil {
		return nil, err
	}

	// Start with a client
	client, err := sarama.NewClient(brokerList, config)
	if err != nil {
		return nil, err
	}

	// Start a new consumer group
	producer, err := sarama.NewAsyncProducerFromClient(client)
	if err != nil {
		return nil, err
	}

	return &AsyncProducer{
		KafkaVersion: config.Version,
		BrokerList:   brokerList,
		Config:       config,
		Client:       client,
		Producer:     producer,
	}, nil
}

func (p *AsyncProducer) Close() error {
	if p.Producer != nil {
		return p.Producer.Close()
	}

	return nil
}

func (p *AsyncProducer) BuildProducerMessageHeader(key string, value string) sarama.RecordHeader {
	return sarama.RecordHeader{
		Key:   []byte(key),
		Value: []byte(value),
	}
}

func (p *AsyncProducer) BuildProducerMessage(topicName string, key string, message string, headers []sarama.RecordHeader) *sarama.ProducerMessage {
	return &sarama.ProducerMessage{
		Topic:     topicName,
		Key:       sarama.StringEncoder(key),
		Value:     sarama.StringEncoder(message),
		Headers:   headers,
		Metadata:  nil,
		Timestamp: time.Now(),
	}
}

func (p *AsyncProducer) Produce(topicName string, message interface{}) (err error) {
	var (
		producerMessage *sarama.ProducerMessage
	)

	// Track error
	go func() {
		for {
			if p.Producer == nil {
				return
			}

			select {
			case success := <-p.Producer.Successes():
				if success != nil {
					log.Debugf("offset: %d, timestamp: %s, partitions: %d",
						success.Offset, success.Timestamp.String(), success.Partition)
				}
			case fail := <-p.Producer.Errors():
				if fail != nil {
					log.Errorf("err: ", fail.Err)
				}
			}
		}
	}()

	switch message.(type) {
	case string:
		producerMessage = p.BuildProducerMessage(topicName, "", message.(string), nil)
	case *sarama.ProducerMessage:
		producerMessage = message.(*sarama.ProducerMessage)
	default:
		return errors.New(
			fmt.Sprintf("message must be either string type or *sarama.ProducerMessage type, but got %s",
				reflect.TypeOf(message).Name()))
	}

	// Produce message to kafka
	p.Producer.Input() <- producerMessage

	return nil
}
