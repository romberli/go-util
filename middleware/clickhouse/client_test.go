package clickhouse

import (
	"fmt"
	"testing"
	"time"

	"github.com/ClickHouse/clickhouse-go"
	"github.com/romberli/log"
	"github.com/stretchr/testify/assert"
)

const (
	addr   = "192.168.137.11:9000"
	dbName = "pmm"
	dbUser = ""
	dbPass = ""
)

var conn = initConn()

type testRow struct {
	Tables []string `middleware:"tables"`
}

func initConn() *Conn {
	config := NewClickhouseConfigWithDefault(addr, dbName, dbUser, dbPass)
	conn, err := NewClickhouseConnWithConfig(config)
	if err != nil {
		log.Error(fmt.Sprintf("init connection failed.\n%s", err.Error()))
		return nil
	}

	return conn
}

func TestAll(t *testing.T) {
	TestConn_Execute(t)
}

func TestConn_Execute(t *testing.T) {
	asst := assert.New(t)

	// create table
	sql := `
		create table t01
		(
			id               Int64,
			name             LowCardinality(String),
			group            Array(LowCardinality(String)),
			type             Enum8('a'=0, 'b'=1, 'c'=2),
			del_flag         Int8,
			create_time      Datetime,
			last_update_time Datetime
		)
			engine = MergeTree PARTITION BY toYYYYMMDD(last_update_time)
			ORDER BY (id, last_update_time) SETTINGS index_granularity = 8192;
	`
	_, err := conn.Execute(sql)
	asst.Nil(err, "test Execute() failed")
	// insert data
	err = conn.Begin()
	asst.Nil(err, "test Execute() failed")
	sql = `insert into t01(id, name, group, type, del_flag, create_time, last_update_time) values(?, ?, ?, ?, ?, ?, ?)`
	_, err = conn.Execute(sql, 1, "aaa", clickhouse.Array([]string{"group1", "group2", "group3"}), "a", 0, time.Now(), time.Now())
	asst.Nil(err, "test Execute() failed")
	err = conn.Commit()
	asst.Nil(err, "test Execute() failed")
	// select data
	sql = `select id, name, group, type, del_flag, create_time, last_update_time from t01`
	result, err := conn.Execute(sql)
	asst.Nil(err, "test Execute() failed")
	asst.Equal(1, result.RowNumber(), "test execute failed")
	// map to struct
	r := &testRow{}
	err = result.MapToStructByRowIndex(r, 0, "middleware")
	// drop table
	sql = `drop table if exists t01;`
	_, err = conn.Execute(sql)
	asst.Nil(err, "test Execute() failed")
}
