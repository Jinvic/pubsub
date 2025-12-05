// Package pubsub provides a lightweight, concurrent-safe publish-subscribe implementation.
package pubsub

import (
	"fmt"

	"github.com/google/uuid"
)

// PubSubID represents a unique identifier for publishers, subscribers, and messages.
type PubSubID string

// randPubID generates a random publisher ID with "pub-" prefix.
func randPubID() PubSubID {
	return PubSubID(fmt.Sprintf("pub-%s", uuid.New().String()))
}

// randSubID generates a random subscriber ID with "sub-" prefix.
func randSubID() PubSubID {
	return PubSubID(fmt.Sprintf("sub-%s", uuid.New().String()))
}

// randMsgID generates a random message ID with "msg-" prefix.
func randMsgID() PubSubID {
	return PubSubID(fmt.Sprintf("msg-%s", uuid.New().String()))
}

// Publishers is a global map of all registered publishers, keyed by their ID.
var Publishers = make(map[PubSubID]*Publisher)

// Subscribers is a global map of all registered subscribers, keyed by their ID.
var Subscribers = make(map[PubSubID]*Subscriber)

// RegisterPublisher registers a publisher in the global Publishers map.
func RegisterPublisher(p *Publisher) {
	Publishers[p.ID] = p
}

// RegisterSubscriber registers a subscriber in the global Subscribers map.
func RegisterSubscriber(s *Subscriber) {
	Subscribers[s.ID] = s
}

// UnregisterPublisher removes a publisher from the global Publishers map.
func UnregisterPublisher(p *Publisher) {
	delete(Publishers, p.ID)
}

// UnregisterSubscriber removes a subscriber from the global Subscribers map.
func UnregisterSubscriber(s *Subscriber) {
	delete(Subscribers, s.ID)
}

// GetPublisher retrieves a publisher by its ID from the global Publishers map.
// Returns nil if the publisher is not found.
func GetPublisher(id PubSubID) *Publisher {
	return Publishers[id]
}

// GetSubscriber retrieves a subscriber by its ID from the global Subscribers map.
// Returns nil if the subscriber is not found.
func GetSubscriber(id PubSubID) *Subscriber {
	return Subscribers[id]
}
