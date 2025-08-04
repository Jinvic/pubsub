package pubsub

import (
	"fmt"

	"github.com/google/uuid"
)

type PubSubID string

func randPubID() PubSubID {
	return PubSubID(fmt.Sprintf("pub-%s", uuid.New().String()))
}

func randSubID() PubSubID {
	return PubSubID(fmt.Sprintf("sub-%s", uuid.New().String()))
}

func randMsgID() PubSubID {
	return PubSubID(fmt.Sprintf("msg-%s", uuid.New().String()))
}

var Publishers = make(map[PubSubID]*Publisher)
var Subscribers = make(map[PubSubID]*Subscriber)

func RegisterPublisher(p *Publisher) {
	Publishers[p.ID] = p
}

func RegisterSubscriber(s *Subscriber) {
	Subscribers[s.ID] = s
}

func UnregisterPublisher(p *Publisher) {
	delete(Publishers, p.ID)
}

func UnregisterSubscriber(s *Subscriber) {
	delete(Subscribers, s.ID)
}

func GetPublisher(id PubSubID) *Publisher {
	return Publishers[id]
}

func GetSubscriber(id PubSubID) *Subscriber {
	return Subscribers[id]
}
