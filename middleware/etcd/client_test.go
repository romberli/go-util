package etcd

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEtcdConnection(t *testing.T) {
	const (
		DefaultKey         = "key001"
		DefaultValue       = "value001"
		DefaultTTL         = 2
		DefaultKeyNotFound = 0
		DefaultMutexKey    = "MutexKey001"
		DefaultMutexValue  = "MutexValue001"
	)

	var (
		err  error
		ok   bool
		conn *Conn
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
	assert.Nil(err, "put key failed.")

	GetResp, err := conn.Get(ctx, DefaultKey)
	assert.Equal(string(GetResp.Kvs[0].Key), DefaultKey, "get key failed.")
	assert.Equal(string(GetResp.Kvs[0].Value), DefaultValue, "get key failed.")

	_, leaseKeepAliveResp, err := conn.PutWithTTL(ctx, DefaultKey, DefaultValue, DefaultTTL)
	assert.Nil(err, "put key with ttl failed.")
	assert.Equal(leaseKeepAliveResp.TTL, int64(DefaultTTL), "put key with ttl failed.")

	time.Sleep((DefaultTTL + 1) * time.Second)

	GetResp, err = conn.Get(ctx, DefaultKey)
	assert.Equal(GetResp.Count, int64(DefaultKeyNotFound), "get expired key failed.")

	ok, err = conn.LockEtcdMutex(ctx, DefaultMutexKey, DefaultMutexValue, DefaultMutexLeaseSeconds)
	assert.Nil(err, "got error when trying to get mutex.")
	assert.True(ok, "this is the first time to get mutex and should success.")

	ok, err = conn.LockEtcdMutex(ctx, DefaultMutexKey, DefaultMutexValue, DefaultMutexLeaseSeconds)
	assert.Nil(err, "got error when trying to lock the mutex from etcd.")
	assert.False(ok, "this is the second time to get mutex and should success.")

	_, err = conn.UnlockEtcdMutex(ctx, DefaultMutexKey)
	assert.Nil(err, "got error when trying to unlock the mutex")
}
