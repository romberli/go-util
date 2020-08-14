package etcd

import (
	"context"
	"time"

	"github.com/etcd-io/etcd/clientv3"

	"github.com/romber2001/go-util/common"
)

const (
	DefaultConnectTimeOut    = 10 * time.Second
	DefaultMutexLeaseSeconds = 3600 * time.Second
	DefaultLeaseID           = 0
)

type EtcdConn struct {
	Endpoints       []string
	Cfg             clientv3.Config
	KeyLeaseMap     map[string]clientv3.Lease
	KeyLeaseRespMap map[string]*clientv3.LeaseGrantResponse
	clientv3.Client
}

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
		Endpoints: endpoints,
		Cfg:       cfg,
		Client:    *client,
	}, nil
}

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

func (conn *EtcdConn) NewLease() clientv3.Lease {
	return clientv3.NewLease(&conn.Client)
}

func (conn *EtcdConn) GetLeaseGrantResponse(ctx context.Context, lease clientv3.Lease, ttl int64) (*clientv3.LeaseGrantResponse, error) {
	leaseResp, err := lease.Grant(ctx, ttl)
	if err != nil {
		return nil, err
	}

	return leaseResp, err
}

func (conn *EtcdConn) LockEtcdMutex(ctx context.Context, mutexKey string, ttl int64) (bool, error) {
	lease := conn.NewLease()
	leaseResp, err := conn.GetLeaseGrantResponse(ctx, lease, ttl)
	if err != nil {
		return false, err
	}

	txn := clientv3.NewKV(&conn.Client).Txn(ctx)
	txn.If(clientv3.Compare(clientv3.CreateRevision(mutexKey), "=", 0)).
		Then(clientv3.OpPut(mutexKey, "", clientv3.WithLease(leaseResp.ID))).
		Else()
	txnResp, err := txn.Commit()
	if err != nil {
		return false, err
	}

	if txnResp.Succeeded {
		return true, nil
	}

	return false, nil
}

func (conn *EtcdConn) UnlockEtcdMutex(ctx context.Context, mutexKey string) error {
	lease, err := conn.GetLeaseByKey(mutexKey)
	if err != nil {
		return err
	}

	leaseResp, err := conn.GetLeaseRespByKey(mutexKey)
	if err != nil {
		return err
	}

	_, err = lease.Revoke(ctx, leaseResp.ID)

	return err
}
