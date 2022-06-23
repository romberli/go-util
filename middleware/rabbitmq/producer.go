package rabbitmq

import (
	"strconv"

	"github.com/pingcap/errors"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	Conn    *Conn
	Channel *amqp.Channel
}

// NewProducer returns a new *Producer
func NewProducer(addr, user, pass, vhost string) (*Producer, error) {
	return NewProducerWithConfig(NewConfig(addr, user, pass, vhost))
}

// NewProducerWithDefault returns a new *Producer with default config
func NewProducerWithDefault(addr, user, pass string) (*Producer, error) {
	return NewProducerWithConfig(NewConfigWithDefault(addr, user, pass))
}

// NewProducerWithConfig returns a new *Producer with given config
func NewProducerWithConfig(config *Config) (*Producer, error) {
	conn, err := NewConnWithConfig(config)
	if err != nil {
		return nil, err
	}

	return NewProducerWithConn(conn)
}

// NewProducerWithConn returns a new *Producer with given connection
func NewProducerWithConn(conn *Conn) (*Producer, error) {
	channel, err := conn.GetConnection().Channel()
	if err != nil {
		return nil, errors.Trace(err)
	}

	return &Producer{
		Conn:    conn,
		Channel: channel,
	}, nil
}

// GetConn returns the connection
func (p *Producer) GetConn() *Conn {
	return p.Conn
}

// GetChannel returns the channel
func (p *Producer) GetChannel() *amqp.Channel {
	return p.Channel
}

// Close closes the channel
func (p *Producer) Close() error {
	return errors.Trace(p.GetChannel().Close())
}

// Disconnect disconnects the rabbitmq server
func (p *Producer) Disconnect() error {
	err := p.Close()
	if err != nil {
		return err
	}

	return p.Conn.Close()
}

// ExchangeDeclare declares an exchange
func (p *Producer) ExchangeDeclare(name, kind string) error {
	return errors.Trace(p.GetChannel().ExchangeDeclare(name, kind, true, false, false, false, nil))
}

// QueueDeclare declares a queue
func (p *Producer) QueueDeclare(name string) (amqp.Queue, error) {
	queue, err := p.GetChannel().QueueDeclare(name, true, false, false, false, nil)
	if err != nil {
		return amqp.Queue{}, errors.Trace(err)
	}

	return queue, nil
}

// QueueBind binds a queue to an exchange
func (p *Producer) QueueBind(queue, exchange, key string) error {
	return errors.Trace(p.GetChannel().QueueBind(queue, key, exchange, false, nil))
}

// BuildMessage builds an amqp.Publishing with given message and content type
func (p *Producer) BuildMessage(message, contentType string) amqp.Publishing {
	return amqp.Publishing{
		ContentType: contentType,
		Body:        []byte(message),
	}
}

// BuildMessageWithExpiration builds an amqp.Publishing with given message and content type and expiration
func (p *Producer) BuildMessageWithExpiration(message, contentType string, expiration int) amqp.Publishing {
	return amqp.Publishing{
		ContentType: contentType,
		Body:        []byte(message),
		Expiration:  strconv.Itoa(expiration),
	}
}

// Publish publishes a message to an exchange
func (p *Producer) Publish(exchange, key string, msg amqp.Publishing) error {
	return errors.Trace(p.GetChannel().Publish(exchange, key, false, false, msg))
}
