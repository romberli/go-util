package mysql

import (
	"fmt"

	"github.com/siddontang/go-mysql/client"
	"github.com/siddontang/go-mysql/mysql"

	"github.com/romberli/go-util/constant"
)

const (
	DefaultCharSet     = "utf8mb4"
	ReplicationMaster  = "master"
	ReplicationSlave   = "slave"
	ReplicationRelay   = "relay" // it has master and slave roles at the same time
	HostString         = "host"
	PortString         = "port"
	ShowSlaveStatusSQL = "show slave status"
)

type Config struct {
	Addr   string
	DBName string
	DBUser string
	DBPass string
}

func NewMySQLConfig(addr string, dbName string, dbUser string, dbPass string) Config {
	return Config{
		Addr:   addr,
		DBName: dbName,
		DBUser: dbUser,
		DBPass: dbPass,
	}
}

type Conn struct {
	Config
	client.Conn
}

// NewMySQLConn returns connection to mysql database, be aware that addr is host:port style, default charset is utf8mb4
func NewMySQLConn(addr string, dbName string, dbUser string, dbPass string) (*Conn, error) {
	var (
		err  error
		conn *client.Conn
	)

	config := NewMySQLConfig(addr, dbName, dbUser, dbPass)

	// connect to mysql
	conn, err = client.Connect(addr, dbUser, dbPass, dbName)
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
		Conn:   *conn,
	}, nil
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

func (conn *Conn) GetReplicationSlaveList() (slaveList []string, err error) {
	var (
		result *mysql.Result
		host   string
		port   int64
	)

	slaveList = []string{}

	result, err = conn.Execute(ShowSlaveStatusSQL)
	if err != nil {
		return nil, err
	}

	for i := 0; i < result.RowNumber(); i++ {
		host, err = result.GetStringByName(i, HostString)
		if err != nil {
			return nil, err
		}

		port, err = result.GetIntByName(i, PortString)
		if err != nil {
			return nil, err
		}

		addr := fmt.Sprintf("%s:%d", host, port)
		slaveList = append(slaveList, addr)
	}

	return slaveList, nil
}

func (conn *Conn) GetReplicationSlavesStatus() (result *mysql.Result, err error) {
	return conn.Execute(ShowSlaveStatusSQL)
}

func (conn *Conn) GetReplicationRole() (role string, err error) {
	var (
		slaveList []string
		result    *mysql.Result
	)

	role = ReplicationMaster

	result, err = conn.GetReplicationSlavesStatus()
	if err != nil {
		return "", err
	}

	if result.RowNumber() != 0 {
		role = ReplicationSlave
	}

	slaveList, err = conn.GetReplicationSlaveList()
	if err != nil {
		return "", err
	}

	if len(slaveList) != 0 && role == ReplicationSlave {
		role = ReplicationRelay
	}

	return role, nil
}
