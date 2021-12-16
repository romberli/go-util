package coredns

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEtcdConnection(t *testing.T) {
	const (
		DefaultPath         = "/coredns"
		DefaultURLPrefix    = "test001.example.com"
		DefaultMasterURL    = "master.test001.example.com"
		DefaultSlaveURL     = "slave.test001.example.com"
		DefaultMasterHostIP = "192.168.137.11"
		DefaultSlaveHostIP0 = "192.168.137.11"
		DefaultSlaveHostIP1 = "192.168.137.11"
		DefaultSleepTime    = 5 * time.Second
		DefaultKey          = "key001"
		DefaultValue        = "value001"
		DefaultTTL          = 3
		DefaultKeyNotFound  = 0
		DefaultMutexKey     = "MutexKey001"
		DefaultMutexValue   = "MutexValue001"
	)

	var (
		err  error
		conn *Conn
	)

	assert := assert.New(t)

	endpoints := []string{"192.168.137.11:2379"}
	ctx := context.Background()

	conn, err = NewCoreDNSConn(endpoints, DefaultPath)
	assert.Nil(err, "connect to etcd failed. endpoints: %v", endpoints)
	defer func() {
		err = conn.Close()
		assert.Nil(err, "close connection failed.")
	}()

	// no A record
	ipList, err := conn.Resolve(ctx, DefaultMasterURL)
	assert.Nil(err, fmt.Sprintf("got error when resolving %s", DefaultMasterURL))
	assert.Nil(ipList, "there should not exist A record for now.")

	m0 := "0." + DefaultMasterURL
	// add master A record
	err = conn.PutARecord(ctx, m0, DefaultMasterHostIP, DefaultTTL)
	assert.Nil(err, fmt.Sprintf("got error when putting a record of %s", DefaultMasterURL))

	ipList, err = conn.Resolve(ctx, DefaultMasterURL)
	assert.Nil(err, fmt.Sprintf("got error when resolving %s", DefaultMasterURL))
	assert.Equal([]string{DefaultMasterHostIP}, ipList, "there should not exist master A record now.")

	// master A record expired
	time.Sleep(DefaultSleepTime)
	ipList, err = conn.Resolve(ctx, DefaultMasterURL)
	assert.Nil(err, fmt.Sprintf("got error when resolving %s", DefaultMasterURL))
	assert.Nil(ipList, "master A record should expire.")

	// add master A record and keep alive
	err = conn.PutARecordAndKeepAlive(ctx, m0, DefaultMasterHostIP, DefaultTTL)
	assert.Nil(err, fmt.Sprintf("got error when putting a record of %s", DefaultMasterURL))
	time.Sleep(DefaultSleepTime)
	ipList, err = conn.Resolve(ctx, DefaultMasterURL)
	assert.Nil(err, fmt.Sprintf("got error when resolving %s", DefaultMasterURL))
	assert.Equal([]string{DefaultMasterHostIP}, ipList, "master A record should still exist.")

	// delete master A record
	err = conn.DeleteARecord(ctx, m0)
	assert.Nil(err, fmt.Sprintf("got error when deleting a record of %s", DefaultMasterURL))
	ipList, err = conn.Resolve(ctx, DefaultMasterURL)
	assert.Nil(err, fmt.Sprintf("got error when resolving %s", DefaultMasterURL))
	assert.Nil(ipList, "master A record should be deleted.")

	// add slave A record
	s0 := "0." + DefaultSlaveURL
	err = conn.PutARecord(ctx, s0, DefaultSlaveHostIP0, DefaultTTL)
	assert.Nil(err, fmt.Sprintf("got error when putting a record of %s", DefaultSlaveURL))
	s1 := "1." + DefaultSlaveURL
	err = conn.PutARecord(ctx, s1, DefaultSlaveHostIP1, DefaultTTL)
	assert.Nil(err, fmt.Sprintf("got error when putting a record of %s", DefaultSlaveURL))
	ipList, err = conn.Resolve(ctx, DefaultSlaveURL)
	assert.Nil(err, fmt.Sprintf("got error when resolving %s", DefaultSlaveURL))
	assert.Equal(len(ipList), 2, "slave A record should exist 2 records.")

	// delete slave A record
	err = conn.DeleteARecord(ctx, s1)
	assert.Nil(err, fmt.Sprintf("got error when deleting a record of %s", DefaultSlaveURL))
	ipList, err = conn.Resolve(ctx, DefaultSlaveURL)
	assert.Nil(err, fmt.Sprintf("got error when resolving %s", DefaultSlaveURL))
	assert.Equal(len(ipList), 1, "slave A record should exist 1 records.")
}
