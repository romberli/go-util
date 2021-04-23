package clickhouse

import (
	"context"
	"testing"
	"time"

	"github.com/ClickHouse/clickhouse-go"
	"github.com/romberli/go-util/constant"
	"github.com/stretchr/testify/assert"
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
	err = conn.Begin()
	asst.Nil(err, "test ExecuteContext() failed")
	sql := `insert into t01(id, name, group, type, del_flag, create_time, last_update_time) values(?, ?, ?, ?, ?, ?, ?);`
	stmt, err := conn.Prepare(sql)
	asst.Nil(err, "test Execute() failed")
	_, err = stmt.Execute(1, constant.DefaultRandomString, clickhouse.Array([]string{"group1", "group2", "group3"}), "a", 0, constant.DefaultRandomTime, time.Now())
	asst.Nil(err, "test Execute() failed")
	err = conn.Commit()
	asst.Nil(err, "test Execute() failed")
	// select data
	sql = `select id, name, group, type, del_flag, create_time, last_update_time from t01 where id = ?;`
	stmt, err = conn.Prepare(sql)
	asst.Nil(err, "test Execute() failed")
	result, err := stmt.Execute(1)
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
	err = conn.Begin()
	asst.Nil(err, "test ExecuteContext() failed")
	sql := `insert into t01(id, name, group, type, del_flag, create_time, last_update_time) values(?, ?, ?, ?, ?, ?, ?);`
	stmt, err := conn.PrepareContext(ctx, sql)
	asst.Nil(err, "test ExecuteContext() failed")
	_, err = stmt.ExecuteContext(ctx, 1, constant.DefaultRandomString, clickhouse.Array([]string{"group1", "group2", "group3"}), "a", 0, constant.DefaultRandomTime, time.Now())
	asst.Nil(err, "test ExecuteContext() failed")
	err = conn.Commit()
	asst.Nil(err, "test ExecuteContext() failed")
	// select data
	sql = `select id, name, group, type, del_flag, create_time, last_update_time from t01 where id = ?;`
	stmt, err = conn.PrepareContext(ctx, sql)
	asst.Nil(err, "test ExecuteContext() failed")
	result, err := stmt.ExecuteContext(ctx, 1)
	asst.Nil(err, "test ExecuteContext() failed")
	asst.Equal(1, result.RowNumber(), "test ExecuteContext() failed")
	// drop table
	err = dropTable()
	asst.Nil(err, "test ExecuteContext() failed")
}
