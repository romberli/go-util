package mysql

import (
	"database/sql/driver"
	"strings"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/go-mysql-org/go-mysql/mysql"

	"github.com/romberli/go-util/constant"
	"github.com/romberli/go-util/middleware/result"
)

func MockClientNewConn() *gomonkey.Patches {
	patches := gomonkey.ApplyFunc(NewConn, func(addr, dbName, dbUser, dbPass string) (*Conn, error) {
		config := NewConfig(addr, dbName, dbUser, dbPass)
		conn := &Conn{
			Config: config,
			Conn:   nil,
		}

		return conn, nil
	})

	return patches
}

func MockClientExecute(conn *Conn, sql string, args ...interface{}) *gomonkey.Patches {
	if strings.Trim(strings.ToLower(sql), constant.SemicolonString) == "select 1" {
		patches := gomonkey.ApplyMethod(conn, "Execute", func(_ *Conn, sql string, args ...interface{}) (*Result, error) {
			return &Result{
				Raw: &mysql.Result{
					Resultset: &mysql.Resultset{},
				},
				Rows: result.NewRows([]string{"1"}, map[string]int{"1": 1}, [][]driver.Value{{1}}),
				Map:  nil,
			}, nil
		})

		return patches
	}

	patches := gomonkey.ApplyMethod(conn, "Execute", func(_ *Conn, sql string, args ...interface{}) (*Result, error) {
		return &Result{
			Raw: &mysql.Result{
				Resultset: &mysql.Resultset{},
			},
			Rows: result.NewRows([]string{}, map[string]int{}, [][]driver.Value{}),
			Map:  nil,
		}, nil
	})

	return patches
}
