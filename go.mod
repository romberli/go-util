module github.com/romberli/go-util

go 1.16

replace github.com/ClickHouse/clickhouse-go v1.4.7 => github.com/romberli/clickhouse-go v1.4.4-0.20210902113008-bb38dc6f756d

require (
	github.com/ClickHouse/clickhouse-go v1.4.7
	github.com/Shopify/sarama v1.26.1
	github.com/go-mysql-org/go-mysql v1.3.0
	github.com/hashicorp/go-multierror v1.1.0
	github.com/json-iterator/go v1.1.11
	github.com/percona/go-mysql v0.0.0-20210427141028-73d29c6da78c
	github.com/pingcap/parser v0.0.0-20211004012448-687005894c4e
	github.com/pkg/errors v0.9.1
	github.com/pkg/sftp v1.12.0
	github.com/prometheus/client_golang v1.11.0
	github.com/prometheus/common v0.26.0
	github.com/romberli/dynamic-struct v1.2.1
	github.com/romberli/log v1.0.20
	github.com/shirou/gopsutil/v3 v3.20.11
	github.com/siddontang/go v0.0.0-20180604090527-bdc77568d726
	github.com/stretchr/testify v1.7.0
	go.etcd.io/etcd/client/v3 v3.5.1
	go.uber.org/zap v1.19.1
	golang.org/x/crypto v0.0.0-20201002170205-7f63de1d35b0
)
