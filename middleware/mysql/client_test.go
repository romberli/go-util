package mysql

import (
	"github.com/siddontang/go-mysql/mysql"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMySQLConnection(t *testing.T) {
	var (
		err       error
		conn      *Conn
		repRole   string
		slaveList []string
		result    *mysql.Result
	)

	assert := assert.New(t)

	addr := "192.168.137.11:3306"
	dbName := "test"
	dbUser := "root"
	dbPass := "xxx"

	conn, err = NewMySQLConn(addr, dbName, dbUser, dbPass)
	assert.Nil(err, "connect to mysql failed. addr: %s, dbName: %s, dbUser: %s, dbPass: %s",
		addr, dbName, dbUser, dbPass)
	defer func() {
		err = conn.Close()
		assert.Nil(err, "close connection failed.")
	}()

	slaveList, err = conn.GetReplicationSlaveList()
	assert.Nil(err, "get replication slave list failed.")
	t.Logf("replication slave list: %v", slaveList)

	result, err = conn.GetReplicationSlavesStatus()
	assert.Nil(err, "get replication slave status failed.")
	if result.RowNumber() > 0 {
		t.Logf("show slave status: %v", result.Values)
	} else {
		t.Logf("this is not a slave node.")
	}
	repRole, err = conn.GetReplicationRole()
	assert.Nil(err, "get replication role failed.")
	t.Logf("replication role: %s", repRole)
}
