package rabbitmq

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/pingcap/errors"

	"github.com/romberli/go-util/constant"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	DefaultVhost = "/"

	ExchangeTyeDirect   = "direct"
	ExchangeTypeFanout  = "fanout"
	ExchangeTypeTopic   = "topic"
	ExchangeTypeHeaders = "headers"
	DefaultExchangeType = ExchangeTypeTopic
)

type Config struct {
	Addr  string
	User  string
	Pass  string
	Vhost string
	Tag   string
}

// NewConfig returns a new *Config
func NewConfig(addr, user, pass, vhost, tag string) *Config {
	return newConfig(addr, user, pass, vhost, tag)
}

// NewConfigWithDefault returns a new *Config with default values
func NewConfigWithDefault(addr, user, pass string) *Config {
	return NewConfig(addr, user, pass, DefaultVhost, uuid.New().String())
}

// newConfig returns a new *Config
func newConfig(addr, user, pass, vhost, tag string) *Config {
	if tag == constant.EmptyString {
		tag = uuid.New().String()
	}

	return &Config{
		Addr:  addr,
		User:  user,
		Pass:  pass,
		Vhost: vhost,
		Tag:   tag,
	}
}

// GetAddr returns the address
func (c *Config) GetAddr() string {
	return c.Addr
}

// GetUser returns the user
func (c *Config) GetUser() string {
	return c.User
}

// GetPass returns the password
func (c *Config) GetPass() string {
	return c.Pass
}

// GetVhost returns the vhost
func (c *Config) GetVhost() string {
	return c.Vhost
}

// GetTag returns the tag
func (c *Config) GetTag() string {
	return c.Tag
}

// GetURL returns the URL
func (c *Config) GetURL() string {
	return fmt.Sprintf("amqp://%s:%s@%s%s", c.GetUser(), c.GetPass(), c.GetAddr(), c.GetVhost())
}

type Conn struct {
	Config *Config
	*amqp.Connection
}

// NewConn returns a new *Conn
func NewConn(addr, user, pass, vhost, tag string) (*Conn, error) {
	config := NewConfig(addr, user, pass, vhost, tag)

	return NewConnWithConfig(config)
}

// NewConnWithDefault returns a new *Conn with default values
func NewConnWithDefault(addr, user, pass string) (*Conn, error) {
	return NewConnWithConfig(NewConfigWithDefault(addr, user, pass))
}

// NewConnWithConfig returns a new *Conn with given config
func NewConnWithConfig(config *Config) (*Conn, error) {
	conn, err := amqp.Dial(config.GetURL())
	if err != nil {
		return nil, errors.Trace(err)
	}

	return &Conn{
		Config:     config,
		Connection: conn,
	}, nil
}

// GetConfig returns the config
func (c *Conn) GetConfig() *Config {
	return c.Config
}

// GetConnection returns the connection
func (c *Conn) GetConnection() *amqp.Connection {
	return c.Connection
}

// Close closes the connection
func (c *Conn) Close() error {
	return errors.Trace(c.GetConnection().Close())
}
