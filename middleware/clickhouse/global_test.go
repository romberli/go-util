package clickhouse

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/romberli/go-util/middleware"
)

func TestGlobalPool(t *testing.T) {
	var (
		err    error
		conn   *PoolConn
		result middleware.Result
	)

	asst := assert.New(t)

	// create pool
	err = InitGlobalPoolWithDefault(addr, dbName, dbUser, dbPass)
	asst.Nil(err, "create pool failed. addr: %s, dbName: %s, dbUser: %s, dbPass: %s", addr, dbName, dbUser, dbPass)

	// get connection from the pool
	conn, err = Get()
	asst.Nil(err, "get connection from pool failed")

	// test connection
	ok := conn.CheckInstanceStatus()
	asst.True(ok, "check instance status failed")

	err = conn.Close()
	asst.Nil(err, "close connection failed")

	sql := "select 1 as ok;"
	result, err = Execute(sql)
	asst.Nil(err, "execute sql with global pool failed")
	actual, err := result.(*Result).GetIntByName(0, "ok")
	asst.Nil(err, "execute sql with global pool failed")
	asst.Equal(1, actual, "expected and actual values are not equal")

	// sleep to test maintain mechanism
	time.Sleep(10 * time.Second)

	err = Close()
	asst.Nil(err, "close global pool failed")
}
