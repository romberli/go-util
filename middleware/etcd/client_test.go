package etcd

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEtcdConnection(t *testing.T) {
	const (
		DefaultKey        = "key001"
		DefaultValue      = "value001"
		DefaultMutexKey   = "MutexKey001"
		DefaultMutexValue = "MutexValue001"
	)

	var (
		err  error
		ok   bool
		conn *EtcdConn
	)

	assert := assert.New(t)

	endpoints := []string{"192.168.137.11:2379"}
	ctx := context.Background()

	conn, err = NewEtcdConn(endpoints)
	assert.Nil(err, "connect to etcd failed. endpoints: %v", endpoints)
	defer func() {
		err = conn.Close()
		assert.Nil(err, "close connection failed.")
	}()

	_, err = conn.Put(ctx, DefaultKey, DefaultValue)
	assert.Nil(err, "put key to etcd failed.")

	GetResp, err := conn.Get(ctx, DefaultKey)
	assert.Equal(string(GetResp.Kvs[0].Key), DefaultKey, "put key to etcd failed.")
	assert.Equal(string(GetResp.Kvs[0].Value), DefaultValue, "put key to etcd failed.")

	ok, err = conn.LockEtcdMutex(ctx, DefaultMutexKey, DefaultMutexValue, DefaultMutexLeaseSeconds)
	assert.Nil(err, "got error when trying to get mutex from etcd.")
	assert.True(ok, "this is the first time to get mutex and should success.")

	ok, err = conn.LockEtcdMutex(ctx, DefaultMutexKey, DefaultMutexValue, DefaultMutexLeaseSeconds)
	assert.Nil(err, "got error when trying to lock the mutex from etcd.")
	assert.False(ok, "this is the second time to get mutex and should success.")

	err = conn.UnlockEtcdMutex(ctx, DefaultMutexKey)
	assert.Nil(err, "got error when trying to unlock the mutex from etcd")
}
