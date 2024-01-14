package router

import (
	"connector/internal/filereader"
	bgplsTypes "connector/pkg/bgpd/bgpls/types"
	"connector/pkg/bgpd/events"
	"log"
)

// Message templates paths.
const (
	openMessagePath        = "bgpOpenMessage.txt"
	keepAliveMessagePath   = "bgpKeepAliveMessage.txt"
	emptyUpdateMessagePath = "bgpEmptyUpdate.txt"
)

// RunControlThread runs a control thread for the router and listens for the router events.
func (router *Router) RunControlThread(bgpdEventCh chan events.BGPdEvent) {
	for {
		select {
		case <-router.running:
			return
		case event := <-router.RouterEventCh:
			switch event.Type {
			case events.NewRouter:
				router := event.Data.(Router)
				log.Printf("New router: %s\n", &router)
			case events.OpenMessageReceived:
				router.writeMessageFromFile(openMessagePath)
				router.writeMessageFromFile(keepAliveMessagePath)
			case events.KeepAliveMessageReceived:
				router.writeMessageFromFile(keepAliveMessagePath)
			case events.UpdateMessageReceived:
				router.writeMessageFromFile(emptyUpdateMessagePath)
			case events.SendBGPLSNodeNLRI:
				nodeNlri := event.Data.([]bgplsTypes.NodeNLRI)
				err := router.sendBGPLSNodeNLRI(nodeNlri)
				if err != nil {
					router.Exit()
					return
				}
			case events.SendBGPLSPrefixNLRI:
				prefixNlri := event.Data.([]bgplsTypes.PrefixNLRI)
				err := router.sendBGPLSPrefixNLRI(prefixNlri)
				if err != nil {
					router.Exit()
					return
				}
			case events.SendBGPLSLinkNLRI:
				linkNlri := event.Data.([]bgplsTypes.LinkNLRI)
				err := router.sendBGPLSLinkNLRI(linkNlri)
				if err != nil {
					router.Exit()
					return
				}
			case events.Quit, events.NotificationMessageReceived:
				router.Exit()
				return
			default:
				log.Fatalf("Unknown event type: %d\n", event.Type)
			}
		}
	}
}

// WriteMessageFromFile reads a message from a file and writes it to the connection.
func (router *Router) writeMessageFromFile(path string) {
	messageBytes, err := filereader.ReadMessage(path)
	if err != nil {
		router.Exit()
		return
	}

	_, err = router.conn.Write(messageBytes)
	if err != nil {
		router.Exit()
		return
	}
}
