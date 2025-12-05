package pubsub

import (
	"context"
	"errors"
	"sync"
)

// Subscriber represents a subscriber in the pubsub system.
// It maintains subscriptions to publishers and their topics, and provides
// methods to consume messages.
type Subscriber struct {
	// ID is the unique identifier for this subscriber.
	ID PubSubID
	// Subscriptions is a map of publisher ID to topic subscriptions.
	Subscriptions map[PubSubID]map[string]struct{}
	// MsgChan is the channel through which messages are received.
	MsgChan chan Message
	// mu provides concurrent-safe access to the subscriber's fields.
	mu sync.RWMutex
}

// ErrSubscriberClosed is returned when trying to consume from a closed subscriber.
var ErrSubscriberClosed = errors.New("subscriber closed")

// NewSubscriber creates a new subscriber.
// bufferSize specifies the size of the message channel buffer. If bufferSize <= 0, it defaults to 10.
// If autoRegister is true, the subscriber is automatically registered in the global Subscribers map.
func NewSubscriber(bufferSize int, autoRegister bool) *Subscriber {
	if bufferSize <= 0 {
		bufferSize = 10
	}
	sub := &Subscriber{
		ID:            randSubID(),
		Subscriptions: make(map[PubSubID]map[string]struct{}),
		MsgChan:       make(chan Message, bufferSize),
	}
	if autoRegister {
		RegisterSubscriber(sub)
	}
	return sub
}

// Subscribe subscribes to a topic of a specific publisher.
// This method is safe for concurrent use.
func (s *Subscriber) Subscribe(pubID PubSubID, topic string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.Subscriptions[pubID]; !ok {
		s.Subscriptions[pubID] = make(map[string]struct{})
	}

	s.Subscriptions[pubID][topic] = struct{}{}
	pub := GetPublisher(pubID)
	if pub != nil {
		pub.AddSubscriber(s.ID, topic)
	}
}

// UnsubscribePublisher unsubscribes from all topics of a specific publisher.
// This method is safe for concurrent use.
func (s *Subscriber) UnsubscribePublisher(pubID PubSubID) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.Subscriptions[pubID]; !ok {
		return
	}

	topics := make([]string, 0, len(s.Subscriptions[pubID]))
	for topic := range s.Subscriptions[pubID] {
		topics = append(topics, topic)
	}
	delete(s.Subscriptions, pubID)

	pub := GetPublisher(pubID)
	if pub != nil {
		pub.RemoveSubscriber(s.ID, topics)
	}
}

// UnsubscribeTopic unsubscribes from specific topics of a publisher.
// This method is safe for concurrent use.
func (s *Subscriber) UnsubscribeTopic(pubID PubSubID, topics []string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.Subscriptions[pubID]; !ok {
		return
	}

	for _, topic := range topics {
		delete(s.Subscriptions[pubID], topic)

		if len(s.Subscriptions[pubID]) == 0 {
			delete(s.Subscriptions, pubID)
		}
	}

	pub := GetPublisher(pubID)
	if pub != nil {
		pub.RemoveSubscriber(s.ID, topics)
	}
}

// Receive receives a message and sends it to the subscriber's message channel.
// This method is safe for concurrent use.
func (s *Subscriber) Receive(msg Message) {
	s.MsgChan <- msg
}

// Consume blocks until a message is received from the message channel.
// Returns the received message and any error that occurred.
func (s *Subscriber) Consume() (Message, error) {
	return s.ConsumeContext(context.Background())
}

// ConsumeContext blocks until a message is received or the context is cancelled.
// Returns the received message and any error that occurred.
func (s *Subscriber) ConsumeContext(ctx context.Context) (Message, error) {
	select {
	case <-ctx.Done():
		return Message{}, ctx.Err()
	case msg, ok := <-s.MsgChan:
		if !ok {
			return Message{}, ErrSubscriberClosed
		}
		return msg, nil
	}
}

// TryConsume attempts to receive a message without blocking.
// Returns the received message, a boolean indicating if a message was received, and any error that occurred.
func (s *Subscriber) TryConsume() (Message, bool, error) {
	select {
	case msg, ok := <-s.MsgChan:
		if !ok {
			return Message{}, false, ErrSubscriberClosed
		}
		return msg, true, nil
	default:
		return Message{}, false, nil
	}
}

// Unregister closes the message channel and removes the subscriber from the global Subscribers map.
func (s *Subscriber) Unregister() {
	close(s.MsgChan)
	UnregisterSubscriber(s)
}
