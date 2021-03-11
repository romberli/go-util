package mysql

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/romberli/go-util/constant"
)

type testStruct struct {
	Name string
	Col1 int
	Col2 float64
}

func newTestStruct(name string, col1 int, col2 float64) *testStruct {
	return &testStruct{
		Name: name,
		Col1: col1,
		Col2: col2,
	}
}

func newTestStructWithDefault() *testStruct {
	return &testStruct{
		Name: constant.DefaultRandomString,
		Col1: constant.DefaultRandomInt,
		Col2: float64(constant.DefaultRandomInt),
	}
}

func TestMySQLConnection(t *testing.T) {
	var (
		err       error
		conn      *Conn
		repRole   string
		slaveList []string
		result    *Result
	)

	asst := assert.New(t)

	addr := "192.168.137.11:3306"
	dbName := "test"
	dbUser := "root"
	dbPass := "root"

	conn, err = NewMySQLConn(addr, dbName, dbUser, dbPass)
	asst.Nil(err, "connect to mysql failed. addr: %s, dbName: %s, dbUser: %s, dbPass: %s",
		addr, dbName, dbUser, dbPass)
	defer func() {
		err = conn.Close()
		asst.Nil(err, "close connection failed")
	}()

	ts := newTestStructWithDefault()
	sql := `insert into t05(name, col1, col2) values(?, ?, ?);`
	result, err = conn.Execute(sql, ts.Name, ts.Col1, ts.Col2)
	asst.Nil(err, "execute insert sql failed")

	slaveList, err = conn.GetReplicationSlaveList()
	asst.Nil(err, "get replication slave list failed")
	t.Logf("replication slave list: %v", slaveList)

	result, err = conn.GetReplicationSlavesStatus()
	asst.Nil(err, "get replication slave status failed")
	if result.RowNumber() > 0 {
		t.Logf("show slave status: %v", result.Values)
	} else {
		t.Logf("this is not a slave node.")
	}
	repRole, err = conn.GetReplicationRole()
	asst.Nil(err, "get replication role failed")
	t.Logf("replication role: %s", repRole)
}
