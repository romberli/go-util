package mysql

import (
	"fmt"
	"testing"
	"time"

	"github.com/romberli/go-util/common"
	"github.com/romberli/go-util/constant"
	"github.com/romberli/go-util/middleware"
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
		create table if not exists t05(
			id int(11) auto_increment primary key,
			name varchar(100),
			col1 int(11),
			col2 decimal(16, 4),
			last_update_time datetime(6) not null default current_timestamp(6) on update current_timestamp(6)
		) engine=innodb character set utf8mb4;
	`
	_, err := conn.Execute(sql)
	return err
}

func dropTable() error {
	sql := `drop table if exists t05;`
	_, err := conn.Execute(sql)
	return err
}

func TestMySQLConnection(t *testing.T) {
	var (
		err       error
		repRole   ReplicationRole
		slaveList []string
		result    *Result
		inClause  string
	)

	asst := assert.New(t)

	defer func() {
		err = conn.Close()
		asst.Nil(err, "close connection failed")
	}()

	// drop table
	err = dropTable()
	asst.Nil(err, "execute drop table sql failed")
	// create table
	err = createTable()
	asst.Nil(err, "execute create table sql failed")
	// insert data
	ts := newTestStruct("aa", 1, 3.14)
	tsEmpty := newTestStructWithDefault()
	sql := `insert into t05(name, col1, col2) values(?, ?, ?), (?, ?, ?);`
	result, err = conn.Execute(sql, ts.Name, ts.Col1, ts.Col2, tsEmpty.Name, tsEmpty.Col1, tsEmpty.Col2)
	asst.Nil(err, "execute insert sql failed")

	// select data
	interfaces, err := common.ConvertInterfaceToSliceInterface([]string{ts.Name, "bb"})
	asst.Nil(err, "execute select sql failed")
	inClause, err = middleware.ConvertSliceToString(interfaces...)
	timeStr := time.Now().Add(-time.Hour).Format(constant.DefaultTimeLayout)
	sql = `select id, name, col1, col2, last_update_time from t05 where name in (%s) and last_update_time >= ?`
	sql = fmt.Sprintf(sql, inClause)
	result, err = conn.Execute(sql, timeStr)
	asst.Nil(err, "execute select sql failed")
	asst.Equal(1, result.RowNumber(), "execute select sql failed")

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
