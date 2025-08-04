package pubsub

import (
	"time"
)

type Message struct {
	ID        PubSubID
	Data      []byte
	Topics    []string
	Publisher PubSubID
	CreatedAt time.Time
}
