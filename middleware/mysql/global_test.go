package mysql

import (
	"testing"
	"time"

	"github.com/romberli/log"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"

	"github.com/romberli/go-util/middleware"
)

func TestGlobalPool(t *testing.T) {
	var (
		err       error
		conn      *PoolConn
		slaveList []string
		result    middleware.Result
	)

	asst := assert.New(t)

	log.SetLevel(zapcore.DebugLevel)

	addr := "192.168.137.11:3306"
	dbName := "test"
	dbUser := "root"
	dbPass := "root"

	// create pool
	err = InitGlobalPoolWithDefault(addr, dbName, dbUser, dbPass)
	asst.Nil(err, "create pool failed. addr: %s, dbName: %s, dbUser: %s, dbPass: %s", addr, dbName, dbUser, dbPass)

	// get connection from the pool
	conn, err = Get()
	asst.Nil(err, "get connection from pool failed")

	// testStruct connection
	slaveList, err = conn.GetReplicationSlaveList()
	asst.Nil(err, "get replication slave list failed")
	t.Logf("replication slave list: %v", slaveList)

	err = conn.Close()
	asst.Nil(err, "close connection failed")

	sql := "select ? as ok;"
	result, err = Execute(sql, 1)
	asst.Nil(err, "execute sql with global pool failed")
	actual, err := result.(*Result).GetIntByName(0, "ok")
	asst.Nil(err, "execute sql with global pool failed")
	asst.Equal(1, actual, "expected and actual values are not equal")

	// sleep to testStruct maintain mechanism
	time.Sleep(10 * time.Second)

	err = Close()
	asst.Nil(err, "close global pool failed")
}
