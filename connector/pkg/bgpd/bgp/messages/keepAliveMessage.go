package messages

import (
	"fmt"
)

// KeepAliveMessageType is type of BGP KeepAlive message defined in [RFC 4271].
//
// [RFC 4271]: https://datatracker.ietf.org/doc/html/rfc4271#section-4.4
const KeepAliveMessageType = 4

// KeepAliveMessage is a message sent by a BGP speaker to maintain the connection.
// Struct contains fields related to BGP Keepalive message.
type KeepAliveMessage struct {
	Type uint8
}

// NewKeepAliveMessage parses raw bytes into a new KeepAliveMessage.
func NewKeepAliveMessage(data []byte) (*KeepAliveMessage, error) {
	return &KeepAliveMessage{Type: KeepAliveMessageType}, nil
}

// String returns a string representation of the message.
func (m *NotificationMessage) String() string {
	return fmt.Sprintf("KeepAliveMessage")
}
