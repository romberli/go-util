package mysql

import (
	"fmt"
	"testing"
	"time"

	"github.com/romberli/log"
	"github.com/stretchr/testify/assert"

	"github.com/romberli/go-util/common"
	"github.com/romberli/go-util/constant"
	"github.com/romberli/go-util/middleware"
)

const (
	testGlobalVariableName  = "read_only"
	testGlobalVariableValue = "OFF"
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

var conn *Conn

func init() {
	conn = initConn()
}

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
		CREATE TABLE IF NOT EXISTS t05(
			id int(11) AUTO_INCREMENT PRIMARY KEY,
			name varchar(100),
			col1 int(11),
			col2 decimal(16, 4),
			last_update_time datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6)
		) ENGINE=innodb CHARACTER SET utf8mb4;
	`
	_, err := conn.Execute(sql)
	return err
}

func dropTable() error {
	sql := `DROP TABLE IF EXISTS t05;`
	_, err := conn.Execute(sql)
	return err
}

func TestConn_All(t *testing.T) {
	TestMySQLConnection(t)
	TestConn_Execute(t)
	TestConn_IsReplicationSlave(t)
	TestConn_IsMater(t)
	TestConn_IsMGR(t)
	TestConn_IsReadOnly(t)
	TestConn_IsSuperReadOnly(t)
	TestConn_SetReadOnly(t)
	TestConn_SetSuperReadOnly(t)
	TestConn_ShowGlobalVariable(t)
	TestConn_SetGlobalVariable(t)
	TestConn_SetGlobalVariables(t)
}

func TestMySQLConnection(t *testing.T) {
	var (
		err       error
		sql       string
		repRole   ReplicationRole
		slaveList []string
		result    *Result
		inClause  string
	)

	asst := assert.New(t)

	// defer func() {
	//	err = conn.Close()
	//	asst.Nil(err, "close connection failed")
	// }()

	// drop table
	err = dropTable()
	asst.Nil(err, "execute drop table sql failed")
	// create table
	err = createTable()
	asst.Nil(err, "execute create table sql failed")
	// insert data
	ts := newTestStruct("aa", 1, 3.14)
	tsEmpty := newTestStructWithDefault()
	sql = `INSERT INTO t05(name, col1, col2) VALUES(?, ?, ?), (?, ?, ?);`
	result, err = conn.Execute(sql, ts.Name, ts.Col1, ts.Col2, tsEmpty.Name, tsEmpty.Col1, tsEmpty.Col2)
	asst.Nil(err, "execute insert sql failed")

	// select data
	interfaces, err := common.ConvertInterfaceToSliceInterface([]string{ts.Name, "bb"})
	asst.Nil(err, "execute select sql failed")
	inClause, err = middleware.ConvertSliceToString(interfaces...)
	timeStr := time.Now().Add(-time.Hour).Format(constant.DefaultTimeLayout)
	sql = `SELECT id, name, col1, col2, last_update_time FROM t05 WHERE name IN (%s) AND last_update_time >= ?`
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
	if result.RowNumber() > constant.ZeroInt {
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

func TestConn_Execute(t *testing.T) {
	asst := assert.New(t)
	// create table
	err := createTable()
	asst.Nil(err, "execute create table sql failed")
	// select data
	sql := "SELECT col2 FROM test.t05 WHERE id = 1;"
	result, err := conn.Execute(sql)
	asst.Nil(err, "execute sql failed")
	col2, err := result.GetFloat(0, 0)
	asst.Equal(0.0, col2, "execute sql failed")
	// drop table
	err = dropTable()
	asst.Nil(err, "execute drop table sql failed")
}

func TestConn_IsReplicationSlave(t *testing.T) {
	asst := assert.New(t)

	isSlave, err := conn.IsReplicationSlave()
	asst.Nil(err, "test IsReplicationSlave() failed")
	asst.False(isSlave, "test IsReplicationSlave() failed")
}

func TestConn_IsMater(t *testing.T) {
	asst := assert.New(t)

	isMaster, err := conn.IsMaster()
	asst.Nil(err, "test IsMaster() failed")
	asst.True(isMaster, "test IsMaster() failed")
}

func TestConn_IsMGR(t *testing.T) {
	asst := assert.New(t)

	isMGR, err := conn.IsMGR()
	asst.Nil(err, "test IsMGR() failed")
	asst.False(isMGR, "test IsMGR() failed")
}

func TestConn_IsReadOnly(t *testing.T) {
	asst := assert.New(t)

	status, err := conn.IsReadOnly()
	asst.Nil(err, "test IsReadOnly() failed")
	asst.False(status, "test IsReadOnly() failed")
}

func TestConn_IsSuperReadOnly(t *testing.T) {
	asst := assert.New(t)

	status, err := conn.IsSuperReadOnly()
	asst.Nil(err, "test IsSuperReadOnly() failed")
	asst.False(status, "test IsSuperReadOnly() failed")
}

func TestConn_SetReadOnly(t *testing.T) {
	asst := assert.New(t)

	err := conn.SetReadOnly(true)
	asst.Nil(err, "test SetReadOnly() failed")
	status, err := conn.IsReadOnly()
	asst.Nil(err, "test SetReadOnly() failed")
	asst.True(status, "test SetReadOnly() failed")
	err = conn.SetReadOnly(false)
	asst.Nil(err, "test SetReadOnly() failed")
	status, err = conn.IsReadOnly()
	asst.Nil(err, "test SetReadOnly() failed")
	asst.False(status, "test SetReadOnly() failed")
}

func TestConn_SetSuperReadOnly(t *testing.T) {
	asst := assert.New(t)

	err := conn.SetSuperReadOnly(true)
	asst.Nil(err, "test SetSuperReadOnly() failed")
	status, err := conn.IsSuperReadOnly()
	asst.Nil(err, "test SetSuperReadOnly() failed")
	asst.True(status, "test SetSuperReadOnly() failed")
	err = conn.SetSuperReadOnly(false)
	asst.Nil(err, "test SetSuperReadOnly() failed")
	status, err = conn.IsSuperReadOnly()
	asst.Nil(err, "test SetSuperReadOnly() failed")
	asst.False(status, "test SetSuperReadOnly() failed")
	err = conn.SetReadOnly(false)
	asst.Nil(err, "test SetSuperReadOnly() failed")
}

func TestConn_ShowGlobalVariable(t *testing.T) {
	asst := assert.New(t)

	value, err := conn.ShowGlobalVariable(testGlobalVariableName)
	asst.Nil(err, "test ShowGlobalVariable() failed")
	asst.Equal(testGlobalVariableValue, value, "test ShowGlobalVariable() failed")
}

func TestConn_SetGlobalVariable(t *testing.T) {
	asst := assert.New(t)

	err := conn.SetGlobalVariable(testGlobalVariableName, testGlobalVariableValue)
	asst.Nil(err, "test SetGlobalVariable() failed")
	value, err := conn.ShowGlobalVariable(testGlobalVariableName)
	asst.Nil(err, "test SetGlobalVariable() failed")
	asst.Equal(testGlobalVariableValue, value, "test SetGlobalVariable() failed")
}

func TestConn_SetGlobalVariables(t *testing.T) {
	asst := assert.New(t)

	variables := map[string]string{
		testGlobalVariableName: testGlobalVariableValue,
	}
	err := conn.SetGlobalVariables(variables)
	asst.Nil(err, "test SetGlobalVariables() failed")
	value, err := conn.ShowGlobalVariable(testGlobalVariableName)
	asst.Nil(err, "test SetGlobalVariables() failed")
	asst.Equal(testGlobalVariableValue, value, "test SetGlobalVariables() failed")
}

func TestTemp(t *testing.T) {
	asst := assert.New(t)

	addr := "192.168.137.12:2883"
	dbName := ""
	dbUser := "root"
	dbPass := "root"

	patches := MockClientNewConn()
	defer patches.Reset()

	conn, err := NewConn(addr, dbName, dbUser, dbPass)
	asst.Nil(err, "new conn failed")

	sql := `CREATE TENANT %s RESOURCE_POOL_LIST = ('%s'), CHARSET = '%s', PRIMARY_ZONE = '%s' set OB_TCP_INVITED_NODES = '%%' ;`
	// sql := `CREATE TENANT ? RESOURCE_POOL_LIST = ?, CHARSET = ?, PRIMARY_ZONE = ? set OB_TCP_INVITED_NODES = '%%' ;`
	zones := "zone1,zone2,zone3"
	sql = fmt.Sprintf(sql, "test", "test", "utf8mb4", zones)

	sql = `drop tenant test;`
	result, err := conn.Execute(sql)
	asst.Nil(err, "execute sql failed")
	t.Logf("result: %v", result)
}

func TestMock(t *testing.T) {
	asst := assert.New(t)

	addr := "192.168.137.12:2883"
	dbName := ""
	dbUser := "root"
	dbPass := "root"

	patches := MockClientNewConn()
	defer patches.Reset()

	conn, err := NewConn(addr, dbName, dbUser, dbPass)
	asst.Nil(err, "new conn failed")

	patches.Reset()
	sql := `select 1;`
	patches = MockClientExecute(conn, sql)
	defer patches.Reset()

	result, err := conn.Execute(sql)
	asst.Nil(err, "execute sql failed")
	asst.Equal(1, result.RowNumber(), "execute sql failed")
	r, err := result.GetInt(constant.ZeroInt, constant.ZeroInt)
	asst.Nil(err, "execute sql failed")
	asst.Equal(1, r, "execute sql failed")
}
