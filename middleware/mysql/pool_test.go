package mysql

import (
	"testing"
	"time"

	"github.com/romberli/log"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"

	"github.com/romberli/go-util/middleware"
)

func TestPool(t *testing.T) {
	var (
		err       error
		pool      *Pool
		conn      middleware.PoolConn
		repRole   string
		slaveList []string
		result    *Result
	)

	asst := assert.New(t)

	log.SetLevel(zapcore.DebugLevel)

	addr := "192.168.137.11:3306"
	dbName := "testStruct"
	dbUser := "root"
	dbPass := "root"

	// create pool
	pool, err = NewMySQLPoolWithDefault(addr, dbName, dbUser, dbPass)
	asst.Nil(err, "create pool failed. addr: %s, dbName: %s, dbUser: %s, dbPass: %s", addr, dbName, dbUser, dbPass)

	// get connection from the pool
	conn, err = pool.Get()
	asst.Nil(err, "get connection from pool failed")

	// testStruct connection
	slaveList, err = conn.(*PoolConn).GetReplicationSlaveList()
	asst.Nil(err, "get replication slave list failed")
	t.Logf("replication slave list: %v", slaveList)

	err = conn.Close()
	asst.Nil(err, "close connection failed")
	conn, err = pool.Get()
	asst.Nil(err, "get connection from pool failed")

	result, err = conn.(*PoolConn).GetReplicationSlavesStatus()
	asst.Nil(err, "get replication slave status failed")
	if result.RowNumber() > 0 {
		t.Logf("show slave status: %v", result.Values)
	} else {
		t.Logf("this is not a slave node")
	}
	repRole, err = conn.(*PoolConn).GetReplicationRole()
	asst.Nil(err, "get replication role failed")
	t.Logf("replication role: %s", repRole)

	// sleep to testStruct maintain mechanism
	time.Sleep(20 * time.Second)

	err = pool.Close()
	asst.Nil(err, "close pool failed")
}
