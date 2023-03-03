package clickhouse

import (
	"context"
	"errors"

	"github.com/romberli/go-util/middleware"
)

var _globalPool *Pool

// InitGlobalPool returns a new *Pool and replaces it as global pool
func InitGlobalPool(addr, dbName, dbUser, dbPass string, debug bool,
	maxConnections, initConnections, maxIdleConnections, maxIdleTime, maxWaitTime, maxRetryCount, keepAliveInterval int, altHosts ...string) error {
	cfg := NewPoolConfig(addr, dbName, dbUser, dbPass, debug, maxConnections, initConnections,
		maxIdleConnections, maxIdleTime, maxWaitTime, maxRetryCount, keepAliveInterval, altHosts...)

	return InitGlobalPoolWithPoolConfig(cfg)
}

// InitGlobalPoolWithDefault returns a new *Pool with default configuration and replaces it as global pool
func InitGlobalPoolWithDefault(addr, dbName, dbUser, dbPass string) error {
	return InitGlobalPool(addr, dbName, dbUser, dbPass, false,
		DefaultMaxConnections, DefaultInitConnections, DefaultMaxIdleConnections, DefaultMaxIdleTime,
		DefaultMaxWaitTime, DefaultMaxRetryCount, DefaultKeepAliveInterval)
}

// InitGlobalPoolWithConfig returns a new *Pool with a Config object and replaces it as global pool
func InitGlobalPoolWithConfig(config Config, maxConnections, initConnections, maxIdleConnections, maxIdleTime,
	maxWaitTime, maxRetryCount, keepAliveInterval int) error {
	cfg := NewPoolConfigWithConfig(config, maxConnections, initConnections, maxIdleConnections, maxIdleTime,
		maxWaitTime, maxRetryCount, keepAliveInterval)

	return InitGlobalPoolWithPoolConfig(cfg)
}

// InitGlobalPoolWithPoolConfig returns a new *Pool with a PoolConfig object and replaces it as global pool
func InitGlobalPoolWithPoolConfig(config PoolConfig) error {
	pool, err := NewPoolWithPoolConfig(config)
	if err != nil {
		return err
	}

	return ReplaceGlobalPool(pool)
}

// ReplaceGlobalPool replaces given pool as global pool
func ReplaceGlobalPool(pool *Pool) error {
	if _globalPool != nil {
		err := _globalPool.Close()
		if err != nil {
			return err
		}
	}

	_globalPool = pool
	return nil
}

// IsClosed returns if global pool had been closed
func IsClosed() bool {
	return _globalPool.IsClosed()
}

// Supply creates given number of connections and add them to free connection channel of global pool
func Supply(num int) error {
	return _globalPool.Supply(num)
}

// Close closes global pool, it sets global pool to nil pointer
func Close() error {
	if _globalPool == nil {
		return nil
	}

	err := _globalPool.Close()
	_globalPool = nil
	return err
}

// Get gets a connection from pool and validate it,
// if there is no valid connection in the pool, it will create a new connection
func Get() (*PoolConn, error) {
	conn, err := _globalPool.Get()
	if err != nil {
		return nil, err
	}

	return conn.(*PoolConn), nil
}

// Release releases given number of connections of global pool, each connection will disconnect with database
func Release(num int) error {
	return _globalPool.Release(num)
}

// Execute executes given command
func Execute(command string, args ...interface{}) (middleware.Result, error) {
	return executeContext(context.Background(), command, args...)
}

// ExecuteContext executes given command with context
func ExecuteContext(ctx context.Context, command string, args ...interface{}) (middleware.Result, error) {
	return executeContext(ctx, command, args...)
}

// executeContext executes given command with context
func executeContext(ctx context.Context, command string, args ...interface{}) (middleware.Result, error) {
	if _globalPool == nil {
		return nil, errors.New("global pool is nil, please initiate it first")
	}

	pc, err := _globalPool.Get()
	if err != nil {
		return nil, err
	}
	defer func() { _ = pc.Close() }()

	return pc.ExecuteContext(ctx, command, args...)
}
