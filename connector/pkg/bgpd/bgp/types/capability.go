package types

import "fmt"

const (
	CapabilityMultiprotocol = 1
	CapabilityHostname      = 73
)

const (
	BgplsAfi  = 16388
	BgplsSafi = 71
)

type Capabilities []Capability

type Capability interface {
	GetType() uint8
}

type BaseCapability struct {
	Type uint8
}

type MultiprotocolCapability struct {
	BaseCapability
	AFI  uint16
	SAFI uint8
}

type HostnameCapability struct {
	BaseCapability
	Hostname string
}

func NewMultiprotocolCapability(capabilityBytes []byte) *MultiprotocolCapability {
	return &MultiprotocolCapability{
		BaseCapability: BaseCapability{Type: CapabilityMultiprotocol},
		AFI:            uint16(capabilityBytes[2])<<8 | uint16(capabilityBytes[3]),
		SAFI:           capabilityBytes[5],
	}
}

func NewHostnameCapability(capabilityBytes []byte) *HostnameCapability {
	hostnameLength := capabilityBytes[2]
	hostname := string(capabilityBytes[3 : 3+int(hostnameLength)])

	return &HostnameCapability{
		BaseCapability: BaseCapability{Type: CapabilityHostname},
		Hostname:       hostname,
	}
}

func (capabilities *Capabilities) HasAfiSafi(afi uint16, safi uint8) bool {
	for _, capability := range *capabilities {
		switch capability.(type) {
		case MultiprotocolCapability:
			c := capability.(MultiprotocolCapability)
			if c.AFI == afi && c.SAFI == safi {
				return true
			}
		}
	}

	return false
}

func (capabilities *Capabilities) GetHostname() string {
	for _, capability := range *capabilities {
		switch capability.(type) {
		case HostnameCapability:
			c := capability.(HostnameCapability)
			return c.Hostname
		}
	}

	return ""
}

// GetType

func (bc BaseCapability) GetType() uint8 {
	return bc.Type
}

// String

func (bc BaseCapability) String() string {
	return fmt.Sprintf("Capability{Type: %d}", bc.Type)
}

func (mc *MultiprotocolCapability) String() string {
	return fmt.Sprintf("Capability{Type: %d, AFI: %d, SAFI: %d}", mc.Type, mc.AFI, mc.SAFI)
}

func (hc *HostnameCapability) String() string {
	return fmt.Sprintf("Capability{Type: %d, Hostname: %s}", hc.Type, hc.Hostname)
}
