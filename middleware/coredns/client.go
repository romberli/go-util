package coredns

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"go.etcd.io/etcd/clientv3"

	"github.com/romber2001/go-util/middleware/etcd"
)

const (
	DefaultMutexPrefix = "/mutex/coredns"
	DefaultTimeFormat  = "2020-08-20 11:55:00.000000"
	InfiniteTTLTime    = 0
	MinimumTTLTime     = 1
	MaximumTTLTime     = 10
)

var ErrInvalidTTL = errors.New(
	fmt.Sprintf("TTL must be larger then %d and smaller than %d.", MinimumTTLTime, MaximumTTLTime))
var ErrLockFailed = errors.New(fmt.Sprintf("lock mutex failed, please try again later."))

// ARecordValue saves the unmarshalled json data which presents A record value
type ARecordValue struct {
	Host string `json:"host"`
	TTL  int64  `json:"ttl"`
}

// NewARecordValue return *ARecordValue, it uses core dns A record value which is a json string as the input
func NewARecordValue(value []byte) (*ARecordValue, error) {
	a := &ARecordValue{}
	err := json.Unmarshal(value, a)
	if err != nil {
		return nil, err
	}

	return a, nil

}

type Conn struct {
	Endpoints []string
	Path      string
	EtcdConn  *etcd.Conn
}

// NewCoreDNSConn returns a coredns connection, coredns must use etcd as the backend
func NewCoreDNSConn(endpoints []string, path string) (*Conn, error) {
	conn, err := etcd.NewEtcdConn(endpoints)
	if err != nil {
		return nil, err
	}

	return &Conn{
		Endpoints: endpoints,
		Path:      path,
		EtcdConn:  conn,
	}, nil
}

// Close close the connection with etcd
func (conn *Conn) Close() error {
	return conn.EtcdConn.Close()
}

// GetEtcdKeyNameFromURL transfers url string to backward and add coredns path as the prefix
func (conn *Conn) GetEtcdKeyNameFromURL(url string) (string, error) {
	result := ""

	strList := strings.Split(url, ".")
	for _, str := range strList {
		str = strings.TrimSpace(str)
		if str == "." || str == "" {
			return "", errors.New(
				fmt.Sprintf("some part of url is empty or contains invaid characters [%s]", str))
		}

		result = str + "/" + result
	}

	result = conn.Path + "/" + result
	result = strings.TrimSuffix(result, "/")

	return result, nil
}

// GetEtcdValue checks if ttl if valid and returns a json string with input ip and ttl
func (conn *Conn) GetEtcdValue(ip string, ttl int64) (string, error) {
	if ttl == InfiniteTTLTime {
		return fmt.Sprintf(`{"host": "%s"}`, ip), nil
	}

	if ttl > MaximumTTLTime || ttl < MinimumTTLTime {
		return "", ErrInvalidTTL
	}

	return fmt.Sprintf(`{"host": "%s", "ttl": %d}`, ip, ttl), nil
}

// Resolve returns host ip slice which is resolved as the given url
func (conn *Conn) Resolve(ctx context.Context, url string) ([]string, error) {
	var result []string

	key, err := conn.GetEtcdKeyNameFromURL(url)
	if err != nil {
		return nil, err
	}

	getResp, err := conn.EtcdConn.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	for _, kv := range getResp.Kvs {
		value, err := NewARecordValue(kv.Value)
		if err != nil {
			return nil, err
		}

		result = append(result, value.Host)
	}

	return result, nil
}

// ResolveWithTTL works like Resolve, but it returns both host ip and ttl,
// the host ip is the key and ttl is the value in the map
func (conn *Conn) ResolveWithTTL(ctx context.Context, url string) ([]map[string]int64, error) {
	var result []map[string]int64

	key, err := conn.GetEtcdKeyNameFromURL(url)
	if err != nil {
		return nil, err
	}

	getResp, err := conn.EtcdConn.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	for _, kv := range getResp.Kvs {
		value, err := NewARecordValue(kv.Value)
		if err != nil {
			return nil, err
		}

		v := make(map[string]int64)
		v[value.Host] = value.TTL

		result = append(result, v)
	}

	return result, nil
}

// PutARecord used to add or modify the A record of coredns,
// it will lock mutex before really putting key into etcd,
// if lock mutex failed, it will return an error
func (conn *Conn) PutARecord(ctx context.Context, url, ip string, ttl int64) error {
	key, err := conn.GetEtcdKeyNameFromURL(url)
	if err != nil {
		return err
	}

	value, err := conn.GetEtcdValue(ip, ttl)
	if err != nil {
		return err
	}

	mutexKey := DefaultMutexPrefix + key
	mutexValue := time.Now().Format(DefaultTimeFormat)

	ok, err := conn.EtcdConn.LockEtcdMutex(ctx, mutexKey, mutexValue, etcd.DefaultMutexLeaseSeconds)
	if err != nil {
		return err
	}

	if ok {
		defer func() {
			_, err = conn.EtcdConn.UnlockEtcdMutex(ctx, mutexKey)
		}()

		_, _, err = conn.EtcdConn.PutWithTTL(ctx, key, value, ttl)

		return err
	}

	return ErrLockFailed
}

// PutARecordAndKeepAlive works like PutARecord, used to add or modify the A record of coredns,
// it will lock mutex before really putting key into etcd,
// if lock mutex failed, it will return an error,
// it also keep alive the key in etcd,
func (conn *Conn) PutARecordAndKeepAlive(ctx context.Context, url, ip string, ttl int64) error {
	err := conn.PutARecord(ctx, url, ip, ttl)
	if err != nil {
		return err
	}

	key, err := conn.GetEtcdKeyNameFromURL(url)
	if err != nil {
		return err
	}

	leaseID, err := conn.EtcdConn.GetLeaseIDByKey(key)
	if err != nil {
		return err
	}

	_, err = conn.EtcdConn.KeepAlive(ctx, leaseID)

	return err
}

// DeleteARecord delete the A record by deleting the concerned key in the etcd,
// it will lock mutex before really doing this,
// if lock failed, it will return an error
func (conn *Conn) DeleteARecord(ctx context.Context, url string) error {
	key, err := conn.GetEtcdKeyNameFromURL(url)
	if err != nil {
		return err
	}

	mutexKey := DefaultMutexPrefix + key
	mutexValue := time.Now().Format(DefaultTimeFormat)

	ok, err := conn.EtcdConn.LockEtcdMutex(ctx, mutexKey, mutexValue, etcd.DefaultMutexLeaseSeconds)
	if err != nil {
		return err
	}

	if ok {
		defer func() {
			_, err = conn.EtcdConn.UnlockEtcdMutex(ctx, mutexKey)
		}()

		_, err = conn.EtcdConn.Delete(ctx, key)

		return err
	}

	return ErrLockFailed
}
