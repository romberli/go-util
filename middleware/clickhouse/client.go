package clickhouse

import (
	"context"
	"database/sql"
	"github.com/romberli/go-util/middleware/sql/statement"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/pingcap/errors"
	"github.com/romberli/go-util/common"
	"github.com/romberli/go-util/constant"
	"github.com/romberli/go-util/middleware"
)

const (
	DefaultDatabase     = "default"
	MaxExecutionTimeStr = "max_execution_time"

	DefaultDialTimeout = 30 * time.Second
	DefaultKeepAlive   = 30 * time.Second

	DefaultReadTimeout  = 60
	DefaultWriteTimeout = 60

	DefaultMaxConnectionsPerConn     = 2
	DefaultMaxIdleConnectionsPerConn = 1
	DefaultMaxLifetime               = time.Hour
	DefaultMaxExecutionTime          = time.Minute
	DefaultBlockBufferSize           = 10

	DefaultCheckInstanceStatusSQL = "select 1 as ok;"
	DefaultGetTimeZoneSQL         = "select timezone();"
)

type Config struct {
	Addr     string
	DBName   string
	DBUser   string
	DBPass   string
	Debug    bool
	AltHosts []string
}

// NewConfig returns a new Config
func NewConfig(addr, dbName, dbUser, dbPass string, debug bool, altHosts ...string) Config {
	return Config{
		Addr:     addr,
		DBName:   dbName,
		DBUser:   dbUser,
		DBPass:   dbPass,
		Debug:    debug,
		AltHosts: altHosts,
	}
}

// NewConfigWithDefault returns a new Config with default value
func NewConfigWithDefault(addr, dbName, dbUser, dbPass string, altHosts ...string) Config {
	return Config{
		Addr:     addr,
		DBName:   dbName,
		DBUser:   dbUser,
		DBPass:   dbPass,
		Debug:    false,
		AltHosts: altHosts,
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

// GetAddrs returns a slice of addresses
func (c *Config) GetAddrs() []string {
	addrs := []string{c.Addr}
	if c.AltHostsExist() {
		addrs = append(addrs, c.AltHosts...)
	}

	return addrs
}

// GetOptions gets *clickhouse.Options
func (c *Config) GetOptions() *clickhouse.Options {
	return &clickhouse.Options{
		Addr: c.GetAddrs(),
		Auth: clickhouse.Auth{
			Database: c.DBName,
			Username: c.DBUser,
			Password: c.DBPass,
		},
		Settings: clickhouse.Settings{
			MaxExecutionTimeStr: DefaultMaxExecutionTime,
		},
		DialTimeout:     DefaultDialTimeout,
		Debug:           c.Debug,
		BlockBufferSize: DefaultBlockBufferSize,
	}
}

type Conn struct {
	Config
	Conn *sql.DB
}

// NewConnWithConfig returns connection to mysql database with given Config
func NewConnWithConfig(config Config) (*Conn, error) {
	conn := clickhouse.OpenDB(config.GetOptions())
	conn.SetMaxOpenConns(DefaultMaxConnectionsPerConn)
	conn.SetMaxIdleConns(DefaultMaxIdleConnectionsPerConn)
	conn.SetConnMaxLifetime(DefaultMaxLifetime)
	conn.SetConnMaxIdleTime(DefaultMaxIdleTime)

	err := conn.Ping()
	if err != nil {
		return nil, errors.Trace(err)
	}

	return &Conn{
		config,
		conn,
	}, nil
}

// NewConn returns connection to Clickhouse database, be aware that addr is host:port style
func NewConn(addr, dbName, dbUser, dbPass string, debug bool, altHosts ...string) (*Conn, error) {
	config := NewConfig(addr, dbName, dbUser, dbPass, debug, altHosts...)

	return NewConnWithConfig(config)
}

func NewConnWithDefault(addr, dbName, dbUser, dbPass string, altHosts ...string) (*Conn, error) {
	config := NewConfigWithDefault(addr, dbName, dbUser, dbPass, altHosts...)

	return NewConnWithConfig(config)
}

// Close closes the connection
func (conn *Conn) Close() error {
	return conn.Conn.Close()
}

// Ping checks if the connection is alive
func (conn *Conn) Ping() error {
	return conn.pingContext(context.Background())
}

// PingContext checks if the connection is alive with context
func (conn *Conn) PingContext(ctx context.Context) error {
	return conn.pingContext(ctx)
}

// pingContext checks if the connection is alive with context
func (conn *Conn) pingContext(ctx context.Context) error {
	return conn.Conn.PingContext(ctx)
}

// Prepare prepares given sql with arguments and return a statement
func (conn *Conn) Prepare(command string) (middleware.Statement, error) {
	return conn.prepareContext(context.Background(), command)
}

// PrepareContext prepares given sql with arguments and return a statement
func (conn *Conn) PrepareContext(ctx context.Context, command string) (middleware.Statement, error) {
	return conn.prepareContext(ctx, command)
}

// prepareContext prepares given sql with arguments and return a statement
func (conn *Conn) prepareContext(ctx context.Context, command string) (middleware.Statement, error) {
	tx, err := conn.Conn.Begin()
	if err != nil {
		return nil, errors.Trace(err)
	}

	stmt, err := tx.PrepareContext(ctx, command)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return NewStatement(stmt, tx, command), nil
}

// Execute executes given sql with arguments and return a result
func (conn *Conn) Execute(command string, args ...interface{}) (middleware.Result, error) {
	return conn.executeContext(context.Background(), command, args...)
}

// ExecuteContext executes given sql with arguments and context then return a result
func (conn *Conn) ExecuteContext(ctx context.Context, command string, args ...interface{}) (middleware.Result, error) {
	return conn.executeContext(ctx, command, args...)
}

// execute executes given sql with arguments and context then return a result
func (conn *Conn) executeContext(ctx context.Context, command string, args ...interface{}) (middleware.Result, error) {
	// set random value to nil
	err := common.SetRandomValueToNil(args...)
	if err != nil {
		return nil, err
	}

	// get sql type
	sqlType := statement.GetType(command)
	if sqlType == statement.Select {
		// this is a select sql
		rows, err := conn.Conn.QueryContext(ctx, command, args...)
		if err != nil {
			return nil, errors.Trace(err)
		}
		defer func() { _ = rows.Close() }()

		return NewResult(rows)
	}
	// this is not a select sql
	_, err = conn.Conn.ExecContext(ctx, command, args...)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return NewEmptyResult(), nil
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
