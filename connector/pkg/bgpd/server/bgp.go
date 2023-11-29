package server

import (
	bgpTypes "connector/pkg/bgpd/bgp/types"
	"connector/pkg/bgpd/bgpRouter"
	bgplsTypes "connector/pkg/bgpd/bgpls/types"
	"connector/pkg/bgpd/events"
	"log"
	"net"
)

type BgpControlThread struct {
	bgpdEventCh chan events.BGPdEvent
	routers     map[string]*bgprouter.Router
}

func NewBgpControlThread(bgpdEventCh chan events.BGPdEvent) *BgpControlThread {
	return &BgpControlThread{
		bgpdEventCh: bgpdEventCh,
		routers:     make(map[string]*bgprouter.Router),
	}
}

func (th *BgpControlThread) Run() {
	go func() {
		for {
			event := <-th.bgpdEventCh

			switch event.Type {
			case events.NewRouter:
				router := event.Data.(*bgprouter.Router)
				th.routers[router.BGPIdentifier] = router
				log.Printf("New router: %s\n", router)

				if th.isCompleteTopology() {
					log.Printf("Complete topology detected")

					for _, r := range th.routers {
						if r.HasMultiprotocolCapability(bgpTypes.BgplsAfi, bgpTypes.BgplsSafi) {
							r.RouterEventCh <- events.RouterEvent{
								Type: events.SendBGPLSNodeNLRI,
								Data: th.CalculateLinkstateNodeNLRI(),
							}
							r.RouterEventCh <- events.RouterEvent{
								Type: events.SendBGPLSPrefixNLRI,
								Data: th.CalculateLinkstatePrefixNLRI(),
							}
							r.RouterEventCh <- events.RouterEvent{
								Type: events.SendBGPLSLinkNLRI,
								Data: th.CalculateLinkstateLinkNLRI(),
							}
						}
					}
				}
			default:
				log.Fatalf("Unknown event type: %d\n", event.Type)
			}
		}
	}()
}

func (th *BgpControlThread) CalculateLinkstateNodeNLRI() []bgplsTypes.NodeNLRI {

	result := make([]bgplsTypes.NodeNLRI, 0)

	for _, router := range th.routers {

		if router.HasMultiprotocolCapability(bgpTypes.BgplsAfi, bgpTypes.BgplsSafi) {
			continue
		}

		result = append(result, bgplsTypes.NodeNLRI{AS: router.AS, RouterId: router.BGPIdentifier, Hostname: router.GetHostname()})
	}

	return result
}

func (th *BgpControlThread) CalculateLinkstatePrefixNLRI() []bgplsTypes.PrefixNLRI {
	result := make([]bgplsTypes.PrefixNLRI, 0)

	for _, router := range th.routers {
		result = append(result, bgplsTypes.PrefixNLRI{AS: router.AS, RouterId: router.BGPIdentifier})
	}

	return result
}

func (th *BgpControlThread) CalculateLinkstateLinkNLRI() []bgplsTypes.LinkNLRI {
	linkNLRIs := make(map[string]bgplsTypes.LinkNLRI)

	prefixNodes := make(map[string]int)

	for _, router := range th.routers {
		if router.HasMultiprotocolCapability(bgpTypes.BgplsAfi, bgpTypes.BgplsSafi) {
			continue
		}

		adminPrefix, found := router.GetAdministrativePrefix()

		if !found {
			log.Printf("E! Router %s has no administrative prefix", router.BGPIdentifier)
			continue
		}

		for _, prefix := range router.Prefixes {
			if prefix == adminPrefix {
				continue
			}

			ip, _, err := net.ParseCIDR(prefix)
			if err != nil {
				log.Printf("E! Error parsing prefix %s: %s", prefix, err)
				continue
			}

			ip = ip.To4()

			prefixNodes[prefix]++

			ip[3] += byte(prefixNodes[prefix])

			if _, ok := linkNLRIs[prefix]; !ok {
				linkNLRIs[prefix] = bgplsTypes.LinkNLRI{Endpoints: make([]bgplsTypes.LinkEndpoint, 0)}
			}

			endpoint := linkNLRIs[prefix]

			endpoint.Endpoints = append(endpoint.Endpoints, bgplsTypes.LinkEndpoint{AS: router.AS, RouterId: router.BGPIdentifier, NeighborPrefix: ip.String()})

			linkNLRIs[prefix] = endpoint
		}
	}

	nlris := make([]bgplsTypes.LinkNLRI, 0)
	for _, nlri := range linkNLRIs {
		// a -> b
		nlris = append(nlris, nlri)
		// b -> a
		nlris = append(nlris, bgplsTypes.LinkNLRI{Endpoints: []bgplsTypes.LinkEndpoint{nlri.Endpoints[1], nlri.Endpoints[0]}})
	}

	return nlris
}

func (th *BgpControlThread) isCompleteTopology() bool {

	// number of nodes on links between routers
	prefixNodes := make(map[string]int)

	for _, router := range th.routers {

		if router.HasMultiprotocolCapability(bgpTypes.BgplsAfi, bgpTypes.BgplsSafi) {
			continue
		}

		adminPrefix, found := router.GetAdministrativePrefix()

		if !found {
			return false
		}

		for _, prefix := range router.Prefixes {
			if prefix == adminPrefix {
				continue
			}

			prefixNodes[prefix]++
		}
	}

	// check if all links have two nodes (all routers are connected)
	for _, nodes := range prefixNodes {
		if nodes != 2 {
			return false
		}
	}

	return true
}
