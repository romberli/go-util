package mysql

import (
	"errors"
)

var _globalPool *Pool

// InitMySQLGlobalPool returns a new *Pool and replaces it as global pool
func InitMySQLGlobalPool(addr, dbName, dbUser, dbPass string,
	maxConnections, initConnections, maxIdleConnections, maxIdleTime, keepAliveInterval int) error {
	cfg := NewPoolConfig(addr, dbName, dbUser, dbPass, maxConnections, initConnections, maxIdleConnections, maxIdleTime, keepAliveInterval)

	return InitMySQLGlobalPoolWithPoolConfig(cfg)
}

// InitMySQLGlobalPoolWithDefault returns a new *Pool with default configuration and replaces it as global pool
func InitMySQLGlobalPoolWithDefault(addr, dbName, dbUser, dbPass string) error {
	return InitMySQLGlobalPool(addr, dbName, dbUser, dbPass,
		DefaultMaxConnections, DefaultInitConnections, DefaultMaxIdleConnections, DefaultMaxIdleTime, DefaultKeepAliveInterval)
}

// InitMySQLGlobalPoolWithConfig returns a new *Pool with a Config object and replaces it as global pool
func InitMySQLGlobalPoolWithConfig(config Config, maxConnections, initConnections, maxIdleConnections, maxIdleTime, keepAliveInterval int) error {
	cfg := NewPoolConfigWithConfig(config, maxConnections, initConnections, maxIdleConnections, maxIdleTime, keepAliveInterval)

	return InitMySQLGlobalPoolWithPoolConfig(cfg)
}

// InitMySQLGlobalPoolWithPoolConfig returns a new *Pool with a PoolConfig object and replaces it as global pool
func InitMySQLGlobalPoolWithPoolConfig(config PoolConfig) error {
	pool, err := NewMySQLPoolWithPoolConfig(config)
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

// Get get gets a connection from pool and validate it,
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

// Execute execute given sql statement
func Execute(sql string, args ...interface{}) (interface{}, error) {
	if _globalPool == nil {
		return nil, errors.New("global pool is nil, please initiate it first")
	}

	pc, err := _globalPool.Get()
	if err != nil {
		return nil, err
	}
	defer func() { _ = pc.Close() }()

	return pc.Execute(sql, args...)
}
