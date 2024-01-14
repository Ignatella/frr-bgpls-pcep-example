package messages

import "fmt"

// NotificationMessageType is type of BGP Notification message defined in [RFC 4271].
//
// [RFC 4271]: https://datatracker.ietf.org/doc/html/rfc4271#section-4.5
const NotificationMessageType = 3

// NotificationMessage is a message sent by a BGP speaker to inform about an error.
type NotificationMessage struct {
	Type uint8
}

// NewNotificationMessage parses raw bytes into a new NotificationMessage.
func NewNotificationMessage(data []byte) (*NotificationMessage, error) {
	return &NotificationMessage{Type: KeepAliveMessageType}, nil
}

// String returns a string representation of the message.
func (m *KeepAliveMessage) String() string {
	return fmt.Sprintf("NotificationMessage")
}
