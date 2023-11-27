package bgprouter

import (
	"bytes"
	"connector/pkg/bgpd/bgp/types"
	"fmt"
	"net"
)

type Router struct {
	conn    net.Conn
	running chan bool

	RouterEventCh chan types.RouterEvent

	AS            uint16
	BGPIdentifier string
	Prefixes      []string
	Capabilities  types.Capabilities
}

func NewRouter(conn net.Conn) *Router {
	return &Router{
		conn:          conn,
		RouterEventCh: make(chan types.RouterEvent),
		running:       make(chan bool),
		Prefixes:      make([]string, 0),
	}
}

func (router *Router) AddPrefix(prefix string) {
	for _, p := range router.Prefixes {
		if p == prefix {
			return
		}
	}

	router.Prefixes = append(router.Prefixes, prefix)
}

func (router *Router) RemovePrefix(prefix string) {
	for i, p := range router.Prefixes {
		if p == prefix {
			router.Prefixes = append(router.Prefixes[:i], router.Prefixes[i+1:]...)
			return
		}
	}
}

func (router *Router) GetAdministrativePrefix() (string, bool) {
	for _, prefix := range router.Prefixes {
		_, ipnet, err := net.ParseCIDR(prefix)

		if err != nil {
			continue
		}

		if bytes.Equal(ipnet.Mask, net.CIDRMask(32, 32)) {
			return prefix, true
		}
	}

	return "", false
}

func (router *Router) HasMultiprotocolCapability(afi uint16, safi uint8) bool {
	for _, capability := range router.Capabilities {
		if c, ok := capability.(*types.MultiprotocolCapability); ok {
			if c.AFI == afi && c.SAFI == safi {
				return true
			}
		}
	}

	return false
}

func (router *Router) GetHostname() string {
	for _, capability := range router.Capabilities {
		if c, ok := capability.(*types.HostnameCapability); ok {
			return c.Hostname
		}
	}

	return ""
}

func (router *Router) Exit() {
	select {
	case <-router.running:
		return
	default:
		close(router.running)
	}
}

func (router *Router) String() string {
	return fmt.Sprintf("Received Router{AS: %d, BGPIdentifier: %s, Prefixes: %v, Capabilities: %v}", router.AS, router.BGPIdentifier, router.Prefixes, router.Capabilities)
}
