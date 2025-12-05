package pubsub

import (
	"sync"
	"time"
)

// Publisher represents a publisher in the pubsub system.
// It maintains a list of subscribers for each topic.
type Publisher struct {
	// ID is the unique identifier for this publisher.
	ID PubSubID
	// Subscribers is a map of topic to subscriber IDs.
	Subscribers map[string]map[PubSubID]struct{}
	// mu provides concurrent-safe access to the publisher's fields.
	mu sync.RWMutex
}

// NewPublisher creates a new publisher.
// If autoRegister is true, the publisher is automatically registered in the global Publishers map.
func NewPublisher(autoRegister bool) *Publisher {
	pub := &Publisher{
		ID:          randPubID(),
		Subscribers: make(map[string]map[PubSubID]struct{}),
	}
	if autoRegister {
		RegisterPublisher(pub)
	}
	return pub
}

// AddSubscriber adds a subscriber to a specific topic.
// This method is safe for concurrent use.
func (p *Publisher) AddSubscriber(subID PubSubID, topic string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if _, ok := p.Subscribers[topic]; !ok {
		p.Subscribers[topic] = make(map[PubSubID]struct{})
	}

	p.Subscribers[topic][subID] = struct{}{}
}

// RemoveSubscriber removes a subscriber from the specified topics.
// This method is safe for concurrent use.
func (p *Publisher) RemoveSubscriber(subID PubSubID, topics []string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, topic := range topics {
		delete(p.Subscribers[topic], subID)
	}
}

// Publish sends a message to all subscribers of the specified topics.
// The message is sent asynchronously using goroutines.
// If an empty string is a topic, all subscribers (regardless of their specific topic subscriptions) will receive the message.
func (p *Publisher) Publish(topics []string, data []byte) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	subIDs := make(map[PubSubID]struct{})
	for _, topic := range topics {
		for subID := range p.Subscribers[topic] {
			subIDs[subID] = struct{}{}
		}
	}
	for subID := range p.Subscribers[""] {
		subIDs[subID] = struct{}{}
	}

	if len(subIDs) == 0 {
		return
	}

	msg := Message{
		ID:        randMsgID(),
		Data:      data,
		Topics:    topics,
		Publisher: p.ID,
		CreatedAt: time.Now(),
	}
	for subID := range subIDs {
		go func(subID PubSubID) {
			sub := GetSubscriber(subID)
			if sub != nil {
				sub.Receive(msg)
			}
		}(subID)
	}
}

// Unregister removes this publisher from the global Publishers map.
func (p *Publisher) Unregister() {
	UnregisterPublisher(p)
}
