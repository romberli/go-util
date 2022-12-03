package prometheus

import (
	"testing"

	"github.com/romberli/log"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"

	"github.com/romberli/go-util/middleware"
)

func TestPool(t *testing.T) {
	var (
		err  error
		pool *Pool
		conn middleware.PoolConn
	)

	asst := assert.New(t)

	log.SetLevel(zapcore.DebugLevel)

	// create pool
	config := NewConfigWithBasicAuth(defaultAddr, defaultUser, defaultPass)
	pool, err = NewPoolWithConfig(config, DefaultMaxConnections, DefaultInitConnections, DefaultMaxIdleConnections, DefaultMaxIdleTime, DefaultMaxWaitTime, DefaultMaxRetryCount, DefaultKeepAliveInterval)
	asst.Nil(err, "create pool failed. addr: %s, user: %s, pass: %s", defaultAddr, defaultAddr, defaultAddr)

	// get connection from the pool
	conn, err = pool.Get()
	asst.Nil(err, "get connection from pool failed")

	// test connection
	ok := conn.(*PoolConn).CheckInstanceStatus()
	asst.True(ok, "check instance failed")

	err = conn.Close()
	asst.Nil(err, "close connection failed")

	err = pool.Close()
	asst.Nil(err, "close pool failed")
}
