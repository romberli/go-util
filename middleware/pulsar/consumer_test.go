package pulsar

import (
	"context"
	"testing"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
)

func TestConsumer_Receive(t *testing.T) {
	producer := newTestProducer(t)
	defer producer.Disconnect()

	consumer := newTestConsumer(t, SubscriptionTypeExclusive)
	defer consumer.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	testPayload := "test receive message"
	msg := &pulsar.ProducerMessage{
		Payload: []byte(testPayload),
	}

	_, err := producer.Send(ctx, msg)
	if err != nil {
		t.Fatalf("failed to send message: %v", err)
	}

	receivedMsg, err := consumer.Receive(ctx)
	if err != nil {
		t.Fatalf("failed to receive message: %v", err)
	}

	if string(receivedMsg.Payload()) != testPayload {
		t.Fatalf("received payload mismatch: expected %s, got %s", testPayload, string(receivedMsg.Payload()))
	}

	err = consumer.Ack(receivedMsg)
	if err != nil {
		t.Fatalf("failed to ack message: %v", err)
	}

	t.Logf("message received successfully, payload: %s", string(receivedMsg.Payload()))
}

func TestConsumer_ReceiveAllSubscriptionTypes(t *testing.T) {
	subscriptionTypes := []SubscriptionType{
		SubscriptionTypeExclusive,
		SubscriptionTypeShared,
		SubscriptionTypeFailover,
		SubscriptionTypeKeyShared,
	}

	for _, subType := range subscriptionTypes {
		t.Run(subType.String(), func(t *testing.T) {
			subscriptionName := testSubscription + "-" + subType.String()

			producer := newTestProducer(t)
			defer producer.Disconnect()

			config := newTestConfig()
			consumer, err := NewConsumerWithDefault(config, subscriptionName, subType)
			if err != nil {
				t.Fatalf("failed to create consumer with type %v: %v", subType, err)
			}
			defer consumer.Close()

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			testPayload := "test message for " + subType.String()
			msg := &pulsar.ProducerMessage{
				Payload: []byte(testPayload),
			}

			_, err = producer.Send(ctx, msg)
			if err != nil {
				t.Fatalf("failed to send message: %v", err)
			}

			receivedMsg, err := consumer.Receive(ctx)
			if err != nil {
				t.Fatalf("failed to receive message: %v", err)
			}

			if string(receivedMsg.Payload()) != testPayload {
				t.Fatalf("received payload mismatch: expected %s, got %s", testPayload, string(receivedMsg.Payload()))
			}

			err = consumer.Ack(receivedMsg)
			if err != nil {
				t.Fatalf("failed to ack message: %v", err)
			}

			t.Logf("subscription type %v: message received successfully", subType)
		})
	}
}

func TestConsumer_Ack(t *testing.T) {
	producer := newTestProducer(t)
	defer producer.Disconnect()

	consumer := newTestConsumer(t, SubscriptionTypeExclusive)
	defer consumer.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	msg := &pulsar.ProducerMessage{
		Payload: []byte("test ack message"),
	}

	_, err := producer.Send(ctx, msg)
	if err != nil {
		t.Fatalf("failed to send message: %v", err)
	}

	receivedMsg, err := consumer.Receive(ctx)
	if err != nil {
		t.Fatalf("failed to receive message: %v", err)
	}

	err = consumer.Ack(receivedMsg)
	if err != nil {
		t.Fatalf("failed to ack message: %v", err)
	}

	t.Logf("message acked successfully")
}

func TestConsumer_Nack(t *testing.T) {
	producer := newTestProducer(t)
	defer producer.Disconnect()

	consumer := newTestConsumer(t, SubscriptionTypeShared)
	defer consumer.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	msg := &pulsar.ProducerMessage{
		Payload: []byte("test nack message"),
	}

	_, err := producer.Send(ctx, msg)
	if err != nil {
		t.Fatalf("failed to send message: %v", err)
	}

	receivedMsg, err := consumer.Receive(ctx)
	if err != nil {
		t.Fatalf("failed to receive message: %v", err)
	}

	consumer.Nack(receivedMsg)

	t.Logf("message nacked successfully")
}

func TestConsumer_Seek(t *testing.T) {
	producer := newTestProducer(t)
	defer producer.Disconnect()

	subscriptionName := testSubscription + "-seek"
	config := newTestConfig()
	consumer, err := NewConsumerWithDefault(config, subscriptionName, SubscriptionTypeExclusive)
	if err != nil {
		t.Fatalf("failed to create consumer: %v", err)
	}
	defer consumer.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	msgIDs := []pulsar.MessageID{}
	for i := 0; i < 3; i++ {
		msg := &pulsar.ProducerMessage{
			Payload: []byte("seek test message " + string(rune('A'+i))),
		}
		msgID, err := producer.Send(ctx, msg)
		if err != nil {
			t.Fatalf("failed to send message %d: %v", i, err)
		}
		msgIDs = append(msgIDs, msgID)
	}

	for i := 0; i < 3; i++ {
		receivedMsg, err := consumer.Receive(ctx)
		if err != nil {
			t.Fatalf("failed to receive message %d: %v", i, err)
		}
		err = consumer.Ack(receivedMsg)
		if err != nil {
			t.Fatalf("failed to ack message %d: %v", i, err)
		}
	}

	err = consumer.Seek(msgIDs[0])
	if err != nil {
		t.Fatalf("failed to seek to first message: %v", err)
	}

	receivedMsg, err := consumer.Receive(ctx)
	if err != nil {
		t.Fatalf("failed to receive message after seek: %v", err)
	}

	expectedPayload := "seek test message A"
	if string(receivedMsg.Payload()) != expectedPayload {
		t.Fatalf("after seek, payload mismatch: expected %s, got %s", expectedPayload, string(receivedMsg.Payload()))
	}

	err = consumer.Ack(receivedMsg)
	if err != nil {
		t.Fatalf("failed to ack message after seek: %v", err)
	}

	t.Logf("seek test passed, received first message again after seek")
}

func TestConsumer_SeekByTime(t *testing.T) {
	producer := newTestProducer(t)
	defer producer.Disconnect()

	subscriptionName := testSubscription + "-seekbytime"
	config := newTestConfig()
	consumer, err := NewConsumerWithDefault(config, subscriptionName, SubscriptionTypeExclusive)
	if err != nil {
		t.Fatalf("failed to create consumer: %v", err)
	}
	defer consumer.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for i := 0; i < 3; i++ {
		msg := &pulsar.ProducerMessage{
			Payload: []byte("seekbytime test message " + string(rune('A'+i))),
		}
		_, err := producer.Send(ctx, msg)
		if err != nil {
			t.Fatalf("failed to send message %d: %v", i, err)
		}
	}

	for i := 0; i < 3; i++ {
		receivedMsg, err := consumer.Receive(ctx)
		if err != nil {
			t.Fatalf("failed to receive message %d: %v", i, err)
		}
		err = consumer.Ack(receivedMsg)
		if err != nil {
			t.Fatalf("failed to ack message %d: %v", i, err)
		}
	}

	seekTime := time.Now().Add(-5 * time.Second)
	err = consumer.SeekByTime(seekTime)
	if err != nil {
		t.Fatalf("failed to seek by time: %v", err)
	}

	receivedMsg, err := consumer.Receive(ctx)
	if err != nil {
		t.Fatalf("failed to receive message after seek by time: %v", err)
	}

	payload := string(receivedMsg.Payload())
	if payload == "" {
		t.Fatal("payload should not be empty after seek by time")
	}

	err = consumer.Ack(receivedMsg)
	if err != nil {
		t.Fatalf("failed to ack message after seek by time: %v", err)
	}

	t.Logf("seek by time test passed, received message: %s", payload)
}
