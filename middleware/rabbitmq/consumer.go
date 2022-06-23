package rabbitmq

import (
	"github.com/pingcap/errors"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	Conn    *Conn
	Channel *amqp.Channel
	Queue   amqp.Queue
}

// NewConsumer returns a new *Consumer
func NewConsumer(addr, user, pass, vhost string) (*Consumer, error) {
	return NewConsumerWithConfig(NewConfig(addr, user, pass, vhost))
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
		Conn:    conn,
		Channel: channel,
	}, nil
}

// GetConn returns the connection
func (c *Consumer) GetConn() *Conn {
	return c.Conn
}

// GetChannel returns the channel
func (c *Consumer) GetChannel() *amqp.Channel {
	return c.Channel
}

// GetQueue returns the queue
func (c *Consumer) GetQueue() amqp.Queue {
	return c.Queue
}

// Close closes the channel
func (c *Consumer) Close() error {
	return errors.Trace(c.GetChannel().Close())
}

// Disconnect disconnects the rabbitmq server
func (c *Consumer) Disconnect() error {
	err := c.Close()
	if err != nil {
		return err
	}

	return c.Conn.Close()
}

// Consume consumes messages from the queue
func (c *Consumer) Consume(queue, consumer string, exclusive bool) (<-chan amqp.Delivery, error) {
	deliveryChan, err := c.GetChannel().Consume(queue, consumer, false, exclusive, false, false, nil)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return deliveryChan, nil
}

// Qos controls how many messages or how many bytes the server will try to keep on
// the network for consumers before receiving delivery acks.
func (c *Consumer) Qos(prefetchCount, prefetchSize int, global bool) error {
	return errors.Trace(c.GetChannel().Qos(prefetchCount, prefetchSize, global))
}

// Ack acknowledges a delivery
func (c *Consumer) Ack(tag uint64, multiple bool) error {
	return errors.Trace(c.GetChannel().Ack(tag, multiple))
}

// Nack negatively acknowledge a delivery
func (c *Consumer) Nack(tag uint64, multiple bool, requeue bool) error {
	return errors.Trace(c.GetChannel().Nack(tag, multiple, requeue))
}
