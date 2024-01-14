// Package bgpd implements the BGP daemon.
package bgpd

import (
	"connector/pkg/bgpd/events"
	"connector/pkg/bgpd/server"
)

// StartDaemon starts the BGP daemon.
func StartDaemon(config *server.Config) {
	eventCh := make(chan events.BGPdEvent)

	server.NewBgpControlThread(eventCh).Run()
	server.NewTCPServer(config).Run(eventCh)
}
