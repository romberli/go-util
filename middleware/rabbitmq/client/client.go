package client

import (
	"github.com/pingcap/errors"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/romberli/go-util/constant"
	"github.com/romberli/go-util/middleware/rabbitmq"
)

const (
	ExchangeTyeDirect   = "direct"
	ExchangeTypeFanout  = "fanout"
	ExchangeTypeTopic   = "topic"
	ExchangeTypeHeaders = "headers"
	DefaultExchangeType = ExchangeTypeTopic
)

type Conn struct {
	Config *rabbitmq.Config
	*amqp.Connection
}

// NewConn returns a new *Conn
func NewConn(addr, user, pass, vhost string) (*Conn, error) {
	config := rabbitmq.NewConfigWithDefault(addr, user, pass, vhost)

	return NewConnWithConfig(config)
}

// NewConnWithDefault returns a new *Conn with default values
func NewConnWithDefault(addr, user, pass string) (*Conn, error) {
	return NewConnWithConfig(rabbitmq.NewConfigWithDefault(addr, user, pass, constant.DefaultVhost))
}

// NewConnWithConfig returns a new *Conn with given config
func NewConnWithConfig(config *rabbitmq.Config) (*Conn, error) {
	conn, err := amqp.Dial(config.GetURL())
	if err != nil {
		return nil, errors.Trace(err)
	}

	return &Conn{
		Config:     config,
		Connection: conn,
	}, nil
}

// Close closes the connection
func (c *Conn) Close() error {
	return errors.Trace(c.Connection.Close())
}
