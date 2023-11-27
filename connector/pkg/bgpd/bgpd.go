package bgpd

import (
	"connector/pkg/bgpd/bgp/types"
	"connector/pkg/bgpd/server"
)

func StartDaemon(config *server.Config) {
	eventCh := make(chan types.BGPdEvent)

	server.NewBgpControlThread(eventCh).Run()
	server.NewTCPServer(config).Run(eventCh)
}
