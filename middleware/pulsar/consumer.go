package pulsar

import (
	"context"
	"time"

	"github.com/pingcap/errors"

	"github.com/apache/pulsar-client-go/pulsar"

	"github.com/romberli/go-util/common"
	"github.com/romberli/go-util/constant"
)

const (
	DefaultConsumerNamePrefix                  = "pulsar-consumer-"
	SubscriptionTypeExclusive SubscriptionType = iota
	SubscriptionTypeShared
	SubscriptionTypeFailover
	SubscriptionTypeKeyShared

	SubscriptionTypeExclusiveName = "Exclusive"
	SubscriptionTypeSharedName    = "Shared"
	SubscriptionTypeFailoverName  = "Failover"
	SubscriptionTypeKeySharedName = "KeyShared"
	SubscriptionTypeUnknownName   = "Unknown"
)

type SubscriptionType int

type Consumer struct {
	Config           *Config
	SubscriptionName string
	SubscriptionType SubscriptionType
	Name             string
	Client           pulsar.Client
	Consumer         pulsar.Consumer
}

func subscriptionTypeToPulsar(t SubscriptionType) pulsar.SubscriptionType {
	switch t {
	case SubscriptionTypeExclusive:
		return pulsar.Exclusive
	case SubscriptionTypeShared:
		return pulsar.Shared
	case SubscriptionTypeFailover:
		return pulsar.Failover
	case SubscriptionTypeKeyShared:
		return pulsar.KeyShared
	default:
		return pulsar.Exclusive
	}
}

func (t SubscriptionType) String() string {
	switch t {
	case SubscriptionTypeExclusive:
		return SubscriptionTypeExclusiveName
	case SubscriptionTypeShared:
		return SubscriptionTypeSharedName
	case SubscriptionTypeFailover:
		return SubscriptionTypeFailoverName
	case SubscriptionTypeKeyShared:
		return SubscriptionTypeKeySharedName
	default:
		return SubscriptionTypeUnknownName
	}
}

func NewConsumer(config *Config, subscriptionName string, subscriptionType SubscriptionType, name string) (*Consumer, error) {
	clientOpts := pulsar.ClientOptions{
		URL: config.URL,
	}
	if config.Token != constant.EmptyString {
		clientOpts.Authentication = pulsar.NewAuthenticationToken(config.Token)
	}

	client, err := pulsar.NewClient(clientOpts)
	if err != nil {
		return nil, errors.Trace(err)
	}

	consumerOpts := pulsar.ConsumerOptions{
		Topic:            config.Topic,
		SubscriptionName: subscriptionName,
		Type:             subscriptionTypeToPulsar(subscriptionType),
	}
	if name != "" {
		consumerOpts.Name = name
	}

	consumer, err := client.Subscribe(consumerOpts)
	if err != nil {
		client.Close()
		return nil, errors.Trace(err)
	}

	return &Consumer{
		Config:           config,
		SubscriptionName: subscriptionName,
		SubscriptionType: subscriptionType,
		Name:             consumer.Name(),
		Client:           client,
		Consumer:         consumer,
	}, nil
}

func NewConsumerWithDefault(config *Config, subscriptionName string, subscriptionType SubscriptionType) (*Consumer, error) {
	name := DefaultConsumerNamePrefix + common.GetRandomString(common.DefaultNormalCharString+common.DefaultDigitalCharString, 6)
	return NewConsumer(config, subscriptionName, subscriptionType, name)
}

func (c *Consumer) Close() {
	if c.Consumer != nil {
		c.Consumer.Close()
	}
}

func (c *Consumer) Disconnect() {
	c.Close()
	if c.Client != nil {
		c.Client.Close()
	}
}

func (c *Consumer) Receive(ctx context.Context) (pulsar.Message, error) {
	msg, err := c.Consumer.Receive(ctx)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return msg, nil
}

func (c *Consumer) Chan() <-chan pulsar.ConsumerMessage {
	return c.Consumer.Chan()
}

func (c *Consumer) Ack(msg pulsar.Message) error {
	return errors.Trace(c.Consumer.Ack(msg))
}

func (c *Consumer) Nack(msg pulsar.Message) {
	c.Consumer.Nack(msg)
}

func (c *Consumer) Seek(msgID pulsar.MessageID) error {
	return errors.Trace(c.Consumer.Seek(msgID))
}

func (c *Consumer) SeekByTime(timestamp time.Time) error {
	return errors.Trace(c.Consumer.SeekByTime(timestamp))
}

func (c *Consumer) IsConnected() bool {
	return c.Client != nil && c.Consumer != nil
}
