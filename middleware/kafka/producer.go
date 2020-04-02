package kafka

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/romber2001/log"
)

type AsyncProducer struct {
	Ctx          context.Context
	KafkaVersion sarama.KafkaVersion
	BrokerList   []string
	TopicName    string
	Config       *sarama.Config
	Client       sarama.Client
	Producer     sarama.AsyncProducer
}

func NewAsyncProducer(ctx context.Context, kafkaVersion string, brokerList []string, topicName string) (p *AsyncProducer, err error) {
	// Init config, specify appropriate version
	config := sarama.NewConfig()
	config.Version, err = sarama.ParseKafkaVersion(kafkaVersion)
	if err != nil {
		return nil, err
	}

	// Start with a client
	client, err := sarama.NewClient(brokerList, config)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = client.Close()
		log.Errorf("got error when closing client. topic: %s, message: %s",
			topicName, err.Error())
	}()

	// Start a new consumer group
	producer, err := sarama.NewAsyncProducerFromClient(client)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = producer.Close()
		log.Errorf("got error when closing producer. topic: %s, message: %s",
			topicName, err.Error())
	}()

	return &AsyncProducer{
		Ctx:          ctx,
		KafkaVersion: config.Version,
		BrokerList:   brokerList,
		TopicName:    topicName,
		Config:       config,
		Client:       client,
		Producer:     producer,
	}, nil
}

func (p *AsyncProducer) Produce(message string) (err error) {
	defer func() {
		err = p.Producer.Close()
		log.Errorf("got error when closing producer. group: %s, topic: %s, message: %s",
			p.TopicName, err.Error())
	}()

	// Produce message to kafka
	producerMessage := &sarama.ProducerMessage{Topic: p.TopicName, Key: nil, Value: sarama.StringEncoder(message), Metadata: 0}
	p.Producer.Input() <- producerMessage

	return nil
}
