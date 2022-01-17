package kafka

import (
	"context"

	"github.com/Shopify/sarama"
	"github.com/pingcap/errors"
	"github.com/romberli/go-util/constant"
	"github.com/romberli/log"
)

type DefaultConsumerGroupHandler struct{}

func (h DefaultConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error { return nil }

func (DefaultConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h DefaultConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		var headers []map[string]string
		for _, header := range message.Headers {
			headers = append(headers, ConvertHeaderToMap(*header))
		}

		log.Infof("topic: %s, partition: %d, offset: %d, key: %s, value: %s, headers: %v",
			message.Topic, message.Partition, message.Offset, string(message.Key), string(message.Value), headers)
		sess.MarkMessage(message, constant.EmptyString)
	}

	return nil
}

type ConsumerGroup struct {
	KafkaVersion sarama.KafkaVersion
	BrokerList   []string
	GroupName    string
	Config       *sarama.Config
	Client       sarama.Client
	Group        sarama.ConsumerGroup
}

func NewConsumerGroup(kafkaVersion string, brokerList []string, groupName string, initOffset int64) (cg *ConsumerGroup, err error) {
	// Init config, specify appropriate version
	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = initOffset
	config.Version, err = sarama.ParseKafkaVersion(kafkaVersion)
	if err != nil {
		return nil, errors.Trace(err)
	}

	// Start with a client
	client, err := sarama.NewClient(brokerList, config)
	if err != nil {
		return nil, errors.Trace(err)
	}

	// Start a new consumer group
	group, err := sarama.NewConsumerGroupFromClient(groupName, client)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return &ConsumerGroup{
		KafkaVersion: config.Version,
		BrokerList:   brokerList,
		GroupName:    groupName,
		Config:       config,
		Client:       client,
		Group:        group,
	}, nil
}

func (cg *ConsumerGroup) Close() error {
	if cg.Group != nil {
		return errors.Trace(cg.Group.Close())
	}

	return nil
}

func (cg *ConsumerGroup) Consume(ctx context.Context, topicName string, handler sarama.ConsumerGroupHandler) (err error) {
	defer func() {
		err = cg.Close()
		if err != nil {
			log.Errorf("close consumer failed. topic: %s, error:\n%+v", topicName, err)
		}
	}()

	// Track errors
	go func() {
		if cg.Group == nil {
			return
		}

		for err = range cg.Group.Errors() {
			log.Errorf("got error when consuming topic. group: %s, topic: %s, error:\n%+v",
				cg.GroupName, topicName, errors.Trace(err))
		}
	}()

	// Iterate over consumer sessions.
	topics := []string{topicName}

	for {
		err = cg.Group.Consume(ctx, topics, handler)
		if err != nil {
			return errors.Trace(err)
		}
	}
}
