package consumer

import (
	"fmt"

	"github.com/pingcap/errors"

	"github.com/romberli/go-util/constant"
	"github.com/romberli/go-util/middleware/rabbitmq"
	"github.com/romberli/go-util/middleware/rabbitmq/client"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	ErrQueueExclusiveUseCode           = amqp.AccessRefused
	ErrQueueExclusiveUseReasonTemplate = `ACCESS_REFUSED - queue '%s' in vhost '%s' in exclusive use`
	ErrChannelOrConnectionClosedCode   = amqp.ChannelError
	ErrChannelOrConnectionClosedReason = `channel/connection is not open`
)

type Consumer struct {
	Conn  *client.Conn
	Chan  *amqp.Channel
	Queue amqp.Queue
}

// NewConsumer returns a new *Consumer
func NewConsumer(addr, user, pass, vhost, tag, exchange, queue, key string) (*Consumer, error) {
	return NewConsumerWithConfig(rabbitmq.NewConfig(addr, user, pass, vhost, tag, exchange, queue, key))
}

// NewConsumerWithConfig returns a new *Consumer with given config
func NewConsumerWithConfig(config *rabbitmq.Config) (*Consumer, error) {
	conn, err := client.NewConnWithConfig(config)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Connection.Channel()
	if err != nil {
		return nil, errors.Trace(err)
	}

	return &Consumer{
		Conn: conn,
		Chan: channel,
	}, nil
}

// Close disconnects the rabbitmq server
func (c *Consumer) Close() error {
	if c.Chan != nil && !c.Chan.IsClosed() {
		return errors.Trace(c.Chan.Close())
	}

	return nil
}

// Disconnect disconnects the rabbitmq server
func (c *Consumer) Disconnect() error {
	err := c.Close()
	if err != nil {
		return err
	}

	return c.Conn.Close()
}

// Channel returns the amqp channel,
// if the channel of the consumer is nil or had been closed, a new channel will be opened,
// otherwise the existing channel will be returned
func (c *Consumer) Channel() (*amqp.Channel, error) {
	if c.Chan == nil || c.Chan.IsClosed() {
		var err error
		c.Chan, err = c.Conn.Channel()
		if err != nil {
			return nil, err
		}
	}

	return c.Chan, nil
}

// ExchangeDeclare declares an exchange
func (c *Consumer) ExchangeDeclare(name, kind string) error {
	channel, err := c.Channel()
	if err != nil {
		return err
	}

	return errors.Trace(channel.ExchangeDeclare(name, kind, true, false, false, false, nil))
}

// QueueDeclare declares a queue
func (c *Consumer) QueueDeclare(name string) error {
	channel, err := c.Channel()
	if err != nil {
		return err
	}

	queue, err := channel.QueueDeclare(name, true, false, false, false, nil)
	if err != nil {
		return errors.Trace(err)
	}

	c.Queue = queue

	return nil
}

// QueueBind binds a queue to an exchange
func (c *Consumer) QueueBind(queue, exchange, key string) error {
	channel, err := c.Channel()
	if err != nil {
		return err
	}

	return errors.Trace(channel.QueueBind(queue, key, exchange, false, nil))
}

// Qos controls how many messages or how many bytes the server will try to keep on
// the network for consumers before receiving delivery acks.
func (c *Consumer) Qos(prefetchCount int, global bool) error {
	channel, err := c.Channel()
	if err != nil {
		return err
	}

	return errors.Trace(channel.Qos(prefetchCount, constant.ZeroInt, global))
}

// Consume consumes messages from the queue
func (c *Consumer) Consume(queue string, exclusive bool) (<-chan amqp.Delivery, error) {
	channel, err := c.Channel()
	if err != nil {
		return nil, err
	}

	deliveryChan, err := channel.Consume(queue, c.Conn.Config.Tag, false, exclusive, false, false, nil)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return deliveryChan, nil
}

// Cancel cancels the delivery of a consumer
func (c *Consumer) Cancel() error {
	channel, err := c.Channel()
	if err != nil {
		return err
	}

	return errors.Trace(channel.Cancel(c.Conn.Config.Tag, false))
}

// Ack acknowledges a delivery
func (c *Consumer) Ack(tag uint64, multiple bool) error {
	channel, err := c.Channel()
	if err != nil {
		return err
	}

	return errors.Trace(channel.Ack(tag, multiple))
}

// Nack negatively acknowledge a delivery
func (c *Consumer) Nack(tag uint64, multiple bool, requeue bool) error {
	channel, err := c.Channel()
	if err != nil {
		return err
	}

	return errors.Trace(channel.Nack(tag, multiple, requeue))
}

// IsExclusiveUseError returns true if the error is exclusive use error
func (c *Consumer) IsExclusiveUseError(queue string, err error) bool {
	if errors.HasStack(err) {
		err = errors.Unwrap(err)
	}

	e, ok := err.(*amqp.Error)
	if !ok {
		return false
	}

	message := fmt.Sprintf(ErrQueueExclusiveUseReasonTemplate, queue, c.Conn.Config.Vhost)

	return e.Code == ErrQueueExclusiveUseCode && e.Reason == message
}

// IsNotFoundQueueError returns true if the error is channel or connection closed error
func (c *Consumer) IsChannelOrConnectionClosedError(err error) bool {
	if errors.HasStack(err) {
		err = errors.Unwrap(err)
	}

	e, ok := err.(*amqp.Error)
	if !ok {
		return false
	}

	return e.Code == ErrChannelOrConnectionClosedCode && e.Reason == ErrChannelOrConnectionClosedReason
}
