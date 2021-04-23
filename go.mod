module github.com/romberli/go-util

go 1.15

replace github.com/ClickHouse/clickhouse-go v1.4.3 => github.com/romberli/clickhouse-go v1.4.4-0.20210422094559-b05fc8c4dbe9

require (
	github.com/ClickHouse/clickhouse-go v1.4.3
	github.com/Shopify/sarama v1.26.1
	github.com/hashicorp/go-multierror v1.1.0
	github.com/json-iterator/go v1.1.10
	github.com/pingcap/errors v0.11.4
	github.com/pkg/errors v0.9.1
	github.com/pkg/sftp v1.12.0
	github.com/prometheus/client_golang v1.10.0
	github.com/prometheus/common v0.18.0
	github.com/romberli/dynamic-struct v1.2.1
	github.com/romberli/log v1.0.20
	github.com/shirou/gopsutil/v3 v3.20.11
	github.com/siddontang/go v0.0.0-20180604090527-bdc77568d726
	github.com/siddontang/go-mysql v1.1.0
	github.com/stretchr/testify v1.6.1
	go.etcd.io/etcd v0.0.0-20200520232829-54ba9589114f
	go.uber.org/zap v1.16.0
	golang.org/x/crypto v0.0.0-20200820211705-5c72a883971a
)
