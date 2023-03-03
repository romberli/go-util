package clickhouse

import (
	"context"
	"testing"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/romberli/go-util/constant"
	"github.com/stretchr/testify/assert"
)

const (
	testDateTime = "2105-12-31 00:00:00"
)

func TestStatementAll(t *testing.T) {
	TestStatement_Execute(t)
	TestStatement_ExecuteContext(t)
}

func TestStatement_Execute(t *testing.T) {
	asst := assert.New(t)

	// create table
	err := createTable()
	asst.Nil(err, "test Execute() failed")
	// insert data
	sql := `insert into t01(id, name, group, type, del_flag, create_time, last_update_time) `
	stmt, err := testConn.Prepare(sql)
	asst.Nil(err, "test Execute() failed")
	defer func() { _ = stmt.Close() }()

	_, err = stmt.Execute(int64(1), constant.DefaultRandomString, clickhouse.ArraySet{"group1", "group2", "group3"}, "a", int8(0), testDateTime, time.Now())
	asst.Nil(err, "test Execute() failed")
	_, err = stmt.Execute(int64(2), constant.DefaultRandomString, clickhouse.ArraySet{"group1", "group2", "group3"}, "a", int8(0), testDateTime, time.Now())
	asst.Nil(err, "test Execute() failed")
	err = stmt.Commit()
	asst.Nil(err, "test Execute() failed")
	// select data
	sql = `select id, name, group, type, del_flag, create_time, last_update_time from t01 where id = ?`
	result, err := testConn.Execute(sql, 1)
	asst.Nil(err, "test Execute() failed")
	asst.Equal(1, result.RowNumber(), "test Execute() failed")
	// drop table
	err = dropTable()
	asst.Nil(err, "test Execute() failed")
}

func TestStatement_ExecuteContext(t *testing.T) {
	asst := assert.New(t)
	ctx := context.Background()

	// create table
	err := createTable()
	asst.Nil(err, "test ExecuteContext() failed")
	// insert data

	sql := `insert into t01(id, name, group, type, del_flag, create_time, last_update_time)`
	stmt, err := testConn.PrepareContext(ctx, sql)
	asst.Nil(err, "test ExecuteContext() failed")
	defer func() { _ = stmt.Close() }()

	_, err = stmt.ExecuteContext(ctx, int64(1), constant.DefaultRandomString, clickhouse.ArraySet{"group1", "group2", "group3"}, "a", int8(0), testDateTime, time.Now())
	asst.Nil(err, "test ExecuteContext() failed")
	err = stmt.Commit()
	asst.Nil(err, "test ExecuteContext() failed")

	// select data
	sql = `select id, name, group, type, del_flag, create_time, last_update_time from t01 where id = ?;`
	result, err := testConn.ExecuteContext(ctx, sql, int64(1))
	asst.Nil(err, "test ExecuteContext() failed")
	asst.Equal(1, result.RowNumber(), "test ExecuteContext() failed")
	// drop table
	err = dropTable()
	asst.Nil(err, "test ExecuteContext() failed")
}
