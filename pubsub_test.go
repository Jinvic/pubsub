package pubsub

import (
	"testing"
	"time"
)

func TestPublish(t *testing.T) {
	pub := NewPublisher()
	RegisterPublisher(pub)
	sub := NewSubscriber(10)
	RegisterSubscriber(sub)

	sub.Subscribe(pub.ID, "test-topic")

	msgData := []byte("Hello, World!")
	pub.Publish([]string{"test-topic"}, msgData)

	select {
	case msg := <-sub.MsgChan:
		if string(msg.Data) != string(msgData) {
			t.Errorf("Message data does not match: got %s, want %s", string(msg.Data), string(msgData))
		}
	case <-time.After(time.Second):
		t.Error("Timeout waiting for message")
	}
}

func TestSubscribeAndUnsubscribeTopic(t *testing.T) {
	sub := NewSubscriber(10)
	pubID := randPubID()

	// Subscribe to a topic
	sub.Subscribe(pubID, "topic1")
	if len(sub.Subscriptions[pubID]) != 1 || sub.Subscriptions[pubID]["topic1"] != struct{}{} {
		t.Error("Subscription failed")
	}

	// Unsubscribe from the topic
	sub.UnsubscribeTopic(pubID, []string{"topic1"})
	if _, ok := sub.Subscriptions[pubID]; ok && len(sub.Subscriptions[pubID]) != 0 {
		t.Error("Unsubscription failed")
	}
}

func TestUnsubscribePublisher(t *testing.T) {
	sub := NewSubscriber(10)
	pubID := randPubID()
	sub.Subscribe(pubID, "topic1")
	sub.Subscribe(pubID, "topic2")

	// Unsubscribe from all topics of the publisher
	sub.UnsubscribePublisher(pubID)
	if _, exists := sub.Subscriptions[pubID]; exists {
		t.Error("UnsubscribePublisher failed")
	}
}
