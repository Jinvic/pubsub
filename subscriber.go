package pubsub

import "sync"

type Subscriber struct {
	ID            PubSubID
	Subscriptions map[PubSubID]map[string]struct{}
	MsgChan       chan Message
	mu            sync.RWMutex
}

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

func (s *Subscriber) Receive(msg Message) {
	s.MsgChan <- msg
}

func (s *Subscriber) Consume() Message {
	return <-s.MsgChan
}

func (s *Subscriber) TryConsume() (Message, bool) {
	select {
	case msg := <-s.MsgChan:
		return msg, true
	default:
		return Message{}, false
	}
}

func (s *Subscriber) Close() {
	close(s.MsgChan)
	UnregisterSubscriber(s)
}

func (s *Subscriber) Unregister() {
	UnregisterSubscriber(s)
}
