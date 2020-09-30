package mysql

import (
	"fmt"

	"github.com/siddontang/go-mysql/client"
	"github.com/siddontang/go-mysql/mysql"
)

const (
	DefaultCharSet    = "utf8mb4"
	ReplicationMaster = "master"
	ReplicationSlave  = "slave"
	ReplicationRelay  = "relay" // it has master and slave roles at the same time
)

type Conn struct {
	Addr   string
	DBName string
	DBUser string
	DBPass string
	client.Conn
}

// NewMySQLConn returns connection to mysql database, be aware that addr is host:port style, default charset is utf8mb4
func NewMySQLConn(addr string, dbName string, dbUser string, dbPass string) (*Conn, error) {
	var (
		err  error
		conn *client.Conn
	)

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
		Addr:   addr,
		DBName: dbName,
		DBUser: dbUser,
		DBPass: dbPass,
		Conn:   *conn,
	}, nil
}

func (conn *Conn) GetReplicationSlaveList() (slaveList []string, err error) {
	var (
		result *mysql.Result
		host   string
		port   int64
	)

	slaveList = []string{}

	result, err = conn.Execute("show slave hosts ;")
	if err != nil {
		return nil, err
	}

	for i := 0; i < result.RowNumber(); i++ {
		host, err = result.GetStringByName(i, "Host")
		if err != nil {
			return nil, err
		}

		port, err = result.GetIntByName(i, "Port")
		if err != nil {
			return nil, err
		}

		addr := fmt.Sprintf("%s:%d", host, port)
		slaveList = append(slaveList, addr)
	}

	return slaveList, nil
}

func (conn *Conn) GetReplicationSlavesStatus() (result *mysql.Result, err error) {
	result, err = conn.Execute("show slave status ;")
	if err != nil {
		return nil, err
	}

	return result, nil
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
