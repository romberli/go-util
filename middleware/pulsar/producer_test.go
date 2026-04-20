package pulsar

import (
	"context"
	"testing"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"

	"github.com/romberli/go-util/constant"
)

const (
	testPulsarURL    = "pulsar://192.168.137.11:6650"
	testToken        = ""
	testTopic        = "persistent://public/default/test-topic"
	testSubscription = "test-subscription"
)

func newTestConfig() *Config {
	return NewConfig(testPulsarURL, testToken, testTopic)
}

func newTestProducer(t *testing.T) *Producer {
	config := newTestConfig()
	producer, err := NewProducerWithDefault(config)
	if err != nil {
		t.Fatalf("failed to create producer: %v", err)
	}
	return producer
}

func newTestConsumer(t *testing.T, subscriptionType SubscriptionType) *Consumer {
	config := newTestConfig()
	consumer, err := NewConsumerWithDefault(config, testSubscription, subscriptionType)
	if err != nil {
		t.Fatalf("failed to create consumer: %v", err)
	}
	return consumer
}

func TestProducer_Send(t *testing.T) {
	producer := newTestProducer(t)
	defer producer.Disconnect()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	msg := &pulsar.ProducerMessage{
		Payload: []byte("test message"),
	}

	msgID, err := producer.Send(ctx, msg)
	if err != nil {
		t.Fatalf("failed to send message: %v", err)
	}

	if msgID == nil {
		t.Fatal("message ID should not be nil")
	}

	t.Logf("message sent successfully, ID: %v", msgID)
}

func TestProducer_SendAsync(t *testing.T) {
	producer := newTestProducer(t)
	defer producer.Disconnect()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	msg := &pulsar.ProducerMessage{
		Payload: []byte("test async message"),
	}

	done := make(chan bool, 1)
	producer.SendAsync(ctx, msg, func(msgID pulsar.MessageID, producerMsg *pulsar.ProducerMessage, err error) {
		if err != nil {
			t.Fatalf("async send failed: %v", err)
		}
		if msgID == nil {
			t.Fatal("async message ID should not be nil")
		}
		t.Logf("async message sent successfully, ID: %v", msgID)
		done <- true
	})

	select {
	case <-done:
	case <-ctx.Done():
		t.Fatal("async send timeout")
	}
}

func TestProducer_SendJSON(t *testing.T) {
	producer := newTestProducer(t)
	defer producer.Disconnect()

	consumer := newTestConsumer(t, SubscriptionTypeExclusive)
	defer consumer.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	testData := map[string]string{
		"key": "value",
		"foo": "bar",
	}

	msgID, err := producer.SendJSON(ctx, testData)
	if err != nil {
		t.Fatalf("failed to send JSON message: %v", err)
	}

	if msgID == nil {
		t.Fatal("JSON message ID should not be nil")
	}

	msg, err := consumer.Receive(ctx)
	if err != nil {
		t.Fatalf("failed to receive message: %v", err)
	}

	payload := string(msg.Payload())
	if payload == constant.EmptyString {
		t.Fatal("payload should not be empty")
	}

	t.Logf("JSON message sent and received successfully, ID: %v, payload: %s", msgID, payload)

	if err := consumer.Ack(msg); err != nil {
		t.Fatalf("failed to ack message: %v", err)
	}
}
