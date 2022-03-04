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
	addr   = "192.168.137.11:9000"
	dbName = "pmm"
	dbUser = ""
	dbPass = ""
)

var testConn = initConn()

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
	_, err := testConn.Execute(sql)

	return err
}

func dropTable() error {
	sql := `drop table if exists t01;`
	_, err := testConn.Execute(sql)

	return err
}

func TestConnAll(t *testing.T) {
	TestConn_Execute(t)
	TestConn_GetTimeZone(t)
}

func TestConn_Execute(t *testing.T) {
	asst := assert.New(t)

	// create table
	err := createTable()
	asst.Nil(err, "test Execute() failed")
	// insert data
	err = testConn.Begin()
	asst.Nil(err, "test Execute() failed")
	sql := `insert into t01(id, name, group, type, del_flag, create_time, last_update_time) values(?, ?, ?, ?, ?, ?, ?)`
	_, err = testConn.Execute(sql, 1, constant.DefaultRandomString, clickhouse.Array([]string{"group1", "group2", "group3"}), "a", 0, constant.DefaultRandomTime, time.Now())
	asst.Nil(err, "test Execute() failed")
	err = testConn.Commit()
	asst.Nil(err, "test Execute() failed")
	err = testConn.Begin()
	asst.Nil(err, "test Execute() failed")
	sql = `insert into t01(id, name, group, type, del_flag, create_time, last_update_time) values(?, ?, ?, ?, ?, ?, ?)`
	_, err = testConn.Execute(sql, 2, constant.DefaultRandomString, clickhouse.Array([]string{}), "a", 0, constant.DefaultRandomTime, time.Now())
	asst.Nil(err, "test Execute() failed")
	err = testConn.Commit()
	asst.Nil(err, "test Execute() failed")
	// select data
	sql = `select id, name, group, type, del_flag, create_time, last_update_time from t01 where last_update_time > ? order by id asc limit ?, ?`
	result, err := testConn.Execute(sql, time.Now().Add(-time.Hour), 1, 1)
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

func TestConn_GetTimeZone(t *testing.T) {
	asst := assert.New(t)

	var (
		err    error
		sql    string
		result *Result
	)

	// sql = `select toUnixTimestamp(now());`
	// result, err = testConn.Execute(sql)
	// asst.Nil(err, "test GetTimeZone() failed")
	// now, err := result.GetInt(constant.ZeroInt, constant.ZeroInt)
	// asst.Nil(err, "test GetTimeZone() failed")
	// asst.Equal(int64(now), time.Now().Unix(), "test GetTimeZone() failed")
	// fmt.Println(int64(now), time.Now().Unix())

	sql = `
        select sm.sql_id,
               m.fingerprint,
               m.example,
               m.db_name,
               sm.exec_count,
               sm.total_exec_time,
               sm.avg_exec_time,
               sm.rows_examined_max
        from (
                 select queryid                                               as sql_id,
                        sum(num_queries)                                      as exec_count,
                        truncate(sum(m_query_time_sum), 2)                    as total_exec_time,
                        truncate(sum(m_query_time_sum) / sum(num_queries), 2) as avg_exec_time,
                        max(m_rows_examined_max)                              as rows_examined_max
                 from metrics
                 where service_type = 'mysql'
                   and service_name in ('192-168-137-11-mysql')
                   and period_start >= ?
                   and period_start < ?
                   and m_rows_examined_max >= ?
                 group by queryid
                 order by rows_examined_max desc
                 limit ? offset ? ) sm
                 left join (select queryid          as sql_id,
                                   max(fingerprint) as fingerprint,
                                   max(example)     as example,
                                   max(database)    as db_name
                            from metrics
                            where service_type = 'mysql'
                              and service_name in ('192-168-137-11-mysql')
                              and period_start >= ?
                              and period_start < ?
                              and m_rows_examined_max >= ?
                            group by queryid) m
                           on sm.sql_id = m.sql_id;
	`
	tz, err := testConn.GetTimeZone()
	asst.Nil(err, "test GetTimeZone() failed")
	startTime, err := time.ParseInLocation(constant.TimeLayoutSecond, "2022-03-03 10:00:00", time.Local)
	asst.Nil(err, "test GetTimeZone() failed")
	startTime = startTime.In(tz)
	endTime, err := time.ParseInLocation(constant.TimeLayoutSecond, "2022-03-04 18:00:00", time.Local)
	asst.Nil(err, "test GetTimeZone() failed")
	endTime = endTime.In(tz)

	// endTime := "2022-03-04 10:00:00"
	minRowsExamined := 0
	limit := 5
	offset := 0
	result, err = testConn.Execute(sql, startTime, endTime, minRowsExamined, limit, offset, startTime, endTime, minRowsExamined)
	// result, err := testConn.Execute(sql)
	asst.Nil(err, "test GetTimeZone() failed")

	fmt.Println(result)
}
