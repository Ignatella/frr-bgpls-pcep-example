package messages

import (
	"fmt"
)

const KeepAliveMessageType = 4

type KeepAliveMessage struct {
	Type uint8
}

func NewKeepAliveMessage(data []byte) (*KeepAliveMessage, error) {
	return &KeepAliveMessage{Type: KeepAliveMessageType}, nil
}

func (m *NotificationMessage) String() string {
	return fmt.Sprintf("KeepAliveMessage")
}
