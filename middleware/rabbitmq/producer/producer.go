package producer

import (
	"context"
	"strconv"

	"github.com/pingcap/errors"

	"github.com/romberli/go-util/constant"
	"github.com/romberli/go-util/middleware/rabbitmq"
	"github.com/romberli/go-util/middleware/rabbitmq/client"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	Conn  *client.Conn
	Chan  *amqp.Channel
	Queue amqp.Queue
}

// NewProducer returns a new *Producer
func NewProducer(addr, user, pass, vhost, tag, exchange, queue, key string) (*Producer, error) {
	return NewProducerWithConfig(rabbitmq.NewConfig(addr, user, pass, vhost, tag, exchange, queue, key))
}

// NewProducerWithConfig returns a new *Producer with given config
func NewProducerWithConfig(config *rabbitmq.Config) (*Producer, error) {
	conn, err := client.NewConnWithConfig(config)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Connection.Channel()
	if err != nil {
		return nil, errors.Trace(err)
	}

	return &Producer{
		Conn: conn,
		Chan: channel,
	}, nil
}

// Close closes the channel
func (p *Producer) Close() error {
	if p.Chan != nil && !p.Chan.IsClosed() {
		return errors.Trace(p.Chan.Close())
	}

	return nil
}

// Disconnect disconnects the rabbitmq server
func (p *Producer) Disconnect() error {
	err := p.Close()
	if err != nil {
		return err
	}

	return p.Conn.Close()
}

// Channel returns the channel, it the channel of the producer is nil or had been closed, a new channel will be opened
func (p *Producer) Channel() (*amqp.Channel, error) {
	if p.Chan == nil || p.Chan.IsClosed() {
		channel, err := p.Conn.Connection.Channel()
		if err != nil {
			return nil, errors.Trace(err)
		}

		p.Chan = channel
	}

	return p.Chan, nil
}

// ExchangeDeclare declares an exchange
func (p *Producer) ExchangeDeclare(kind string) error {
	channel, err := p.Channel()
	if err != nil {
		return err
	}

	err = channel.ExchangeDeclare(p.Conn.Config.Exchange, kind, true, false, false, false, nil)
	if err != nil {
		return errors.Trace(err)
	}

	return nil
}

// QueueDeclare declares a queue
func (p *Producer) QueueDeclare() error {
	channel, err := p.Channel()
	if err != nil {
		return err
	}

	queue, err := channel.QueueDeclare(p.Conn.Config.Queue, true, false, false, false, nil)
	if err != nil {
		return errors.Trace(err)
	}

	p.Queue = queue

	return nil
}

// QueueBind binds a queue to an exchange
func (p *Producer) QueueBind() error {
	channel, err := p.Channel()
	if err != nil {
		return err
	}

	err = channel.QueueBind(p.Conn.Config.Queue, p.Conn.Config.Key, p.Conn.Config.Exchange, false, nil)
	if err != nil {
		return errors.Trace(err)
	}

	return nil
}

// BuildMessage builds an amqp.Publishing with given content type and message
func (p *Producer) BuildMessage(contentType, message string) amqp.Publishing {
	return amqp.Publishing{
		ContentType: contentType,
		Body:        []byte(message),
	}
}

// BuildMessageWithExpiration builds an amqp.Publishing with given content type and message and expiration
func (p *Producer) BuildMessageWithExpiration(contentType, message string, expiration int) amqp.Publishing {
	return amqp.Publishing{
		ContentType: contentType,
		Body:        []byte(message),
		Expiration:  strconv.Itoa(expiration),
	}
}

// PublishMessage publishes a json message to an exchange
func (p *Producer) PublishJSON(message string) error {
	return p.Publish(p.BuildMessage(constant.DefaultJSONContentType, message))
}

// Publish publishes a message to an exchange
func (p *Producer) Publish(msg amqp.Publishing) error {
	return p.publishWithContext(context.Background(), msg)
}

// PublishWithContext publishes a message to an exchange with context
func (p *Producer) PublishWithContext(ctx context.Context, msg amqp.Publishing) error {
	return p.publishWithContext(ctx, msg)
}

// Publish publishes a message to an exchange
func (p *Producer) publishWithContext(ctx context.Context, msg amqp.Publishing) error {
	channel, err := p.Channel()
	if err != nil {
		return err
	}

	return errors.Trace(channel.PublishWithContext(ctx, p.Conn.Config.Exchange, p.Conn.Config.Key, false, false, msg))
}
