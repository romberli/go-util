package mysql

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-mysql-org/go-mysql/client"
	"github.com/pingcap/errors"
	"github.com/romberli/go-util/common"
	"github.com/romberli/go-util/constant"
)

type ReplicationRole string

const (
	DefaultAddrSliceLen = 2
	ValueStr            = "Value"
	VariableOn          = "ON"
	VariableOFF         = "OFF"

	DefaultCharSet                = "utf8mb4"
	HostString                    = "host"
	PortString                    = "port"
	DefaultCheckInstanceStatusSQL = "select 1 as ok;"
	SelectVersionSQL              = "select @@version"
	ShowSlaveStatusSQL            = "show slave status"
	ShowReplicaStatusSQL          = "show replica status"
	ShowSlaveHostsSQL             = "show slave hosts"

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

func (c Config) GetAddr() string {
	return c.Addr
}

func (c Config) GetDBName() string {
	return c.DBName
}

func (c Config) GetDBUser() string {
	return c.DBUser
}

func (c Config) GetDBPass() string {
	return c.DBPass
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
		return nil, errors.Trace(err)
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
		return nil, errors.Trace(err)
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

	if result.RowNumber() > constant.ZeroInt {
		role = ReplicationReplica
	}

	slaveList, err := conn.GetReplicationSlaveList()
	if err != nil {
		return constant.EmptyString, err
	}

	if len(slaveList) > constant.ZeroInt && role == ReplicationReplica {
		role = ReplicationRelay
	}

	return role, nil
}

// IsMaster checks if this server is master
func (conn *Conn) IsMaster() (bool, error) {
	isSlave, err := conn.IsReplicationSlave()
	if err != nil {
		return false, err
	}

	if isSlave {
		return false, nil
	}

	isMGR, err := conn.IsMGR()
	if err != nil {
		return false, err
	}
	if isMGR {
		mysqlVersion, err := conn.GetVersion()
		if err != nil {
			return false, err
		}

		addrSlice := strings.Split(conn.GetAddr(), constant.ColonString)
		if len(addrSlice) != DefaultAddrSliceLen {
			return false, errors.Errorf("connection address must be formatted as ip:port, %s is not valid", conn.GetAddr())
		}
		hostIP := addrSlice[constant.ZeroInt]
		portNum := addrSlice[constant.OneInt]

		var sql string

		if mysqlVersion.IsMySQL8() {
			// mysql 8.0
			sql = "SELECT member_role FROM performance_schema.replication_group_members WHERE member_host = ? AND member_port = ? AND member_role = 'SECONDARY';"
			result, err := conn.Execute(sql, hostIP, portNum)
			if err != nil {
				return false, err
			}
			if result.RowNumber() > constant.ZeroInt {
				return false, nil
			}
		} else {
			// mysql 5.7
			// check if single primary mode
			sql = "SHOW STATUS LIKE 'group_replication_primary_member' ;"
			result, err := conn.Execute(sql)
			if err != nil {
				return false, err
			}

			if result.RowNumber() == constant.ZeroInt {
				// maybe multiple primary mode or no primary node exists
				readOnly, err := conn.IsReadOnly()
				if err != nil {
					return false, err
				}

				return !readOnly, nil
			}

			// single primary mode
			primaryMemberUUID, err := result.GetString(constant.ZeroInt, constant.OneInt)
			if err != nil {
				return false, err
			}
			sql = "SELECT member_id FROM performance_schema.replication_group_members WHERE member_id = ? AND member_host = ? AND member_port = ? ;"
			result, err = conn.Execute(sql, primaryMemberUUID, hostIP, portNum)
			if err != nil {
				return false, err
			}

			return result.RowNumber() > constant.ZeroInt, nil
		}
	}

	return true, nil
}

// IsReplicationSlave checks if this server is slave
func (conn *Conn) IsReplicationSlave() (bool, error) {
	result, err := conn.Execute(ShowSlaveStatusSQL)
	if err != nil {
		return false, err
	}

	return result.RowNumber() > constant.ZeroInt, nil
}

// IsMGR checks if this server is a member of MGR cluster
func (conn *Conn) IsMGR() (bool, error) {
	sql := "SELECT COUNT(*) FROM performance_schema.replication_group_members ;"
	result, err := conn.Execute(sql)
	if err != nil {
		return false, err
	}

	count, err := result.GetInt(constant.ZeroInt, constant.ZeroInt)
	if err != nil {
		return false, err
	}

	if count > constant.ZeroInt {
		return true, nil
	}

	return false, nil
}

// IsReadOnly checks if this server is read only
func (conn *Conn) IsReadOnly() (bool, error) {
	sql := "SHOW VARIABLES LIKE 'read_only';"
	result, err := conn.Execute(sql)
	if err != nil {
		return false, err
	}
	status, err := result.GetString(constant.ZeroInt, constant.OneInt)
	if err != nil {
		return false, err
	}

	return status == VariableOn, nil
}

// IsSuperReadOnly checks if this server is super read only
func (conn *Conn) IsSuperReadOnly() (bool, error) {
	sql := "SHOW VARIABLES LIKE 'super_read_only';"
	result, err := conn.Execute(sql)
	if err != nil {
		return false, err
	}
	status, err := result.GetString(constant.ZeroInt, constant.OneInt)
	if err != nil {
		return false, err
	}

	return status == VariableOn, nil
}

// SetReadOnly sets read only
func (conn *Conn) SetReadOnly(readOnly bool) error {
	var value int

	if readOnly {
		value = constant.OneInt
	}
	sql := fmt.Sprintf("SET GLOBAL read_only = %d;", value)
	_, err := conn.Execute(sql)

	return err
}

// SetSuperReadOnly sets super read only
func (conn *Conn) SetSuperReadOnly(superReadOnly bool) error {
	var value int

	if superReadOnly {
		value = constant.OneInt
	}
	sql := fmt.Sprintf("SET GLOBAL super_read_only = %d;", value)
	_, err := conn.Execute(sql)

	return err
}

// ShowGlobalVariable returns the value of the given variable
func (conn *Conn) ShowGlobalVariable(variable string) (string, error) {
	sql := fmt.Sprintf("SHOW GLOBAL VARIABLES LIKE '%s';", variable)
	result, err := conn.Execute(sql)
	if err != nil {
		return constant.EmptyString, err
	}

	return result.GetString(constant.ZeroInt, constant.OneInt)
}

// SetGlobalVariable sets global variable
func (conn *Conn) SetGlobalVariable(variable, value string) error {
	sql := fmt.Sprintf("SET GLOBAL %s = '%s';", variable, value)
	_, err := conn.Execute(sql)

	return err
}

// SetGlobalVariables sets global variables
func (conn *Conn) SetGlobalVariables(variables map[string]string) error {
	for variable, value := range variables {
		err := conn.SetGlobalVariable(variable, value)
		if err != nil {
			return err
		}
	}

	return nil
}
