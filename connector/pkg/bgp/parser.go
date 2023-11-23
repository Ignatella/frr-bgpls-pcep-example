package bgp

import (
	"connector/pkg/bgp/messages"
	"errors"
	"log"
)

func ParseMessage(data []byte) (messages.Message, error) {
	// Check BGP marker
	bgpMarker := [16]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	if !equalSlices(data[:16], bgpMarker[:]) {
		return nil, errors.New("invalid BGP marker")
	}

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
	default:
		return nil, errors.New("unknown BGP message type")
	}
}

func equalSlices(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}

	for i, _ := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
