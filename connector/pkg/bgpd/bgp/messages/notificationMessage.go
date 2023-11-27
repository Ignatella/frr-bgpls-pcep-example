package messages

import "fmt"

const NotificationMessageType = 3

type NotificationMessage struct {
	Type uint8
}

func NewNotificationMessage(data []byte) (*NotificationMessage, error) {
	return &NotificationMessage{Type: KeepAliveMessageType}, nil
}

func (m *KeepAliveMessage) String() string {
	return fmt.Sprintf("NotificationMessage")
}
