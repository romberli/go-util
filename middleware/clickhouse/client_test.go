package clickhouse

import (
	"fmt"
	"testing"

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

	sql := `select tables from metrics limit 3`
	result, err := conn.Execute(sql)
	asst.Nil(err, "test Execute() failed")
	asst.Equal(3, result.RowNumber(), "test execute failed")

	r := &testRow{}

	err = result.MapToStructByRowIndex(r, 0, "middleware")

	fmt.Println(fmt.Sprintf("%v", r))

	for _, row := range result.Values {
		t.Log(fmt.Sprintf("%v", row))
	}
}
