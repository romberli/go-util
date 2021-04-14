package clickhouse

import (
	"database/sql/driver"
	"fmt"
	"strings"

	"github.com/ClickHouse/clickhouse-go"

	"github.com/romberli/go-util/common"
	"github.com/romberli/go-util/constant"
	"github.com/romberli/go-util/middleware/sqls"
)

const (
	DefaultDatabase     = "default"
	DefaultReadTimeout  = 10
	DefaultWriteTimeout = 10
)

type Config struct {
	Addr         string
	DBName       string
	DBUser       string
	DBPass       string
	Debug        bool
	ReadTimeout  int
	WriteTimeout int
	AltHosts     []string
}

// NewClickhouseConfig returns a new Config
func NewClickhouseConfig(addr, dbName, dbUser, dbPass string, debug bool, readTimeout, writeTimeout int, altHosts ...string) Config {
	return Config{
		Addr:         addr,
		DBName:       dbName,
		DBUser:       dbUser,
		DBPass:       dbPass,
		Debug:        debug,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		AltHosts:     altHosts,
	}
}

// NewClickhouseConfigWithDefault returns a new Config with default value
func NewClickhouseConfigWithDefault(addr, dbName, dbUser, dbPass string, altHosts ...string) Config {
	return Config{
		Addr:         addr,
		DBName:       dbName,
		DBUser:       dbUser,
		DBPass:       dbPass,
		Debug:        false,
		ReadTimeout:  DefaultReadTimeout,
		WriteTimeout: DefaultWriteTimeout,
		AltHosts:     altHosts,
	}
}

// AltHostsExist checks if alternative hosts is empty
func (c *Config) AltHostsExist() bool {
	if c.AltHosts != nil && len(c.AltHosts) > constant.ZeroInt {
		return true
	}

	return false
}

// AltHostsString converts AltHosts to string
func (c *Config) AltHostsString() string {
	if c.AltHosts == nil {
		return constant.EmptyString
	}

	switch len(c.AltHosts) {
	case 0:
		return constant.EmptyString
	case 1:
		return c.AltHosts[constant.ZeroInt]
	default:
		s := c.AltHosts[constant.ZeroInt]
		for _, host := range c.AltHosts[1:] {
			s = ", " + host
		}

		return s
	}
}

// GetConnectionString generates connection string to clickhouse
func (c *Config) GetConnectionString() string {
	connStr := fmt.Sprintf("tcp://%s?", c.Addr)

	if c.DBName == constant.EmptyString {
		connStr += fmt.Sprintf("database=%s&", DefaultDatabase)
	} else {
		connStr += fmt.Sprintf("database=%s&", c.DBName)
	}
	if c.DBUser != constant.EmptyString {
		connStr += fmt.Sprintf("username=%s&", c.DBUser)
	}
	if c.DBPass != constant.EmptyString {
		connStr += fmt.Sprintf("password=%s&", c.DBPass)
	}

	if c.Debug {
		connStr += "debug=true&"
	}
	if c.ReadTimeout != constant.ZeroInt {
		connStr += fmt.Sprintf("read_timeout=%d&", c.ReadTimeout)
	}
	if c.WriteTimeout != constant.ZeroInt {
		connStr += fmt.Sprintf("write_time=%d&", c.WriteTimeout)
	}
	if c.AltHostsExist() {
		connStr += fmt.Sprintf("alter_hosts=%s&", c.AltHostsString())
	}

	return strings.Trim(connStr, "&")
}

type Conn struct {
	Config
	clickhouse.Clickhouse
}

// NewClickhouseConnWithConfig returns connection to mysql database with given Config
func NewClickhouseConnWithConfig(config Config) (*Conn, error) {
	// connect to Clickhouse
	client, err := clickhouse.OpenDirect(config.GetConnectionString())
	if err != nil {
		return nil, err
	}

	return &Conn{
		config,
		client,
	}, nil
}

// NewClickhouseConn returns connection to Clickhouse database, be aware that addr is host:port style
func NewClickhouseConn(addr, dbName, dbUser, dbPass string, debug bool, readTimeout, writeTimeout int, altHosts ...string) (*Conn, error) {
	config := NewClickhouseConfig(addr, dbName, dbUser, dbPass, debug, readTimeout, writeTimeout, altHosts...)

	return NewClickhouseConnWithConfig(config)
}

func NewClickhouseConnWithDefault(addr, dbName, dbUser, dbPass string, altHosts ...string) (*Conn, error) {
	config := NewClickhouseConfigWithDefault(addr, dbName, dbUser, dbPass, altHosts...)

	return NewClickhouseConnWithConfig(config)
}

// Execute executes given sql with arguments and return a result
func (conn *Conn) Execute(command string, args ...interface{}) (*Result, error) {
	// prepare
	stmt, err := conn.Prepare(command)
	if err != nil {
		return nil, err
	}

	err = common.SetRandomValueToNil(args...)
	if err != nil {
		return nil, err
	}

	var values []driver.Value

	for _, arg := range args {
		values = append(values, arg)
	}

	sqlType := sqls.GetType(command)
	if sqlType == sqls.Select {
		// this is a select sql
		rows, err := stmt.Query(values)
		if err != nil {
			return nil, err
		}
		defer func() { _ = rows.Close() }()

		return NewResult(rows), nil
	}

	// this is not a select sql
	_, err = stmt.Exec(values)
	if err != nil {
		return nil, err
	}

	return NewEmptyResult(), nil
}

func (conn *Conn) CheckInstanceStatus() bool {
	sql := "select 1 as ok;"
	result, err := conn.Execute(sql)
	if err != nil {
		return false
	}

	ok, err := result.GetIntByName(constant.ZeroInt, "ok")
	if err != nil {
		return false
	}

	return ok == 1
}

func (conn *Conn) Begin() error {
	_, err := conn.Clickhouse.Begin()

	return err
}

func (conn *Conn) Commit() error {
	return conn.Clickhouse.Commit()
}

func (conn *Conn) Rollback() error {
	return conn.Clickhouse.Rollback()
}
