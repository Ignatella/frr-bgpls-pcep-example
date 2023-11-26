package messages

import (
	"connector/pkg/bgpd/bgp/types"
	"fmt"
	"log"
	"net"
)

const OpenMessageType = 1

type OpenMessage struct {
	Type          uint8
	AS            uint16
	BGPIdentifier net.IP
	Capabilities  types.Capabilities
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

	// Capabilities
	capabilities := make(types.Capabilities, 0)

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

		for j := 0; j < len(parameterValue); {
			capabilityType := parameterValue[j]
			capabilityLength := parameterValue[j+1]

			capabilityBytes := parameterValue[j : j+2+int(capabilityLength)]

			switch capabilityType {
			case types.CapabilityHostname:
				c := types.NewHostnameCapability(capabilityBytes)
				capabilities = append(capabilities, c)
			case types.CapabilityMultiprotocol:
				c := types.NewMultiprotocolCapability(capabilityBytes)
				capabilities = append(capabilities, c)
			}

			j += 2 + int(capabilityLength)
		}

		i += 2 + int(parameterLength)
	}

	return &OpenMessage{Type: OpenMessageType, AS: asNumber, BGPIdentifier: bpgIdentifier, Capabilities: capabilities}, nil
}

func (m *OpenMessage) String() string {
	return fmt.Sprintf("OpenMessage{Type: %d, AS: %d, BGPIdentifier: %s, Capabilities: %s}", m.Type, m.AS, m.BGPIdentifier, m.Capabilities)
}
