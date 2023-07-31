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
		repRole   ReplicationRole
		slaveList []string
		result    *Result
	)

	asst := assert.New(t)

	log.SetLevel(zapcore.DebugLevel)
	log.SetDisableEscape(true)

	addr := "192.168.137.11:3306"
	dbName := "test"
	dbUser := "root"
	dbPass := "root"

	// create pool
	pool, err = NewPoolWithDefault(addr, dbName, dbUser, dbPass)
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

func TestPool_Transaction(t *testing.T) {
	var (
		err  error
		pool *Pool
	)

	asst := assert.New(t)

	log.SetLevel(zapcore.DebugLevel)

	addr := "192.168.137.11:3306"
	dbName := "test"
	dbUser := "root"
	dbPass := "root"

	// create pool
	pool, err = NewPool(addr, dbName, dbUser, dbPass, 1, 1, 1, 10000, 1, -1, 1000)
	asst.Nil(err, "create pool failed. addr: %s, dbName: %s, dbUser: %s, dbPass: %s", addr, dbName, dbUser, dbPass)

	// get transaction from the pool
	trx, err := pool.Transaction()
	asst.Nil(err, "get transaction from the pool failed")

	defer func() {
		err = trx.Close()
		asst.Nil(err, "Close transaction failed")
	}()

	// create table
	err = createTable()
	asst.Nil(err, " execute create table sql failed")

	err = trx.Begin()
	asst.Nil(err, "Begin transaction failed")

	sql := "INSERT INTO t05 (id, col1) VALUES (?, ?)"
	_, err = trx.Execute(sql, 1, 1)
	asst.Nil(err, "execute sql failed. sql: %s", sql)
	result, err := trx.Execute("SELECT COUNT(*) FROM t01")
	asst.Nil(err, "execute sql failed. sql: %s", sql)
	count, err := result.GetInt(0, 0)
	asst.Nil(err, "get count failed")
	t.Logf("count: %d", count)

	// err = trx.Close()
	// asst.Nil(err, "Close transaction failed")

	// trx, err = pool.Transaction()
	// asst.Nil(err, "get transaction from the pool failed")
	// err = trx.Begin()
	// asst.Nil(err, "Begin transaction failed")
	_, err = trx.Execute(sql, 2, 2)
	asst.Nil(err, "execute sql failed. sql: %s", sql)
	result, err = trx.Execute("SELECT COUNT(*) FROM t01")
	asst.Nil(err, "execute sql failed. sql: %s", sql)
	count, err = result.GetInt(0, 0)
	asst.Nil(err, "get count failed")
	t.Logf("count: %d", count)

	err = trx.Commit()
	asst.Nil(err, "Begin transaction failed")

	// drop table
	err = dropTable()
	asst.Nil(err, " execute drop table sql failed")
}
