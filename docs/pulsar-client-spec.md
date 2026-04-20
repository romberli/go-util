# Pulsar 客户端规格设计

## 1. 背景与目标

新增 `middleware/pulsar` 包，为上层业务提供 Apache Pulsar 的生产者与消费者能力。设计参考 `middleware/kafka` 的扁平风格：无 Pool、无子包、每个角色一个文件，利用 `pulsar.Client` 原生的连接管理能力替代手动池化。

**依赖库**：`github.com/apache/pulsar-client-go v0.13.x`

---

## 2. 参考模式：Kafka vs RabbitMQ

| 维度 | Kafka（参考） | RabbitMQ（不采用） |
|------|--------------|-------------------|
| 目录结构 | 全部平铺在 `kafka/` | 多级子包（client/producer/consumer） |
| 连接管理 | Client 内建，无 Pool | 手动 Pool + 维护 goroutine |
| 文件数 | 3 个（producer/consumer/util） | 7 个文件 |
| 角色结构 | `AsyncProducer{ Client, Producer }` | `PoolProducer{ Conn, Pool }` |

---

## 3. Pulsar 与 RabbitMQ 概念对照

| RabbitMQ 概念 | Pulsar 概念 | 说明 |
|--------------|-------------|------|
| Connection | Client | 顶层连接对象，线程安全可复用 |
| Channel | Producer / Consumer | 分别创建，各自独立 |
| Exchange + Queue + Binding | Topic | Pulsar Topic 是第一等公民 |
| Consumer Tag | Subscription Name | 标识一个消费订阅 |
| Ack / Nack | Ack / Nack / ReconsumeLater | Pulsar 还支持延迟重投 |
| Exchange Type | Subscription Type | Exclusive / Shared / Failover / KeyShared |
| Vhost | Tenant / Namespace | Topic 全名：`persistent://tenant/namespace/topic` |

---

## 4. 目录结构

```
middleware/pulsar/
├── config.go       # Config：连接 + Topic + 认证 + 订阅参数
├── producer.go     # Producer：封装 pulsar.Client + pulsar.Producer
└── consumer.go     # Consumer：封装 pulsar.Client + pulsar.Consumer
                    #           + DefaultConsumerHandler
```

- 无 Pool、无子包，与 `middleware/kafka` 风格完全对齐
- `pulsar.Message` 已足够完备，不额外定义消息结构体

---

## 5. 各层规格

### 5.1 配置层（`config.go`）

```go
type Config struct {
    URL              string // Pulsar broker URL，如 "pulsar://localhost:6650"
    AuthToken        string // JWT Token，空则不启用认证
    Topic            string // 完整 topic 名，如 "persistent://tenant/ns/topic"
    SubscriptionName string // 订阅名称（Consumer 专属）
    SubscriptionType string // Exclusive | Shared | Failover | KeyShared，默认 Shared
}
```

构造函数：
- `NewConfig(url, authToken, topic, subscriptionName, subscriptionType string) *Config`
- `NewConfigWithDefault(url, topic string) *Config` — 无认证，SubscriptionType 默认 Shared

辅助方法：
- `Clone() *Config`
- `GetSubscriptionType() pulsar.SubscriptionType` — 字符串映射为 SDK 枚举

**SubscriptionType 映射**：
```
"Exclusive"  → pulsar.Exclusive
"Shared"     → pulsar.Shared      （默认）
"Failover"   → pulsar.Failover
"KeyShared"  → pulsar.KeyShared
```

---

### 5.2 生产者（`producer.go`）

```go
type Producer struct {
    Config   *Config
    Client   pulsar.Client
    Producer pulsar.Producer
}
```

构造函数：
- `NewProducer(url, topic string) (*Producer, error)`
- `NewProducerWithConfig(config *Config) (*Producer, error)` — 内部依次调用 `pulsar.NewClient()` → `client.CreateProducer()`

方法：

| 方法签名 | 说明 |
|---------|------|
| `Close() error` | 关闭 Producer 和 Client |
| `Send(ctx context.Context, payload []byte) (pulsar.MessageID, error)` | 同步发送原始字节 |
| `SendJSON(ctx context.Context, v interface{}) (pulsar.MessageID, error)` | JSON 序列化后同步发送 |
| `SendWithProperties(ctx context.Context, payload []byte, properties map[string]string) (pulsar.MessageID, error)` | 携带自定义属性同步发送 |
| `SendAsync(ctx context.Context, payload []byte, callback func(pulsar.MessageID, *pulsar.ProducerMessage, error))` | 异步发送，结果通过 callback 回调 |

---

### 5.3 消费者（`consumer.go`）

```go
type Consumer struct {
    Config   *Config
    Client   pulsar.Client
    Consumer pulsar.Consumer
}
```

构造函数：
- `NewConsumer(url, topic, subscriptionName, subscriptionType string) (*Consumer, error)`
- `NewConsumerWithConfig(config *Config) (*Consumer, error)` — 内部依次调用 `pulsar.NewClient()` → `client.Subscribe()`

方法：

| 方法签名 | 说明 |
|---------|------|
| `Close() error` | 关闭 Consumer 和 Client |
| `Receive(ctx context.Context) (pulsar.Message, error)` | 阻塞接收一条消息 |
| `Chan() <-chan pulsar.Message` | 返回消息 channel，供 range 消费 |
| `Ack(msg pulsar.Message) error` | 确认消费成功 |
| `AckID(msgID pulsar.MessageID) error` | 按消息 ID 确认 |
| `Nack(msg pulsar.Message)` | 否定确认，触发重投 |
| `ReconsumeLater(msg pulsar.Message, delay time.Duration)` | 延迟重投 |
| `Seek(msgID pulsar.MessageID) error` | 按 MessageID 重置消费位点 |
| `SeekByTime(t time.Time) error` | 按时间重置消费位点 |
| `Unsubscribe() error` | 取消订阅 |

**DefaultConsumerHandler**（类比 Kafka 的 `DefaultConsumerGroupHandler`）：

```go
type ConsumerHandler interface {
    Handle(consumer *Consumer, msg pulsar.Message) error
}

type DefaultConsumerHandler struct{}

// Handle 打印消息日志并自动 Ack，业务方可实现 ConsumerHandler 接口覆盖此行为
func (h DefaultConsumerHandler) Handle(consumer *Consumer, msg pulsar.Message) error
```

`Consumer` 提供驱动方法：
```go
// Consume 循环消费，使用给定 handler 处理每条消息，ctx 取消时退出
func (c *Consumer) Consume(ctx context.Context, handler ConsumerHandler) error
```

---

## 6. 关键设计决策

### 6.1 无 Pool
`pulsar.Client` 内建连接管理（连接复用、重连、负载均衡），无需应用层手动池化。每个 `Producer` / `Consumer` 实例各持一个 `pulsar.Client`，生命周期与角色绑定，`Close()` 时一并关闭。

### 6.2 同步 + 异步发送并存
- `Send` / `SendJSON` / `SendWithProperties` — 同步，阻塞等待 broker 确认，返回 `MessageID`
- `SendAsync` — 异步，立即返回，结果通过 callback 通知，适合高吞吐场景

### 6.3 ConsumerHandler 模式
对齐 Kafka 的 `ConsumerGroupHandler` 接口模式：
- 定义 `ConsumerHandler` 接口，业务方实现自定义逻辑
- 提供 `DefaultConsumerHandler`（日志 + 自动 Ack）开箱即用
- `Consumer.Consume(ctx, handler)` 驱动消费循环，ctx 取消时优雅退出

### 6.4 错误处理
遵循项目规范，所有 error 通过 `errors.Trace(err)` 包装后返回，保留完整调用栈。

---

## 7. 测试规划

| 文件 | 说明 |
|------|------|
| `producer_test.go` | Producer 单元测试 |
| `producer_local_test.go` | Producer 集成测试（需本地 Pulsar 实例） |
| `consumer_test.go` | Consumer 单元测试 |
| `consumer_local_test.go` | Consumer 集成测试（需本地 Pulsar 实例） |

本地测试默认连接地址：`pulsar://localhost:6650`

---

## 8. 与 Kafka 实现的对照

| 维度 | Kafka | Pulsar |
|------|-------|--------|
| 生产者结构 | `AsyncProducer{ Client, Producer }` | `Producer{ Config, Client, Producer }` |
| 消费者结构 | `ConsumerGroup{ Client, Group }` | `Consumer{ Config, Client, Consumer }` |
| 消费处理器 | `ConsumerGroupHandler` 接口 + `DefaultConsumerGroupHandler` | `ConsumerHandler` 接口 + `DefaultConsumerHandler` |
| 消费驱动 | `ConsumerGroup.Consume(ctx, topic, handler)` | `Consumer.Consume(ctx, handler)` |
| 异步发送 | `Producer.Input() <- msg`（channel 方式） | `Producer.SendAsync(ctx, payload, callback)` |
| 消息确认 | `sess.MarkMessage(msg, "")` | `consumer.Ack(msg)` |

---

*文档版本：v1.0 | 创建时间：2026-03-19*
