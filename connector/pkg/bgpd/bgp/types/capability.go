// Package types contains types used in BGP protocol.
package types

import "fmt"

// Capability codes
const (
	CapabilityMultiprotocol = 1
	CapabilityHostname      = 73
)

// AFI and SAFI codes
const (
	BgplsAfi  = 16388
	BgplsSafi = 71
)

// Capabilities is a slice of capabilities.
type Capabilities []Capability

// Capability is an abstraction for all BGP capability.
type Capability interface {
	GetType() uint8
}

// BaseCapability is a base struct for all capabilities containing BGP capability type.
type BaseCapability struct {
	Type uint8
}

// MultiprotocolCapability is a struct containing AFI and SAFI.
// It represents BGP Multiprotocol capability defined in [RFC 2858].
//
// [RFC 2858]: https://datatracker.ietf.org/doc/html/rfc2858#section-2
type MultiprotocolCapability struct {
	BaseCapability
	AFI  uint16
	SAFI uint8
}

// HostnameCapability is a struct containing hostname.
type HostnameCapability struct {
	BaseCapability
	Hostname string
}

// NewMultiprotocolCapability parses raw bytes into a new MultiprotocolCapability.
func NewMultiprotocolCapability(capabilityBytes []byte) *MultiprotocolCapability {
	return &MultiprotocolCapability{
		BaseCapability: BaseCapability{Type: CapabilityMultiprotocol},
		AFI:            uint16(capabilityBytes[2])<<8 | uint16(capabilityBytes[3]),
		SAFI:           capabilityBytes[5],
	}
}

// NewHostnameCapability parses raw bytes into a new HostnameCapability.
func NewHostnameCapability(capabilityBytes []byte) *HostnameCapability {
	hostnameLength := capabilityBytes[2]
	hostname := string(capabilityBytes[3 : 3+int(hostnameLength)])

	return &HostnameCapability{
		BaseCapability: BaseCapability{Type: CapabilityHostname},
		Hostname:       hostname,
	}
}

// HasAfiSafi checks if the router has Multiprotocol capability with given AFI and SAFI.
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

// GetHostname returns hostname of the router.
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

// GetType returns capability type.
func (bc BaseCapability) GetType() uint8 {
	return bc.Type
}

// String returns a string representation of the capability.
func (bc BaseCapability) String() string {
	return fmt.Sprintf("Capability{Type: %d}", bc.Type)
}

// String returns a string representation of the multiprotocol capability.
func (mc *MultiprotocolCapability) String() string {
	return fmt.Sprintf("Capability{Type: %d, AFI: %d, SAFI: %d}", mc.Type, mc.AFI, mc.SAFI)
}

// String returns a string representation of the hostname capability.
func (hc *HostnameCapability) String() string {
	return fmt.Sprintf("Capability{Type: %d, Hostname: %s}", hc.Type, hc.Hostname)
}
