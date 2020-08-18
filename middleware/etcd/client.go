package etcd

import (
	"context"
	"time"

	"go.etcd.io/etcd/clientv3"

	"github.com/romber2001/go-util/common"
)

const (
	DefaultConnectTimeOut    = 10 * time.Second
	DefaultMutexLeaseSeconds = 3600
	DefaultLeaseID           = 0
)

type EtcdConn struct {
	Endpoints       []string
	Cfg             clientv3.Config
	KeyLeaseMap     map[string]clientv3.Lease
	KeyLeaseRespMap map[string]*clientv3.LeaseGrantResponse
	clientv3.Client
}

// NewEtcdConn returns connection to etcd, it uses client v3 library api
func NewEtcdConn(endpoints []string) (*EtcdConn, error) {
	cfg := clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: DefaultConnectTimeOut,
	}

	client, err := clientv3.New(cfg)
	if err != nil {
		return nil, err
	}

	return &EtcdConn{
		Endpoints:       endpoints,
		Cfg:             cfg,
		KeyLeaseMap:     make(map[string]clientv3.Lease),
		KeyLeaseRespMap: make(map[string]*clientv3.LeaseGrantResponse),
		Client:          *client,
	}, nil
}

// Close close the etcd connection
func (conn *EtcdConn) Close() error {
	return conn.Client.Close()
}

// NewLease returns lease
func (conn *EtcdConn) NewLease() clientv3.Lease {
	return clientv3.NewLease(&conn.Client)
}

// NewLease returns lease grant response which contains lease id
func (conn *EtcdConn) NewLeaseGrantResponse(ctx context.Context, lease clientv3.Lease, ttl int64) (*clientv3.LeaseGrantResponse, error) {
	leaseResp, err := lease.Grant(ctx, ttl)
	if err != nil {
		return nil, err
	}

	return leaseResp, err
}

// GetLeaseByKey returns lease by mutex key name which was maintained when successfully get the mutex
func (conn *EtcdConn) GetLeaseByKey(key string) (clientv3.Lease, error) {
	keyExists, err := common.KeyInMap(key, conn.KeyLeaseMap)
	if err != nil {
		return nil, err
	}

	if keyExists {
		return conn.KeyLeaseMap[key], nil
	}

	return nil, nil
}

// GetLeaseRespByKey returns lease response by mutex key name which was maintained when successfully get the mutex
func (conn *EtcdConn) GetLeaseRespByKey(key string) (*clientv3.LeaseGrantResponse, error) {
	keyExists, err := common.KeyInMap(key, conn.KeyLeaseRespMap)
	if err != nil {
		return nil, err
	}

	if keyExists {
		return conn.KeyLeaseRespMap[key], nil
	}

	return nil, nil
}

// LockEtcdMutex tries to get a distributed mutex from etcd, if success, return true, nil
func (conn *EtcdConn) LockEtcdMutex(ctx context.Context, mutexKey string, mutexValue string, ttl int64) (bool, error) {
	lease := conn.NewLease()
	leaseResp, err := conn.NewLeaseGrantResponse(ctx, lease, ttl)
	if err != nil {
		return false, err
	}

	txn := clientv3.NewKV(&conn.Client).Txn(ctx)
	txn.If(clientv3.Compare(clientv3.CreateRevision(mutexKey), "=", 0)).
		Then(clientv3.OpPut(mutexKey, mutexValue, clientv3.WithLease(leaseResp.ID))).
		Else()
	txnResp, err := txn.Commit()
	if err != nil {
		return false, err
	}

	// successfully get the mutex
	if txnResp.Succeeded {
		conn.KeyLeaseMap[mutexKey] = lease
		conn.KeyLeaseRespMap[mutexKey] = leaseResp
		return true, nil
	}

	return false, nil
}

// UnlockEtcdMutex release the distributed mutex
func (conn *EtcdConn) UnlockEtcdMutex(ctx context.Context, mutexKey string) error {
	lease, err := conn.GetLeaseByKey(mutexKey)
	if err != nil {
		return err
	}

	if lease == nil {
		// lease does not exist in the key map, this means the connection did not get the mutex successfully before
		return nil
	}

	leaseResp, err := conn.GetLeaseRespByKey(mutexKey)
	if err != nil {
		return err
	}

	if leaseResp == nil {
		// leaseResp does not exist in the key map, this means the connection did not get the mutex successfully before
		return nil
	}

	_, err = lease.Revoke(ctx, leaseResp.ID)
	delete(conn.KeyLeaseMap, mutexKey)
	delete(conn.KeyLeaseRespMap, mutexKey)

	return err
}
