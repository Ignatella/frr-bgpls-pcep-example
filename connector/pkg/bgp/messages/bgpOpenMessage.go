package messages

import (
	"fmt"
	"log"
	"net"
)

const OpenMessageType = 1

type OpenMessage struct {
	Type          uint8
	AS            uint16
	BGPIdentifier net.IP
	Hostname      string
}

func NewOpenMessage(data []byte) (*OpenMessage, error) {

	log.Printf("Parsing open message")

	// AS number
	asNumberOffset := 20
	asNumberBytes := data[asNumberOffset : asNumberOffset+2]
	asNumber := uint16(asNumberBytes[0])<<8 | uint16(asNumberBytes[1])
	log.Printf("AS number: %d\n", asNumber)

	// IP address of the router
	bgpIdentifierOffset := 24
	bgpIdentifierBytes := data[bgpIdentifierOffset : bgpIdentifierOffset+4]
	bpgIdentifier := net.IP(bgpIdentifierBytes)
	log.Printf("BGP identifier: %s\n", bpgIdentifier)

	// Hostname
	hostname := ""

	// Go through optional parameters
	optionalParametersOffset := 29
	optionalParametersBytes := data[optionalParametersOffset:]

	for i := 0; i < len(optionalParametersBytes); {
		parameterType := optionalParametersBytes[i]
		parameterLength := optionalParametersBytes[i+1]
		parameterValue := optionalParametersBytes[i+2 : i+2+int(parameterLength)]

		if parameterType != 2 {
			log.Printf("Parameter type: %d\n met. Not expecting. Expecting 2 (Capability)", parameterType)
			i += 2 + int(parameterLength)
			continue
		}

		// Parse only FQDN capability
		capabilityCode := parameterValue[0]
		if capabilityCode == 73 {
			hostnameLength := parameterValue[2]
			hostname = string(parameterValue[3 : 3+int(hostnameLength)])
		}

		i += 2 + int(parameterLength)
	}

	log.Printf("Hostname: %s\n", hostname)

	return &OpenMessage{Type: OpenMessageType, AS: asNumber, BGPIdentifier: bpgIdentifier, Hostname: hostname}, nil
}

func (m *OpenMessage) String() string {
	return fmt.Sprintf("OpenMessage{Type: %d, AS: %d, BGPIdentifier: %s, Hostname: %s}", m.Type, m.AS, m.BGPIdentifier, m.Hostname)
}
