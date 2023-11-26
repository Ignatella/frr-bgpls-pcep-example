package bgprouter

import (
	"connector/pkg/bgpd/bgp/types"
	"fmt"
	"net"
)

type Router struct {
	conn          net.Conn
	routerEventCh chan types.RouterEvent
	running       chan bool

	AS            uint16
	BGPIdentifier string
	Prefixes      []string
	Capabilities  types.Capabilities
}

func NewRouter(conn net.Conn) *Router {
	return &Router{
		conn:          conn,
		routerEventCh: make(chan types.RouterEvent),
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

func (router *Router) Copy() *Router {
	dst := NewRouter(nil)
	dst.AS = router.AS
	dst.BGPIdentifier = router.BGPIdentifier
	dst.Prefixes = append(dst.Prefixes, router.Prefixes...)
	dst.Capabilities = append(dst.Capabilities, router.Capabilities...)

	return dst
}

func (router *Router) String() string {
	return fmt.Sprintf("Received Router{AS: %d, BGPIdentifier: %s, Prefixes: %v, Capabilities: %v}", router.AS, router.BGPIdentifier, router.Prefixes, router.Capabilities)
}
