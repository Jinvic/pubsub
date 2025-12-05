package pubsub

import (
	"time"
)

// Message represents a message in the pubsub system.
type Message struct {
	// ID is the unique identifier for this message.
	ID PubSubID
	// Data contains the payload of the message.
	Data []byte
	// Topics is a list of topics that this message belongs to.
	Topics []string
	// Publisher is the ID of the publisher that sent this message.
	Publisher PubSubID
	// CreatedAt is the timestamp when the message was created.
	CreatedAt time.Time
}
