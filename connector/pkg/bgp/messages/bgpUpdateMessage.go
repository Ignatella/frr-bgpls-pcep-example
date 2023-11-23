package messages

import (
	"fmt"
	"log"
	"net"
)

const UpdateMessageType = 2

type UpdateMessage struct {
	Type         uint8
	Prefix       net.IP
	PrefixLength uint8
}

func NewUpdateMessage(data []byte) (*UpdateMessage, error) {

	log.Printf("Parsing update message")

	pathAttributeLengthOffset := 21
	pathAttributeLengthBytes := data[pathAttributeLengthOffset : pathAttributeLengthOffset+2]
	pathAttributeLength := uint16(pathAttributeLengthBytes[0])<<8 | uint16(pathAttributeLengthBytes[1])
	log.Printf("Path attribute length: %d\n", pathAttributeLength)

	nlriOffset := pathAttributeLengthOffset + 2 + int(pathAttributeLength)
	prefixLength := data[nlriOffset]
	prefix := net.IP(data[nlriOffset+1:])

	log.Printf("Prefix: %s/%d\n", prefix, prefixLength)

	return &UpdateMessage{Type: UpdateMessageType, Prefix: prefix, PrefixLength: prefixLength}, nil
}

func (m *UpdateMessage) String() string {
	return fmt.Sprintf("UpdateMessage{Type: %d, Prefix: %s/%d}", m.Type, m.Prefix, m.PrefixLength)
}
