package clickhouse

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pingcap/errors"

	"github.com/ClickHouse/clickhouse-go"
	"github.com/romberli/go-util/common"
	"github.com/romberli/go-util/constant"
)

const (
	DefaultDatabase     = "default"
	DefaultReadTimeout  = 10
	DefaultWriteTimeout = 10

	DefaultCheckInstanceStatusSQL = "select 1 as ok;"
	DefaultGetTimeZoneSQL         = "select timezone();"
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

// NewConfig returns a new Config
func NewConfig(addr, dbName, dbUser, dbPass string, debug bool, readTimeout, writeTimeout int, altHosts ...string) Config {
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

// NewConfigWithDefault returns a new Config with default value
func NewConfigWithDefault(addr, dbName, dbUser, dbPass string, altHosts ...string) Config {
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

// NewConnWithConfig returns connection to mysql database with given Config
func NewConnWithConfig(config Config) (*Conn, error) {
	// connect to Clickhouse
	client, err := clickhouse.OpenDirect(config.GetConnectionString())
	if err != nil {
		return nil, errors.Trace(err)
	}

	return &Conn{
		config,
		client,
	}, nil
}

// NewConn returns connection to Clickhouse database, be aware that addr is host:port style
func NewConn(addr, dbName, dbUser, dbPass string, debug bool, readTimeout, writeTimeout int, altHosts ...string) (*Conn, error) {
	config := NewConfig(addr, dbName, dbUser, dbPass, debug, readTimeout, writeTimeout, altHosts...)

	return NewConnWithConfig(config)
}

func NewConnWithDefault(addr, dbName, dbUser, dbPass string, altHosts ...string) (*Conn, error) {
	config := NewConfigWithDefault(addr, dbName, dbUser, dbPass, altHosts...)

	return NewConnWithConfig(config)
}

// Prepare prepares a statement and returns a *Statement
func (conn *Conn) Prepare(command string) (*Statement, error) {
	return conn.prepareContext(context.Background(), command)
}

// PrepareContext prepares a statement with context and returns a *Statement
func (conn *Conn) PrepareContext(ctx context.Context, command string) (*Statement, error) {
	return conn.prepareContext(ctx, command)
}

// prepareContext prepares a statement with context and returns a *Statement
func (conn *Conn) prepareContext(ctx context.Context, command string) (*Statement, error) {
	stmt, err := conn.Clickhouse.PrepareContext(ctx, command)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return NewStatement(stmt), nil
}

// Execute executes given sql with arguments and return a result
func (conn *Conn) Execute(command string, args ...interface{}) (*Result, error) {
	return conn.executeContext(context.Background(), command, args...)
}

// ExecuteContext executes given sql with arguments and context then return a result
func (conn *Conn) ExecuteContext(ctx context.Context, command string, args ...interface{}) (*Result, error) {
	return conn.executeContext(ctx, command, args...)
}

// execute executes given sql with arguments and context then return a result
func (conn *Conn) executeContext(ctx context.Context, command string, args ...interface{}) (*Result, error) {
	// prepare
	stmt, err := conn.prepareContext(ctx, command)
	if err != nil {
		return nil, err
	}
	// set random value to nil
	err = common.SetRandomValueToNil(args...)
	if err != nil {
		return nil, err
	}

	return stmt.executeContext(ctx, args...)
}

// CheckInstanceStatus returns if instance is ok
func (conn *Conn) CheckInstanceStatus() bool {
	result, err := conn.Execute(DefaultCheckInstanceStatusSQL)
	if err != nil {
		return false
	}

	ok, err := result.GetInt(constant.ZeroInt, constant.ZeroInt)
	if err != nil {
		return false
	}

	return ok == 1
}

func (conn *Conn) GetTimeZone() (*time.Location, error) {
	result, err := conn.Execute(DefaultGetTimeZoneSQL)
	if err != nil {
		return nil, err
	}

	tz, err := result.GetString(constant.ZeroInt, constant.ZeroInt)
	if err != nil {
		return nil, err
	}

	t, err := time.LoadLocation(tz)
	if err != nil {
		return nil, err
	}

	return t, nil
}

// Begin begins a new transaction
func (conn *Conn) Begin() error {
	_, err := conn.Clickhouse.Begin()

	return errors.Trace(err)
}

// Commit commits the transaction
func (conn *Conn) Commit() error {
	return errors.Trace(conn.Clickhouse.Commit())
}

// Rollback rollbacks the transaction
func (conn *Conn) Rollback() error {
	return errors.Trace(conn.Clickhouse.Rollback())
}
