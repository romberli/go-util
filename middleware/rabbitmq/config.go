package rabbitmq

import (
	"fmt"

	"github.com/pingcap/errors"

	"github.com/romberli/go-util/constant"
)

const (
	DefaultUnlimitedWaitTime   = -1 // seconds
	DefaultUnlimitedRetryCount = -1
)

type Config struct {
	Addr     string
	User     string
	Pass     string
	Vhost    string
	Tag      string
	Exchange string
	Queue    string
	Key      string
}

// NewConfig returns a new *Config
func NewConfig(addr, user, pass, vhost, tag, exchange, queue, key string) *Config {
	return newConfig(addr, user, pass, vhost, tag, exchange, queue, key)
}

// NewConfigWithDefault returns a new *Config with default values
func NewConfigWithDefault(addr, user, pass, vhost string) *Config {
	return NewConfig(addr, user, pass, vhost, constant.EmptyString,
		constant.EmptyString, constant.EmptyString, constant.EmptyString)
}

// newConfig returns a new *Config
func newConfig(addr, user, pass, vhost, tag, exchange, queue, key string) *Config {
	return &Config{
		Addr:     addr,
		User:     user,
		Pass:     pass,
		Vhost:    vhost,
		Tag:      tag,
		Exchange: exchange,
		Queue:    queue,
		Key:      key,
	}
}

// GetURL returns the URL
func (c *Config) GetURL() string {
	return fmt.Sprintf("amqp://%s:%s@%s%s", c.User, c.Pass, c.Addr, c.Vhost)
}

// Clone returns a new *Config
func (c *Config) Clone() *Config {
	return newConfig(c.Addr, c.User, c.Pass, c.Vhost, c.Tag, c.Exchange, c.Queue, c.Key)
}

type PoolConfig struct {
	*Config
	MaxConnections     int
	InitConnections    int
	MaxIdleConnections int
	MaxIdleTime        int
	MaxWaitTime        int
	MaxRetryCount      int
	KeepAliveInterval  int
}

// NewPoolConfig returns a new PoolConfig
func NewPoolConfig(addr, user, host, vhost, tag, exchange, queue, key string,
	maxConnections, initConnections, maxIdleConnections, maxIdleTime, maxWaitTime, maxRetryCount, keepAliveInterval int) *PoolConfig {
	config := NewConfig(addr, user, host, vhost, tag, exchange, queue, key)

	return &PoolConfig{
		Config:             config,
		MaxConnections:     maxConnections,
		InitConnections:    initConnections,
		MaxIdleConnections: maxIdleConnections,
		MaxIdleTime:        maxIdleTime,
		MaxWaitTime:        maxWaitTime,
		MaxRetryCount:      maxRetryCount,
		KeepAliveInterval:  keepAliveInterval,
	}
}

// NewPoolConfigWithConfig returns a new PoolConfig
func NewPoolConfigWithConfig(config *Config, maxConnections, initConnections, maxIdleConnections,
	maxIdleTime, maxWaitTime, maxRetryCount, keepAliveInterval int) *PoolConfig {
	return &PoolConfig{
		Config:             config,
		MaxConnections:     maxConnections,
		InitConnections:    initConnections,
		MaxIdleConnections: maxIdleConnections,
		MaxIdleTime:        maxIdleTime,
		MaxWaitTime:        maxWaitTime,
		MaxRetryCount:      maxRetryCount,
		KeepAliveInterval:  keepAliveInterval,
	}
}

// Validate validates pool config
func (pc *PoolConfig) Validate() error {
	// validate MaxConnections
	if pc.MaxConnections <= constant.ZeroInt {
		return errors.New("maximum connection argument should larger than 0")
	}
	// validate InitConnections
	if pc.InitConnections < constant.ZeroInt {
		return errors.New("init connection argument should not be smaller than 0")
	}
	if pc.InitConnections > pc.MaxConnections {
		return errors.Errorf("init connections should be less or equal than maximum connections. init_connections: %d, max_connections: %d",
			pc.InitConnections, pc.MaxConnections)
	}
	// validate MaxIdleConnections
	if pc.MaxIdleConnections < constant.ZeroInt {
		return errors.New("maximum idle connection argument should not be smaller than 0")
	}
	if pc.MaxIdleConnections > pc.MaxConnections {
		return errors.New("maximum idle connection argument should not be larger than maximum connection argument")
	}
	// validate MaxIdleTime
	if pc.MaxIdleTime <= constant.ZeroInt {
		return errors.New("maximum idle time argument should be larger than 0")
	}
	// validate MaxWaitTime
	if pc.MaxWaitTime < DefaultUnlimitedWaitTime {
		return errors.New("maximum wait time argument should not be smaller than -1")
	}
	// validate MaxRetryCount
	if pc.MaxRetryCount < DefaultUnlimitedRetryCount {
		return errors.New("maximum retry count argument should not be smaller than -1")
	}
	// validate KeepAliveInterval
	if pc.KeepAliveInterval <= constant.ZeroInt {
		return errors.New("keep alive interval argument should be larger than 0")
	}

	return nil
}

// Clone returns a new PoolConfig with same values
func (pc *PoolConfig) Clone() *PoolConfig {
	return &PoolConfig{
		Config:             pc.Config.Clone(),
		MaxConnections:     pc.MaxConnections,
		InitConnections:    pc.InitConnections,
		MaxIdleConnections: pc.MaxIdleConnections,
		MaxIdleTime:        pc.MaxIdleTime,
		MaxWaitTime:        pc.MaxWaitTime,
		MaxRetryCount:      pc.MaxRetryCount,
		KeepAliveInterval:  pc.KeepAliveInterval,
	}
}
