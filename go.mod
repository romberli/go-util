module github.com/romberli/go-util

go 1.16

replace github.com/ClickHouse/clickhouse-go v1.4.3 => github.com/romberli/clickhouse-go v1.4.4-0.20210422094559-b05fc8c4dbe9

require (
	github.com/ClickHouse/clickhouse-go v1.4.3
	github.com/Shopify/sarama v1.26.1
	github.com/go-mysql-org/go-mysql v1.3.0
	github.com/hashicorp/go-multierror v1.1.0
	github.com/json-iterator/go v1.1.10
	github.com/percona/go-mysql v0.0.0-20210427141028-73d29c6da78c
	github.com/pingcap/parser v0.0.0-20210525032559-c37778aff307
	github.com/pingcap/tidb v1.1.0-beta.0.20210526073135-acf5e52ffc78
	github.com/pkg/errors v0.9.1
	github.com/pkg/sftp v1.12.0
	github.com/prometheus/client_golang v1.10.0
	github.com/prometheus/common v0.18.0
	github.com/romberli/dynamic-struct v1.2.1
	github.com/romberli/log v1.0.20
	github.com/shirou/gopsutil/v3 v3.20.11
	github.com/siddontang/go v0.0.0-20180604090527-bdc77568d726
	github.com/stretchr/testify v1.6.1
	go.etcd.io/etcd v0.5.0-alpha.5.0.20200824191128-ae9734ed278b
	go.uber.org/zap v1.16.0
	golang.org/x/crypto v0.0.0-20200820211705-5c72a883971a
)
