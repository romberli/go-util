package rabbitmq

import (
	"strconv"

	"github.com/pingcap/errors"
	"golang.org/x/net/context"

	"github.com/romberli/go-util/constant"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	Conn     *Conn
	Channel  *amqp.Channel
	exchange string
	queue    string
	key      string
}

// NewProducer returns a new *Producer
func NewProducer(addr, user, pass, vhost, tag string) (*Producer, error) {
	return NewProducerWithConfig(NewConfig(addr, user, pass, vhost, tag))
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

// SetExchange sets the exchange
func (p *Producer) SetExchange(exchange string) {
	p.exchange = exchange
}

// SetQueue sets the queue
func (p *Producer) SetQueue(queue string) {
	p.queue = queue
}

// SetKey sets the key
func (p *Producer) SetKey(key string) {
	p.key = key
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
	err := p.GetChannel().ExchangeDeclare(name, kind, true, false, false, false, nil)
	if err != nil {
		return errors.Trace(err)
	}

	p.SetExchange(name)

	return nil
}

// QueueDeclare declares a queue
func (p *Producer) QueueDeclare(name string) (amqp.Queue, error) {
	queue, err := p.GetChannel().QueueDeclare(name, true, false, false, false, nil)
	if err != nil {
		return amqp.Queue{}, errors.Trace(err)
	}

	p.SetQueue(queue.Name)

	return queue, nil
}

// QueueBind binds a queue to an exchange
func (p *Producer) QueueBind(queue, exchange, key string) error {
	err := p.GetChannel().QueueBind(queue, key, exchange, false, nil)
	if err != nil {
		return errors.Trace(err)
	}

	p.SetKey(key)

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
	return p.publishWithContext(context.Background(), p.exchange, p.key, p.BuildMessage(constant.DefaultJSONContentType, message))
}

// Publish publishes a message to an exchange
func (p *Producer) Publish(exchange, key string, msg amqp.Publishing) error {
	return p.publishWithContext(context.Background(), exchange, key, msg)
}

// PublishWithContext publishes a message to an exchange with context
func (p *Producer) PublishWithContext(ctx context.Context, exchange, key string, msg amqp.Publishing) error {
	return p.publishWithContext(ctx, exchange, key, msg)
}

// Publish publishes a message to an exchange
func (p *Producer) publishWithContext(ctx context.Context, exchange, key string, msg amqp.Publishing) error {
	return errors.Trace(p.GetChannel().PublishWithContext(ctx, exchange, key, false, false, msg))
}
