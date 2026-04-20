# Pulsar Middleware Package Design

## Overview

Create a Pulsar message publishing/subscribing utility package under `middleware/pulsar/`, following a simplified structure similar to `pulsar_old` but with improved design.

## Requirements Summary

- **No connection pool**: Pulsar client (`apache/pulsar-client-go`) has built-in connection pooling via Client object
- **Single Client**: One Client + multiple Producer/Consumer pattern
- **No custom Message struct**: Use pulsar SDK's built-in `ProducerMessage` and `ConsumerMessage`
- **Support all subscription types**: Exclusive, Shared, Failover, KeyShared
- **Token authentication only**: Static token, no refresh mechanism
- **Unit tests required**: Follow rabbitmq package test style

## Directory Structure

```
middleware/pulsar/
├── config.go          # Config struct for connection and subscription
├── config_test.go     # Config unit tests
├── producer.go        # Producer wrapper
├── producer_test.go   # Producer unit tests
├── consumer.go        # Consumer wrapper with handler interface
├── consumer_test.go   # Consumer unit tests
```

## Component Design

### 1. Config (`config.go`)

**Struct:**
```go
type Config struct {
    URL              string  // Pulsar broker URL, e.g. "pulsar://localhost:6650"
    Token            string  // Authentication token (optional)
    Topic            string  // Topic name
    SubscriptionName string  // Subscription name (for Consumer)
    SubscriptionType string  // Subscription type: Exclusive/Shared/Failover/KeyShared
}
```

**Constants:**
- `DefaultURL = "pulsar://localhost:6650"`
- `DefaultSubscriptionType = "Shared"`
- `SubscriptionTypeExclusive = "Exclusive"`
- `SubscriptionTypeShared = "Shared"`
- `SubscriptionTypeFailover = "Failover"`
- `SubscriptionTypeKeyShared = "KeyShared"`

**Methods:**
- `NewConfig(url, token, topic, subscriptionName, subscriptionType) *Config`
- `NewConfigWithDefault(url, topic) *Config` - uses default subscription type, no token
- `Clone() *Config`
- `GetSubscriptionType() pulsar.SubscriptionType` - maps string to SDK enum

### 2. Producer (`producer.go`)

**Struct:**
```go
type Producer struct {
    Config   *Config
    Client   pulsar.Client
    Producer pulsar.Producer
}
```

**Methods:**
- `NewProducer(url, topic) (*Producer, error)` - simple constructor
- `NewProducerWithConfig(config) (*Producer, error)` - full constructor
- `Close() error` - closes producer then client
- `Send(ctx context.Context, payload []byte) (pulsar.MessageID, error)` - sync send
- `SendJSON(ctx context.Context, v interface{}) (pulsar.MessageID, error)` - marshal and send JSON
- `SendWithProperties(ctx context.Context, payload []byte, properties map[string]string) (pulsar.MessageID, error)` - send with custom properties
- `SendAsync(ctx context.Context, payload []byte, callback func(pulsar.MessageID, *pulsar.ProducerMessage, error))` - async send

### 3. Consumer (`consumer.go`)

**Struct:**
```go
type Consumer struct {
    Config   *Config
    Client   pulsar.Client
    Consumer pulsar.Consumer
}

type ConsumerHandler interface {
    Handle(consumer *Consumer, msg pulsar.Message) error
}

type DefaultConsumerHandler struct{} // logs message and auto-acks
```

**Methods:**
- `NewConsumer(url, topic, subscriptionName, subscriptionType) (*Consumer, error)`
- `NewConsumerWithConfig(config) (*Consumer, error)`
- `Close() error` - closes consumer then client
- `Receive(ctx context.Context) (pulsar.Message, error)` - blocking receive
- `Chan() <-chan pulsar.ConsumerMessage` - message channel for range-based consumption
- `Ack(msg pulsar.Message) error`
- `AckID(msgID pulsar.MessageID) error`
- `Nack(msg pulsar.Message)`
- `ReconsumeLater(msg pulsar.Message, delay time.Duration)`
- `Seek(msgID pulsar.MessageID) error`
- `SeekByTime(t time.Time) error`
- `Unsubscribe() error`
- `Consume(ctx context.Context, handler ConsumerHandler) error` - consumption loop

## Dependencies

- `github.com/apache/pulsar-client-go/pulsar` - already in go.mod
- `github.com/pingcap/errors` - error wrapping (consistent with other packages)
- `github.com/romberli/log` - logging (used in consumer handler)
- `github.com/romberli/go-util/constant` - constants

## Testing Strategy

- Each component has corresponding `_test.go` file
- Use `github.com/stretchr/testify` for assertions
- Tests require running Pulsar broker (can use Docker or skip in CI)
- Follow patterns from `rabbitmq` and `pulsar_old` test files