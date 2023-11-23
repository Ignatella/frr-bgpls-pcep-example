package bgp

import (
	"connector/pkg/bgp/messages"
	"errors"
	"log"
)

func ParseMessage(data []byte) (messages.Message, error) {
	// Check BGP message type
	messageTypeOffset := 18
	messageType := data[messageTypeOffset]
	log.Printf("Message type: %d\n", int(messageType))

	switch messageType {
	case messages.OpenMessageType:
		message, err := messages.NewOpenMessage(data)
		log.Printf("Open message: %s\n", message)
		return message, err
	case messages.KeepAliveMessageType:
		message, err := messages.NewKeepAliveMessage(data)
		log.Printf("KeepAlive message: %s\n", message)
		return message, err
	case messages.UpdateMessageType:
		message, err := messages.NewUpdateMessage(data)
		log.Printf("Update message: %s\n", message)
		return message, err
	default:
		return nil, errors.New("unknown BGP message type")
	}
}
