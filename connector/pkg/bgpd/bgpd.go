package bgpd

import (
	"connector/pkg/bgpd/events"
	"connector/pkg/bgpd/server"
)

func StartDaemon(config *server.Config) {
	eventCh := make(chan events.BGPdEvent)

	server.NewBgpControlThread(eventCh).Run()
	server.NewTCPServer(config).Run(eventCh)
}
