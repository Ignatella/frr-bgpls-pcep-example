// Package bgp contains types and logic of BGP protocol.
package bgp

import (
	"connector/pkg/bgpd/bgp/messages"
	"errors"
	"log"
)

// ParseMessage parses raw bytes into a new BGP message.
func ParseMessage(data []byte) (messages.Message, error) {
	messageTypeOffset := 18
	messageType := data[messageTypeOffset]

	log.Printf("Message type: %d\n", int(messageType))

	switch messageType {
	case messages.OpenMessageType:
		message, err := messages.NewOpenMessage(data)
		return message, err
	case messages.KeepAliveMessageType:
		message, err := messages.NewKeepAliveMessage(data)
		return message, err
	case messages.UpdateMessageType:
		message, err := messages.NewUpdateMessage(data)
		return message, err
	case messages.NotificationMessageType:
		message, err := messages.NewNotificationMessage(data)
		return message, err
	default:
		return nil, errors.New("unknown BGP message type")
	}
}
