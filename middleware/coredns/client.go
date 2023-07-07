package coredns

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/pingcap/errors"
	"github.com/romberli/log"

	"github.com/romberli/go-util/constant"

	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/romberli/go-util/middleware/etcd"
)

const (
	DefaultMutexPrefix = "/mutex/coredns"
	DefaultTimeFormat  = "2020-08-20 11:55:00.000000"
	InfiniteTTLTime    = 0
	MinimumTTLTime     = 1
	MaximumTTLTime     = 10
)

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
		return nil, errors.Trace(err)
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

// Close closes the connection with etcd
func (conn *Conn) Close() error {
	return conn.EtcdConn.Close()
}

// GetEtcdKeyNameFromURL transfers url string to backward and add coredns path as the prefix
func (conn *Conn) GetEtcdKeyNameFromURL(url string) (string, error) {
	var result string

	strList := strings.Split(url, constant.DotString)
	for _, str := range strList {
		str = strings.TrimSpace(str)
		if str == constant.DotString || str == constant.EmptyString {
			return constant.EmptyString, errors.Errorf("some part of url is empty or contains invalid characters [%s]", str)
		}

		result = str + constant.SlashString + result
	}

	result = conn.Path + constant.SlashString + result
	result = strings.TrimSuffix(result, constant.SlashString)

	return result, nil
}

// GetEtcdValue checks if ttl if valid and returns a json string with input ip and ttl
func (conn *Conn) GetEtcdValue(ip string, ttl int64) (string, error) {
	if ttl == InfiniteTTLTime {
		return fmt.Sprintf(`{"host": "%s"}`, ip), nil
	}

	if ttl > MaximumTTLTime || ttl < MinimumTTLTime {
		return constant.EmptyString, errors.Errorf("TTL must be larger then %d and smaller than %d.", MinimumTTLTime, MaximumTTLTime)
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
		return nil, errors.Trace(err)
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
		return nil, errors.Trace(err)
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
			if err != nil {
				log.Errorf("unlock etcd mutex failed. error:%n%+v", err)
			}
		}()

		_, _, err = conn.EtcdConn.PutWithTTL(ctx, key, value, ttl)

		return err
	}

	return errors.New("lock mutex failed, please try again later.")
}

// PutARecordAndKeepAlive works like PutARecord, used to add or modify the A record of coredns,
// it will lock mutex before really putting key into etcd,
// if lock mutex failed, it will return an error,
// it also keeps alive the key in etcd,
func (conn *Conn) PutARecordAndKeepAlive(ctx context.Context, url, ip string, ttl int64) error {
	err := conn.PutARecord(ctx, url, ip, ttl)
	if err != nil {
		return err
	}

	key, err := conn.GetEtcdKeyNameFromURL(url)
	if err != nil {
		return err
	}

	leaseID := conn.EtcdConn.GetLeaseIDByKey(key)
	if leaseID == clientv3.NoLease {
		return nil
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
			if err != nil {
				log.Errorf("unlock etcd mutex failed. error:%n%+v", err)
			}
		}()

		_, err = conn.EtcdConn.Delete(ctx, key)

		return err
	}

	return errors.New("lock mutex failed, please try again later.")
}
