// Package router contains types and logic of "Router" (BGP peer).
package router

import (
	"bytes"
	"connector/pkg/bgpd/bgp/types"
	"connector/pkg/bgpd/events"
	"fmt"
	"net"
)

// Router represents a BGP peer.
type Router struct {
	conn    net.Conn
	running chan bool

	RouterEventCh chan events.RouterEvent

	AS            uint16
	BGPIdentifier string
	Prefixes      []string
	Capabilities  types.Capabilities
}

// NewRouter initializes a new Router.
func NewRouter(conn net.Conn) *Router {
	return &Router{
		conn:          conn,
		RouterEventCh: make(chan events.RouterEvent),
		running:       make(chan bool),
		Prefixes:      make([]string, 0),
	}
}

// AddPrefix adds an announced by the peer prefix to the router.
func (router *Router) AddPrefix(prefix string) {
	for _, p := range router.Prefixes {
		if p == prefix {
			return
		}
	}

	router.Prefixes = append(router.Prefixes, prefix)
}

// RemovePrefix removes a prefix from the router.
func (router *Router) RemovePrefix(prefix string) {
	for i, p := range router.Prefixes {
		if p == prefix {
			router.Prefixes = append(router.Prefixes[:i], router.Prefixes[i+1:]...)
			return
		}
	}
}

// GetAdministrativePrefix returns the administrative prefix of the router.
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

// HasMultiprotocolCapability checks if the router has a multiprotocol capability.
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

// GetHostname returns the hostname of the router.
func (router *Router) GetHostname() string {
	for _, capability := range router.Capabilities {
		if c, ok := capability.(*types.HostnameCapability); ok {
			return c.Hostname
		}
	}

	return ""
}

// Exit closes the running channel effectively stopping the router and breaking connection with the peer.
func (router *Router) Exit() {
	select {
	case <-router.running:
		return
	default:
		close(router.running)
	}
}

// String returns a string representation of the router.
func (router *Router) String() string {
	return fmt.Sprintf("Received Router{AS: %d, BGPIdentifier: %s, Prefixes: %v, Capabilities: %v}", router.AS, router.BGPIdentifier, router.Prefixes, router.Capabilities)
}
