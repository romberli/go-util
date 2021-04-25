package mysql

import (
	"context"
	"fmt"

	"github.com/go-mysql-org/go-mysql/client"

	"github.com/romberli/go-util/common"
	"github.com/romberli/go-util/constant"
)

type ReplicationRole string

const (
	DefaultCharSet       = "utf8mb4"
	HostString           = "host"
	PortString           = "port"
	SelectVersionSQL     = "select @@version"
	ShowSlaveStatusSQL   = "show slave status"
	ShowReplicaStatusSQL = "show replica status"
	ShowSlaveHostsSQL    = "show slave hosts"

	// ReplicationSource represents mysql master, it's an alternative name
	ReplicationSource ReplicationRole = "source"
	// ReplicationReplica represents mysql slave, it's an alternative name
	ReplicationReplica ReplicationRole = "replica"
	// ReplicationRelay means this mysql instance has source and replica roles at the same time
	ReplicationRelay ReplicationRole = "relay"
)

type Config struct {
	Addr   string
	DBName string
	DBUser string
	DBPass string
}

// NewConfig returns a new Config
func NewConfig(addr string, dbName string, dbUser string, dbPass string) Config {
	return Config{
		Addr:   addr,
		DBName: dbName,
		DBUser: dbUser,
		DBPass: dbPass,
	}
}

type Conn struct {
	Config
	*client.Conn
}

// NewConn returns connection to mysql database, be aware that addr is host:port style, default charset is utf8mb4
func NewConn(addr string, dbName string, dbUser string, dbPass string) (*Conn, error) {
	config := NewConfig(addr, dbName, dbUser, dbPass)

	// connect to mysql
	conn, err := client.Connect(addr, dbUser, dbPass, dbName)
	if err != nil {
		return nil, err
	}

	// set connection charset
	err = conn.SetCharset(DefaultCharSet)
	if err != nil {
		return nil, err
	}

	// use db
	err = conn.UseDB(dbName)
	if err != nil {
		return nil, err
	}

	return &Conn{
		Config: config,
		Conn:   conn,
	}, nil
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
	stmt, err := conn.Conn.Prepare(command)
	if err != nil {
		return nil, err
	}

	return NewStatement(stmt), nil
}

// Execute executes given sql and placeholders and returns a result
func (conn *Conn) Execute(command string, args ...interface{}) (*Result, error) {
	return conn.executeContext(context.Background(), command, args...)
}

// ExecuteContext executes given sql and placeholders with context and returns a result
func (conn *Conn) ExecuteContext(ctx context.Context, command string, args ...interface{}) (*Result, error) {
	return conn.executeContext(ctx, command, args...)
}

// executeContext executes given sql and placeholders with context and returns a result
func (conn *Conn) executeContext(ctx context.Context, command string, args ...interface{}) (*Result, error) {
	err := common.SetRandomValueToNil(args...)
	if err != nil {
		return nil, err
	}

	result, err := conn.Conn.Execute(command, args...)
	if err != nil {
		return nil, err
	}

	return NewResult(result), nil
}

// GetVersion returns mysql version
func (conn *Conn) GetVersion() (Version, error) {
	result, err := conn.Execute(SelectVersionSQL)
	if err != nil {
		return nil, err
	}

	versionStr, err := result.GetString(constant.ZeroInt, constant.ZeroInt)
	if err != nil {
		return nil, err
	}

	return Parse(versionStr)
}

// CheckInstanceStatus checks mysql instance status
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

// GetReplicationSlaveList returns slave list of this server
func (conn *Conn) GetReplicationSlaveList() (slaveList []string, err error) {
	slaveList = []string{}

	result, err := conn.Execute(ShowSlaveHostsSQL)
	if err != nil {
		return nil, err
	}

	for i := 0; i < result.RowNumber(); i++ {
		host, err := result.GetStringByName(i, HostString)
		if err != nil {
			return nil, err
		}

		port, err := result.GetIntByName(i, PortString)
		if err != nil {
			return nil, err
		}

		addr := fmt.Sprintf("%s:%d", host, port)
		slaveList = append(slaveList, addr)
	}

	return slaveList, nil
}

// GetReplicationSlavesStatus returns replication slave status, like sql: "show slave status;"
func (conn *Conn) GetReplicationSlavesStatus() (result *Result, err error) {
	return conn.executeContext(context.Background(), ShowSlaveStatusSQL)
}

// GetReplicationRole returns replication role
func (conn *Conn) GetReplicationRole() (role ReplicationRole, err error) {
	role = ReplicationSource

	result, err := conn.GetReplicationSlavesStatus()
	if err != nil {
		return constant.EmptyString, err
	}

	if result.RowNumber() != 0 {
		role = ReplicationReplica
	}

	slaveList, err := conn.GetReplicationSlaveList()
	if err != nil {
		return constant.EmptyString, err
	}

	if len(slaveList) != 0 && role == ReplicationReplica {
		role = ReplicationRelay
	}

	return role, nil
}
