package pulsar

import (
	"context"
	"encoding/json"

	"github.com/pingcap/errors"

	"github.com/apache/pulsar-client-go/pulsar"

	"github.com/romberli/go-util/common"
	"github.com/romberli/go-util/constant"
)

const (
	DefaultProducerNamePrefix = "pulsar-producer-"
)

type Producer struct {
	Config   *Config
	Name     string
	Client   pulsar.Client
	Producer pulsar.Producer
}

func NewProducer(config *Config, name string) (*Producer, error) {
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

	producerOpts := pulsar.ProducerOptions{
		Topic: config.Topic,
	}
	if name != constant.EmptyString {
		producerOpts.Name = name
	}

	producer, err := client.CreateProducer(producerOpts)
	if err != nil {
		client.Close()
		return nil, errors.Trace(err)
	}

	return &Producer{
		Config:   config,
		Name:     producer.Name(),
		Client:   client,
		Producer: producer,
	}, nil
}

func NewProducerWithDefault(config *Config) (*Producer, error) {
	name := DefaultProducerNamePrefix + common.GetRandomString(common.DefaultNormalCharString+common.DefaultDigitalCharString, 6)
	return NewProducer(config, name)
}

func (p *Producer) Close() {
	if p.Producer != nil {
		p.Producer.Close()
	}
}

func (p *Producer) Disconnect() {
	p.Close()
	if p.Client != nil {
		p.Client.Close()
	}
}

func (p *Producer) Send(ctx context.Context, msg *pulsar.ProducerMessage) (pulsar.MessageID, error) {
	msgID, err := p.Producer.Send(ctx, msg)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return msgID, nil
}

func (p *Producer) SendAsync(ctx context.Context, msg *pulsar.ProducerMessage, callback func(pulsar.MessageID, *pulsar.ProducerMessage, error)) {
	p.Producer.SendAsync(ctx, msg, callback)
}

func (p *Producer) SendJSON(ctx context.Context, data interface{}) (pulsar.MessageID, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, errors.Trace(err)
	}

	msg := &pulsar.ProducerMessage{
		Payload: jsonBytes,
	}

	return p.Send(ctx, msg)
}

func (p *Producer) IsConnected() bool {
	return p.Client != nil && p.Producer != nil
}
