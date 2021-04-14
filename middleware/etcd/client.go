package etcd

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"go.etcd.io/etcd/clientv3"
)

const (
	DefaultConnectTimeOut    = 10 * time.Second
	DefaultMutexLeaseSeconds = 3600
	MaxTTL                   = 3600 * 24
	MinimumTTL               = 2
	ZeroRevision             = 0
)

type Conn struct {
	Endpoints     []string
	KeyLeaseIDMap sync.Map
	clientv3.Config
	clientv3.Client
	clientv3.Lease
}

// NewEtcdConn returns connection to etcd, it uses client v3 library api
func NewEtcdConn(endpoints []string) (*Conn, error) {
	cfg := clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: DefaultConnectTimeOut,
	}

	client, err := clientv3.New(cfg)
	if err != nil {
		return nil, err
	}

	return &Conn{
		Endpoints: endpoints,
		Config:    cfg,
		Client:    *client,
		Lease:     clientv3.NewLease(client),
	}, nil
}

// NewEtcdConnWithConnectTimeout returns connection to etcd, it uses client v3 library api
func NewEtcdConnWithConnectTimeout(endpoints []string, timeout time.Duration) (*Conn, error) {
	cfg := clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: timeout,
	}

	client, err := clientv3.New(cfg)
	if err != nil {
		return nil, err
	}

	return &Conn{
		Endpoints: endpoints,
		Config:    cfg,
		Client:    *client,
		Lease:     clientv3.NewLease(client),
	}, nil
}

// Close close the connection
func (conn *Conn) Close() error {
	return conn.Client.Close()
}

// GetLeaseIDByKey returns lease response by mutex key name which was maintained when successfully get the mutex
func (conn *Conn) GetLeaseIDByKey(key string) (clientv3.LeaseID, error) {
	leaseID, ok := conn.KeyLeaseIDMap.Load(key)
	if !ok {
		return clientv3.NoLease, nil
	}

	return leaseID.(clientv3.LeaseID), nil
}

func (conn *Conn) CheckLeaseGrantExists(ctx context.Context, leaseID clientv3.LeaseID) (bool, error) {
	if leaseID == clientv3.NoLease {
		return false, nil
	}

	leaseResp, err := conn.TimeToLive(ctx, leaseID)
	if err != nil {
		return false, err
	}

	if leaseResp.TTL <= MinimumTTL {
		return false, nil
	}

	return true, nil
}

// LockEtcdMutex tries to get a distributed mutex from etcd, if success, return true, nil
func (conn *Conn) LockEtcdMutex(ctx context.Context, mutexKey, mutexValue string, ttl int64) (bool, error) {
	if ttl > MaxTTL {
		return false, errors.New(fmt.Sprintf("maximum ttl could not be larger than %d.", MaxTTL))
	}

	leaseResp, err := conn.Grant(ctx, ttl)
	if err != nil {
		return false, err
	}

	txn := clientv3.NewKV(&conn.Client).Txn(ctx)
	txn.If(clientv3.Compare(clientv3.CreateRevision(mutexKey), "=", ZeroRevision)).
		Then(clientv3.OpPut(mutexKey, mutexValue, clientv3.WithLease(leaseResp.ID))).
		Else()
	txnResp, err := txn.Commit()
	if err != nil {
		return false, err
	}

	// successfully get the mutex
	if txnResp.Succeeded {
		conn.KeyLeaseIDMap.Store(mutexKey, leaseResp.ID)
		return true, nil
	}

	return false, nil
}

// UnlockEtcdMutex release the distributed mutex
func (conn *Conn) UnlockEtcdMutex(ctx context.Context, mutexKey string) (*clientv3.LeaseRevokeResponse, error) {
	leaseID, err := conn.GetLeaseIDByKey(mutexKey)
	if err != nil {
		return nil, err
	}

	leaseRevokeResp, err := conn.Revoke(ctx, leaseID)

	conn.KeyLeaseIDMap.Delete(mutexKey)

	return leaseRevokeResp, err
}

// PutWithTTLAndKeepAliveOnce put the key and value and refresh the lease with given ttl once.
func (conn *Conn) PutWithTTLAndKeepAliveOnce(ctx context.Context, key, value string, ttl int64) (*clientv3.PutResponse, *clientv3.LeaseKeepAliveResponse, error) {
	leaseID, err := conn.GetLeaseIDByKey(key)
	if err != nil {
		return nil, nil, err
	}

	leaseGrantExists, err := conn.CheckLeaseGrantExists(ctx, leaseID)
	if err != nil {
		return nil, nil, err
	}

	if !leaseGrantExists {
		leaseResp, err := conn.Grant(ctx, ttl)
		if err != nil {
			return nil, nil, err
		}

		leaseID = leaseResp.ID
		conn.KeyLeaseIDMap.Store(key, leaseID)
	}

	putResp, err := conn.Client.Put(ctx, key, value, clientv3.WithLease(leaseID))
	if err != nil {
		return nil, nil, err
	}

	leaseKeepAliveResp, err := conn.KeepAliveOnce(ctx, leaseID)

	return putResp, leaseKeepAliveResp, err
}

// PutWithTTLAndKeepAlive put the key and value and keep alive the lease with given ttl
func (conn *Conn) PutWithTTLAndKeepAlive(ctx context.Context, key, value string, ttl int64) (*clientv3.PutResponse, <-chan *clientv3.LeaseKeepAliveResponse, error) {
	putResp, leaseKeepAliveResp, err := conn.PutWithTTLAndKeepAliveOnce(ctx, key, value, ttl)
	if err != nil {
		return nil, nil, err
	}

	leaseKeepAliveRespChan, err := conn.KeepAlive(ctx, leaseKeepAliveResp.ID)

	return putResp, leaseKeepAliveRespChan, err
}

// Delete delete the key, but does NOT revoke the concerned lease because the lease may be assigned to other keys.
func (conn *Conn) Delete(ctx context.Context, key string) (*clientv3.DeleteResponse, error) {
	leaseID, err := conn.GetLeaseIDByKey(key)
	if err != nil {
		return nil, err
	}

	if leaseID == clientv3.NoLease {
		return conn.Client.Delete(ctx, key)
	}

	conn.KeyLeaseIDMap.Delete(key)

	return conn.Client.Delete(ctx, key)
}

// DeleteWithPrefix delete the keys of which name started with the given key,
// but does NOT revoke the concerned lease because the lease may be assigned to other keys.
func (conn *Conn) DeleteWithPrefix(ctx context.Context, key string) (*clientv3.DeleteResponse, error) {
	leaseID, err := conn.GetLeaseIDByKey(key)
	if err != nil {
		return nil, err
	}

	if leaseID == clientv3.NoLease {
		return conn.Client.Delete(ctx, key)
	}

	conn.KeyLeaseIDMap.Delete(key)

	return conn.Client.Delete(ctx, key, clientv3.WithPrefix())
}

// PutWithTTL is an alias of PutWithTTLAndKeepAliveOnce
func (conn *Conn) PutWithTTL(ctx context.Context, key, value string, ttl int64) (*clientv3.PutResponse, *clientv3.LeaseKeepAliveResponse, error) {
	return conn.PutWithTTLAndKeepAliveOnce(ctx, key, value, ttl)
}
