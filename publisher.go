package pubsub

import (
	"sync"
	"time"
)

type Publisher struct {
	ID          PubSubID
	Subscribers map[string]map[PubSubID]struct{}
	mu          sync.RWMutex
}

func NewPublisher() *Publisher {
	return &Publisher{
		ID:          randPubID(),
		Subscribers: make(map[string]map[PubSubID]struct{}),
	}
}

func (p *Publisher) AddSubscriber(subID PubSubID, topic string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if _, ok := p.Subscribers[topic]; !ok {
		p.Subscribers[topic] = make(map[PubSubID]struct{})
	}

	p.Subscribers[topic][subID] = struct{}{}
}

func (p *Publisher) RemoveSubscriber(subID PubSubID, topics []string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, topic := range topics {
		delete(p.Subscribers[topic], subID)
	}
}

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
