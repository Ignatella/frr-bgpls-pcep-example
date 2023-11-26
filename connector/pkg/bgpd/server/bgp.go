package server

import (
	"connector/pkg/bgpd/bgp/types"
	"connector/pkg/bgpd/bgpRouter"
	"log"
)

func StartBGPControlThread(bgpdEventCh chan types.BGPdEvent) {
	go func() {
		routers := make(map[string]bgprouter.Router)

		for {
			event := <-bgpdEventCh

			switch event.Type {
			case types.NewRouter:
				r := event.Data.(bgprouter.Router)
				routers[r.BGPIdentifier] = r
				log.Printf("New router: %s\n", &r)
			default:
				log.Fatalf("Unknown event type: %d\n", event.Type)
			}
		}
	}()
}
