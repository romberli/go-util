package kafka

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/romber2001/log"
)

type DefaultConsumerGroupHandler struct{}

func (DefaultConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (DefaultConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (h DefaultConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		fmt.Printf("topic: %s, partition: %d, offset: %d, key: %s, value: %s", message.Topic, message.Partition, message.Offset, string(message.Key), string(message.Value))
		sess.MarkMessage(message, "")
	}

	return nil
}

type ConsumerGroup struct {
	Ctx          context.Context
	KafkaVersion sarama.KafkaVersion
	BrokerList   []string
	GroupName    string
	Config       *sarama.Config
	Client       sarama.Client
	Group        sarama.ConsumerGroup
}

func NewConsumerGroup(ctx context.Context, kafkaVersion string, brokerList []string, groupName string) (cg *ConsumerGroup, err error) {
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
	defer func() { _ = client.Close() }()

	// Start a new consumer group
	group, err := sarama.NewConsumerGroupFromClient(groupName, client)
	if err != nil {
		return nil, err
	}
	defer func() { _ = group.Close() }()

	return &ConsumerGroup{
		Ctx:          ctx,
		KafkaVersion: config.Version,
		BrokerList:   brokerList,
		GroupName:    groupName,
		Config:       config,
		Client:       client,
		Group:        group,
	}, nil
}

func (cg *ConsumerGroup) Consume(topicName string, handler sarama.ConsumerGroupHandler) (err error) {
	// Track errors
	go func() {
		for err = range cg.Group.Errors() {
			log.Errorf("got error when consuming topic. group: %s, topic: %s, message: %s",
				cg.GroupName, topicName, err.Error())
		}
	}()

	// Iterate over consumer sessions.
	for {
		select {
		case <-cg.Ctx.Done():
			return cg.Ctx.Err()
		default:
			topics := []string{topicName}

			err = cg.Group.Consume(cg.Ctx, topics, handler)
			if err != nil {
				return err
			}
		}
	}
}
