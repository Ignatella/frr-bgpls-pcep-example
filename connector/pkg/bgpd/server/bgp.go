package server

import (
	bgpTypes "connector/pkg/bgpd/bgp/types"
	"connector/pkg/bgpd/bgpRouter"
	bgplsTypes "connector/pkg/bgpd/bgpls/types"
	"log"
)

type BgpControlThread struct {
	bgpdEventCh chan bgpTypes.BGPdEvent
	routers     map[string]*bgprouter.Router
}

func NewBgpControlThread(bgpdEventCh chan bgpTypes.BGPdEvent) *BgpControlThread {
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
			case bgpTypes.NewRouter:
				router := event.Data.(*bgprouter.Router)
				th.routers[router.BGPIdentifier] = router
				log.Printf("New router: %s\n", router)

				if th.isCompleteTopology() {
					log.Printf("Complete topology detected")

					for _, r := range th.routers {
						if r.HasMultiprotocolCapability(bgpTypes.BgplsAfi, bgpTypes.BgplsSafi) {
							r.RouterEventCh <- bgpTypes.RouterEvent{
								Type: bgpTypes.SendBGPLSTopology,
								Data: th.CalculateLinkstateNodeNLRI(),
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
		result = append(result, bgplsTypes.NodeNLRI{AS: router.AS, RouterId: router.BGPIdentifier, Hostname: router.GetHostname()})
	}

	return result
}

//
//func (th *BgpControlThread) CalculateLinkstateLinkNLRI() []bgplsTypes.LinkNLRI {
//	// local | as
//	// local | igp router id
//
//	// remote | as
//	// remote | igp router id
//
//	// link local
//	// link remote
//
//	// adj sid
//
//}

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
