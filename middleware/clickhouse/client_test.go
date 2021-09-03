package clickhouse

import (
	"fmt"
	"testing"
	"time"

	"github.com/ClickHouse/clickhouse-go"
	"github.com/romberli/log"
	"github.com/stretchr/testify/assert"

	"github.com/romberli/go-util/constant"
)

const (
	addr   = "192.168.10.219:9000"
	dbName = "default"
	dbUser = ""
	dbPass = ""
)

var conn = initConn()

type testRow struct {
	Tables []string `middleware:"tables"`
}

func initConn() *Conn {
	config := NewConfigWithDefault(addr, dbName, dbUser, dbPass)
	c, err := NewConnWithConfig(config)
	if err != nil {
		log.Error(fmt.Sprintf("init connection failed.\n%s", err.Error()))
		return nil
	}

	return c
}

func createTable() error {
	sql := `
		create table if not exists t01
		(
			id               Int64,
			name             Nullable(String),
			group            Array(String),
			type             Enum8('a'=0, 'b'=1, 'c'=2),
			del_flag         Int8,
			create_time      Nullable(Datetime),
			last_update_time Datetime
		)
			engine = MergeTree PARTITION BY toYYYYMMDD(last_update_time)
			ORDER BY (id, last_update_time) SETTINGS index_granularity = 8192;
	`
	_, err := conn.Execute(sql)

	return err
}

func dropTable() error {
	sql := `drop table if exists t01;`
	_, err := conn.Execute(sql)

	return err
}

func TestConnAll(t *testing.T) {
	TestConn_Execute(t)
}

func TestConn_Execute(t *testing.T) {
	asst := assert.New(t)

	// create table
	err := createTable()
	asst.Nil(err, "test Execute() failed")
	// insert data
	err = conn.Begin()
	asst.Nil(err, "test Execute() failed")
	sql := `insert into t01(id, name, group, type, del_flag, create_time, last_update_time) values(?, ?, ?, ?, ?, ?, ?)`
	_, err = conn.Execute(sql, 1, constant.DefaultRandomString, clickhouse.Array([]string{"group1", "group2", "group3"}), "a", 0, constant.DefaultRandomTime, time.Now())
	asst.Nil(err, "test Execute() failed")
	err = conn.Commit()
	asst.Nil(err, "test Execute() failed")
	err = conn.Begin()
	asst.Nil(err, "test Execute() failed")
	sql = `insert into t01(id, name, group, type, del_flag, create_time, last_update_time) values(?, ?, ?, ?, ?, ?, ?)`
	_, err = conn.Execute(sql, 2, constant.DefaultRandomString, clickhouse.Array([]string{}), "a", 0, constant.DefaultRandomTime, time.Now())
	asst.Nil(err, "test Execute() failed")
	err = conn.Commit()
	asst.Nil(err, "test Execute() failed")
	// select data
	sql = `select id, name, group, type, del_flag, create_time, last_update_time from t01 where last_update_time > ? order by id asc limit ?, ?`
	result, err := conn.Execute(sql, time.Now().Add(-time.Hour), 1, 1)
	asst.Nil(err, "test Execute() failed")
	id, err := result.GetInt(0, 0)
	asst.Nil(err, "test execute failed")
	asst.Equal(2, id, "test execute failed")
	// map to struct
	r := &testRow{}
	err = result.MapToStructByRowIndex(r, 0, "middleware")
	// drop table
	err = dropTable()
	asst.Nil(err, "test Execute() failed")
}
