package rabbitmq

import (
	"fmt"

	"github.com/pingcap/errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/romberli/go-util/constant"
)

const (
	ErrQueueExclusiveUseCode           = amqp.AccessRefused
	ErrQueueExclusiveUseReasonTemplate = `ACCESS_REFUSED - queue '%s' in vhost '%s' in exclusive use`
)

type Consumer struct {
	Conn  *Conn
	Chan  *amqp.Channel
	Queue amqp.Queue
}

// NewConsumer returns a new *Consumer
func NewConsumer(addr, user, pass, vhost, tag string) (*Consumer, error) {
	return NewConsumerWithConfig(NewConfig(addr, user, pass, vhost, tag))
}

// NewConsumerWithDefault returns a new *Consumer with default config
func NewConsumerWithDefault(addr, user, pass string) (*Consumer, error) {
	return NewConsumerWithConfig(NewConfigWithDefault(addr, user, pass))
}

// NewConsumerWithConfig returns a new *Consumer with given config
func NewConsumerWithConfig(config *Config) (*Consumer, error) {
	conn, err := NewConnWithConfig(config)
	if err != nil {
		return nil, err
	}

	return NewConsumerWithConn(conn)
}

// NewConsumerWithConn returns a new *Consumer with given connection
func NewConsumerWithConn(conn *Conn) (*Consumer, error) {
	channel, err := conn.GetConnection().Channel()
	if err != nil {
		return nil, errors.Trace(err)
	}

	return &Consumer{
		Conn: conn,
		Chan: channel,
	}, nil
}

// GetConn returns the connection
func (c *Consumer) GetConn() *Conn {
	return c.Conn
}

// GetChannel returns the channel
func (c *Consumer) GetChannel() *amqp.Channel {
	return c.Chan
}

// GetQueue returns the queue
func (c *Consumer) GetQueue() amqp.Queue {
	return c.Queue
}

// CloseChannel closes the channel
func (c *Consumer) CloseChannel() error {
	return errors.Trace(c.GetChannel().Close())
}

// Close disconnects the rabbitmq server
func (c *Consumer) Close() error {
	if c.GetChannel() != nil && !c.GetChannel().IsClosed() {
		err := c.CloseChannel()
		if err != nil {
			return err
		}
	}

	return c.GetConn().Close()
}

// Channel returns the amqp channel,
// if the channel of the consumer is nil or had been closed, a new channel will be opened,
// otherwise the existing channel will be returned
func (c *Consumer) Channel() (*amqp.Channel, error) {
	if c.GetChannel() == nil || c.GetChannel().IsClosed() {
		var err error
		c.Chan, err = c.GetConn().Channel()
		if err != nil {
			return nil, err
		}
	}

	return c.GetChannel(), nil
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
func (c *Consumer) QueueDeclare(name string) (amqp.Queue, error) {
	channel, err := c.Channel()
	if err != nil {
		return amqp.Queue{}, err
	}

	queue, err := channel.QueueDeclare(name, true, false, false, false, nil)
	if err != nil {
		return amqp.Queue{}, errors.Trace(err)
	}

	return queue, nil
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

	deliveryChan, err := channel.Consume(queue, c.GetConn().GetConfig().GetTag(), false, exclusive, false, false, nil)
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

	return errors.Trace(channel.Cancel(c.GetConn().GetConfig().GetTag(), false))
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

func (c *Consumer) IsExclusiveUseError(queue string, err error) bool {
	if errors.HasStack(err) {
		err = errors.Unwrap(err)
	}

	e, ok := err.(*amqp.Error)
	if !ok {
		return false
	}

	message := fmt.Sprintf(ErrQueueExclusiveUseReasonTemplate, queue, c.GetConn().GetConfig().GetVhost())

	return e.Code == ErrQueueExclusiveUseCode && e.Reason == message
}
