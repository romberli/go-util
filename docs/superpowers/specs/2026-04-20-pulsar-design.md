# Pulsar 工具包设计文档

## 概述

在 `middleware/pulsar` 目录下创建 Pulsar 消息队列工具包，用于消息发布和订阅。参考现有 `middleware/rabbitmq` 包的设计风格，但根据 Pulsar 的特性进行简化。

## 背景

- 项目已依赖 `apache/pulsar-client-go v0.13.0`
- Pulsar 内置连接管理，Producer/Consumer 通过共享 Client 连接，无需额外连接池
- 参考 rabbitmq 包结构但简化：无 pool.go、无 message.go、无 client 子目录

## 设计决策

| 决策项 | 选择 | 原因 |
|--------|------|------|
| 连接池 | 无池化，使用 Pulsar 内置 | Pulsar Client 本身管理连接，Producer/Consumer 轻量 |
| 订阅模式 | 全部四种 | Exclusive/Shared/Failover/Key_Shared |
| 消息结构 | 无特殊结构 | 直接使用 []byte 或 string |
| 认证方式 | Token 认证 | 仅支持 Token 字符串，无 TLS |
| Config 参数 | 核心参数 | URL、Topic、Subscription、SubscriptionType、Token |

## 目录结构

```
middleware/pulsar/
├── config.go          # Config 配置结构体
├── producer/
│   └── producer.go    # Producer 发布者
├── consumer/
│   └── consumer.go    # Consumer 消费者
```

## 组件设计

### config.go

```go
type Config struct {
    URL             string                 // Pulsar broker URL, e.g. "pulsar://localhost:6650"
    Topic           string                 // Topic name
    Subscription    string                 // Subscription name (consumer only)
    SubscriptionType pulsar.SubscriptionType // Exclusive/Shared/Failover/Key_Shared, 默认 Exclusive
    Token           string                 // 认证 Token，空字符串表示无认证
}

// 方法
NewConfig(url, topic, subscription, subscriptionType, token) *Config
NewConfigWithDefault(url, topic) *Config  // 无认证，Exclusive 模式，无订阅
Clone() *Config
```

### producer/producer.go

```go
type Producer struct {
    Client   pulsar.Client
    Producer pulsar.Producer
    Config   *Config
}

// 方法
NewProducer(config) (*Producer, error)
NewProducerWithDefault(url, topic) (*Producer, error)  // 无认证快速创建
Close() error
Send(ctx context.Context, message []byte) (pulsar.MessageID, error)  // 同步发送
SendAsync(ctx context.Context, message []byte, callback func(pulsar.MessageID, error))  // 异步发送
SendJSON(message string) (pulsar.MessageID, error)  // 发送 JSON 消息
```

### consumer/consumer.go

```go
type Consumer struct {
    Client   pulsar.Client
    Consumer pulsar.Consumer
    Config   *Config
}

// 方法
NewConsumer(config) (*Consumer, error)
NewConsumerWithDefault(url, topic, subscription) (*Consumer, error)  // 无认证 Exclusive
Close() error
Receive(ctx context.Context) (pulsar.Message, error)  // 接收单条消息
Chan() <-chan pulsar.Message  // 返回消息 channel
Ack(message pulsar.Message) error  // 確认消息
Nack(message pulsar.Message) error  // 否认消息
SeekByTime(time time.Time) error  // 按时间回溯
SeekByMessageID(messageID pulsar.MessageID) error  // 按 MessageID 回溯
```

## 错误处理

使用 `github.com/pingcap/errors` 包，遵循项目现有风格：
- 错误使用 `errors.Trace()` 包装
- 关闭操作忽略已关闭状态

## 使用示例

### Producer 示例

```go
config := pulsar.NewConfig("pulsar://localhost:6650", "my-topic", "", pulsar.Exclusive, "my-token")
producer, err := pulsarproducer.NewProducer(config)
if err != nil {
    log.Fatal(err)
}
defer producer.Close()

msgID, err := producer.Send(context.Background(), []byte("hello"))
```

### Consumer 示例

```go
config := pulsar.NewConfig("pulsar://localhost:6650", "my-topic", "my-sub", pulsar.Exclusive, "my-token")
consumer, err := pulsarconsumer.NewConsumer(config)
if err != nil {
    log.Fatal(err)
}
defer consumer.Close()

for msg := range consumer.Chan() {
    fmt.Println(string(msg.Payload()))
    consumer.Ack(msg)
}
```

## 依赖

- `github.com/apache/pulsar-client-go v0.13.0` - 已在 go.mod 中
- `github.com/pingcap/errors` - 已在 go.mod 中
- `github.com/romberli/go-util/constant` - 常量包