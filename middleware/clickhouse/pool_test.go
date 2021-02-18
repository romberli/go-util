package clickhouse

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
	pool, err = NewClickhousePoolWithDefault(addr, dbName, dbUser, dbPass)
	asst.Nil(err, "create pool failed. addr: %s, dbName: %s, dbUser: %s, dbPass: %s", addr, dbName, dbUser, dbPass)

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
