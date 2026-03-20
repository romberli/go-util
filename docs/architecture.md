# go-util 架构设计文档

## 1. 项目概述

`go-util` 是一个 Go 语言基础设施工具库，为上层业务项目提供标准化的数据库连接、加密、鉴权、系统调用、配置管理等能力，旨在减少重复代码、统一技术规范。

- **模块路径**：`github.com/romberli/go-util`
- **Go 版本**：1.24.0（toolchain go1.24.2）
- **开源协议**：见 LICENSE

---

## 2. 技术栈

| 类别     | 依赖                                                                |
|--------|-------------------------------------------------------------------|
| 错误处理   | `github.com/pingcap/errors`（带堆栈追踪）                                |
| 日志     | `github.com/romberli/log`（封装 logrus/zap）                          |
| 配置     | `github.com/spf13/viper`                                          |
| 测试     | `github.com/stretchr/testify` + `github.com/agiledragon/gomonkey` |
| 数据库    | MySQL（go-mysql/Percona）、ClickHouse v2                             |
| 消息队列   | Kafka（IBM Sarama）、RabbitMQ（amqp091-go）                            |
| SQL 解析 | TiDB Parser（`github.com/pingcap/tidb/pkg/parser`）                 |
| 加密     | `golang.org/x/crypto`、SM2/SM3（gmsm tjfoc）、RSA/AES 标准库             |
| 鉴权     | `github.com/golang-jwt/jwt/v5`                                    |
| 分布式    | etcd v3、Prometheus client                                         |
| 唯一 ID  | Snowflake（`github.com/bwmarrin/snowflake`）                        |
| 泛型工具   | Go 1.18+ 泛型类型约束                                                   |

---

## 3. 目录结构

```
go-util/
├── auth/           # JWT 鉴权模块
├── common/         # 通用工具函数
├── config/         # 配置管理
├── constant/       # 全局常量定义
├── crypto/         # 加密/解密模块
├── docs/           # 设计文档
├── http/           # HTTP 客户端工具
├── linux/          # Linux 系统级操作
├── middleware/     # 中间件抽象层
│   ├── clickhouse/ #   ClickHouse 客户端
│   ├── etcd/       #   etcd 客户端
│   ├── kafka/      #   Kafka 生产者/消费者
│   ├── mysql/      #   MySQL 客户端
│   ├── prometheus/ #   Prometheus 指标
│   ├── rabbitmq/   #   RabbitMQ 客户端
│   └── sql/        #   SQL 解析工具
├── types/          # 泛型类型约束
├── uid/            # 唯一 ID 生成
└── viper/          # Viper 配置加载封装
```

---

## 4. 模块详解

### 4.1 auth — JWT 鉴权

提供 JWT Token 的签发与验证，支持多种签名算法：

- **RSA**：非对称加密，适用于跨服务信任场景
- **ECDSA**：椭圆曲线签名，密钥更短、性能更好
- **HMAC**：对称签名，适用于内部服务

### 4.2 common — 通用工具

体量最大的模块（25+ 文件），涵盖：

| 子功能     | 说明                                   |
|---------|--------------------------------------|
| 类型转换    | 字符串、数字、布尔、时间等互转，处理 nil 边界            |
| JSON 处理 | 序列化/反序列化，字段掩码（`MaskJSON`）            |
| 敏感数据掩码  | 基于正则的字段替换（`mask.go`），支持 SQL 语句中的密码字段 |
| 重试逻辑    | 可配置次数和间隔的重试工具函数                      |
| 排序工具    | 通用排序辅助                               |
| 时间工具    | 时区处理、格式化                             |

**敏感数据掩码示例**：

```go
// 自动替换 SQL 中 IDENTIFIED BY 后的密码
masked := common.MaskString("CREATE USER 'u'@'%' IDENTIFIED BY 'secret'")
// → "CREATE USER 'u'@'%' IDENTIFIED BY '***'"

// 屏蔽 JSON 中指定字段
masked := common.MaskJSON(jsonStr, []string{"password", "token"})
```

### 4.3 crypto — 加密模块

支持三套加密体系：

| 算法  | 用途                     |
|-----|------------------------|
| RSA | 非对称加密/解密，密钥对管理         |
| AES | 对称加密，高性能数据加密           |
| SM2 | 中国国密标准（GM/T 0003），合规场景 |

### 4.4 middleware — 中间件抽象层

核心设计，通过统一接口屏蔽底层差异。

#### 接口定义（`middleware/middleware.go`）

```
Pool        连接池接口（Get/Put/Close）
Transaction 事务接口（Begin/Commit/Rollback）
Statement   语句接口（Execute/Query）
Result      结果集接口（行迭代、列读取）
```

#### 各实现模块

| 模块            | 特性                                                                      |
|---------------|-------------------------------------------------------------------------|
| `mysql/`      | 连接池、主从角色检测、版本信息获取                                                       |
| `clickhouse/` | ClickHouse 批量写入/查询                                                      |
| `kafka/`      | 生产者/消费者，支持分区策略                                                          |
| `rabbitmq/`   | Exchange/Queue 声明，消息确认                                                  |
| `etcd/`       | 分布式 KV 读写，Watch 监听                                                      |
| `prometheus/` | 指标注册与暴露                                                                 |
| `sql/`        | SQL 语句解析（基于 TiDB Parser），提取表名/列名/索引，识别语句类型（SELECT/INSERT/CREATE USER 等） |

### 4.5 linux — 系统级操作

| 子模块      | 功能          |
|----------|-------------|
| `ssh.go` | 远程 SSH 命令执行 |
| 进程管理     | 进程启停、PID 操作 |
| 网络工具     | 网卡、IP、端口检测  |
| 文件操作     | 文件读写、权限管理   |
| 挂载管理     | 磁盘挂载/卸载     |

### 4.6 其他模块

| 模块          | 说明                                        |
|-------------|-------------------------------------------|
| `config/`   | 应用配置结构体与加载逻辑                              |
| `constant/` | 全局字符串和数字常量，避免魔法字符串                        |
| `http/`     | HTTP 客户端封装（超时、重试、请求构建）                    |
| `types/`    | 泛型类型约束：`Primitive`、`Number`、`Int`、`Float` |
| `uid/`      | 基于 Snowflake 算法的分布式唯一 ID 生成               |
| `viper/`    | Viper 配置加载封装，支持热重载（fsnotify）              |

---

## 5. 核心设计原则

### 5.1 接口优先

`middleware/` 层严格面向接口编程，调用方依赖抽象而非具体实现，便于替换底层组件和单元测试。

### 5.2 敏感数据防泄漏

在日志输出、JSON 序列化、SQL 打印等环节，通过 `common/mask.go` 统一处理敏感字段（密码、Token、密钥），防止信息泄漏。

### 5.3 常量集中管理

所有魔法字符串和数字均在 `constant/` 包中统一定义，便于维护和全局替换。

### 5.4 错误可追溯

使用 `github.com/pingcap/errors` 在每个错误传播点调用 `errors.Trace()`，保留完整调用栈，方便定位问题。

### 5.5 测试分层

- `xxx_test.go`：纯单元测试，无外部依赖，CI 中直接运行
- `xxx_local_test.go`：需要本地数据库/中间件的集成测试，本地开发时手动运行

---

## 6. 数据流示意

```
业务项目
    │
    ├── auth/         签发/验证 JWT
    ├── crypto/       加密敏感配置
    ├── config/       加载应用配置
    │   └── viper/
    │
    ├── middleware/   统一中间件接口
    │   ├── mysql     → MySQL 客户端
    │   ├── clickhouse → ClickHouse 客户端
    │   ├── kafka     → 消息发布/订阅
    │   ├── rabbitmq  → RabbitMQ 客户端
    │   ├── etcd      → 分布式配置
    │   └── prometheus→ 指标上报
    │
    ├── common/       工具函数（转换、掩码、重试）
    ├── linux/        系统操作（SSH、进程、网络）
    └── uid/          生成唯一 ID
```

---

## 7. 扩展指南

### 新增中间件

1. 在 `middleware/` 下创建子包
2. 实现 `middleware.go` 中定义的接口（`Pool`、`Statement` 等）
3. 编写 `xxx_test.go` 单元测试 和 `xxx_local_test.go` 集成测试

### 新增工具函数

- 通用函数放入 `common/`
- 新增常量在 `constant/` 中定义
- 若功能独立且较大，考虑新建顶层包

---

*文档生成时间：2026-03-19*
