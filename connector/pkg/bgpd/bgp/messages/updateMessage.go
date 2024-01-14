package messages

import (
	"fmt"
	"log"
	"net"
)

// UpdateMessageType is type of BGP Update message defined in [RFC 4271].
//
// [RFC 4271]: https://datatracker.ietf.org/doc/html/rfc4271#section-4.3
const UpdateMessageType = 2

// UpdateMessage is a message sent by a BGP speaker to update the routing information.
// It contains information about the speaker's prefixes.
type UpdateMessage struct {
	Type   uint8
	Prefix net.IPNet
}

// NewUpdateMessage parses raw bytes into a new UpdateMessage.
func NewUpdateMessage(data []byte) (*UpdateMessage, error) {

	log.Printf("Parsing update message")

	pathAttributeLengthOffset := 21
	pathAttributeLengthBytes := data[pathAttributeLengthOffset : pathAttributeLengthOffset+2]
	pathAttributeLength := uint16(pathAttributeLengthBytes[0])<<8 | uint16(pathAttributeLengthBytes[1])

	nlriOffset := pathAttributeLengthOffset + 2 + int(pathAttributeLength)
	prefixLength := data[nlriOffset]
	prefix := net.IP(data[nlriOffset+1:])

	log.Printf("Prefix: %s/%d\n", prefix, prefixLength)

	return &UpdateMessage{
		Type: UpdateMessageType,
		Prefix: net.IPNet{
			IP:   prefix,
			Mask: net.CIDRMask(int(prefixLength), 32),
		},
	}, nil
}

// String returns a string representation of the message.
func (m *UpdateMessage) String() string {
	return fmt.Sprintf("UpdateMessage{Type: %d, Prefix: %s}", m.Type, m.Prefix)
}
