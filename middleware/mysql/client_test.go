package mysql

import (
	"fmt"
	"testing"

	"github.com/romberli/go-util/constant"
	"github.com/romberli/log"
	"github.com/stretchr/testify/assert"
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

var conn = initConn()

func initConn() *Conn {
	addr := "192.168.137.11:3306"
	dbName := "test"
	dbUser := "root"
	dbPass := "root"

	c, err := NewConn(addr, dbName, dbUser, dbPass)
	if err != nil {
		log.Error(fmt.Sprintf("init connection failed.\n%s", err.Error()))
		return nil
	}

	return c
}

func createTable() error {
	sql := `
		create table if not exists t10(
			id int(11) auto_increment primary key,
			name varchar(100),
			col1 int(11),
			col2 decimal(16, 4)
		) engine=innodb character set utf8mb4;
	`
	_, err := conn.Execute(sql)
	return err
}

func dropTable() error {
	sql := `drop table t10;`
	_, err := conn.Execute(sql)
	return err
}

func TestMySQLConnection(t *testing.T) {
	var (
		err       error
		repRole   ReplicationRole
		slaveList []string
		result    *Result
	)

	asst := assert.New(t)

	// defer func() {
	// 	err = conn.Close()
	// 	asst.Nil(err, "close connection failed")
	// }()
	// create table
	err = createTable()
	asst.Nil(err, "execute create sql failed")
	// insert data
	ts := newTestStructWithDefault()
	sql := `insert into t05(name, col1, col2) values(?, ?, ?);`
	result, err = conn.Execute(sql, ts.Name, ts.Col1, ts.Col2)
	asst.Nil(err, "execute insert sql failed")
	// check replication
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
	// drop table
	err = dropTable()
	asst.Nil(err, "execute drop sql failed")
}
